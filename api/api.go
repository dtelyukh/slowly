package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"slowly/logger"

	"github.com/gorilla/mux"
)

type Config struct {
	MaxTimeout uint
}

type handler struct {
	config *Config
	log    logger.Logger
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc(
		"/api/slow",
		h.throttle(h.slow())).
		Methods("POST")

	router.ServeHTTP(w, r)
}

func NewHandler(log logger.Logger, cfg *Config) Handler {
	if cfg.MaxTimeout == 0 {
		panic("Invalid MaxTimeout value")
	}
	return &handler{config: cfg, log: log}
}

func (h *handler) slow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			h.sendJsonResponse(w, r, http.StatusInternalServerError, err)
			h.log.Errorf("Parsing request body error: %s", err)
			return
		}

		req := SlowRequest{}
		json.Unmarshal(content, &req)
		if req.Timeout == 0 {
			h.sendJsonResponse(w, r, http.StatusBadRequest, errors.New("Invalid timeout value"))
			return
		}

		select {
		case <-time.After(time.Millisecond * time.Duration(req.Timeout)):
			h.sendJsonResponse(w, r, http.StatusOK, err)
			return
		case <-r.Context().Done():
			return
		}

	}
}

func (h *handler) throttle(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.config.MaxTimeout)*time.Millisecond)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))

		if len(w.Header().Get("Content-type")) == 0 {
			h.sendJsonResponse(w, r, http.StatusBadRequest, errors.New("timeout too long"))
		}
	}

	return http.HandlerFunc(fn)
}

func (h *handler) sendJsonResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	var jsonStr []byte
	resp := JsonResponse{}
	if err == nil {
		resp.Status = "ok"
	} else {
		resp.Error = err.Error()
	}
	jsonStr, _ = json.Marshal(resp)

	h.log.Infof("%s %s %s %d %s", r.Method, r.RequestURI, r.Proto, code, r.UserAgent())

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(jsonStr)
}
