package api

type SlowRequest struct {
	Timeout uint `json:"timeout"`
}

type JsonResponse struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}
