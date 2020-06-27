PROJECT_PKGS := $$(go list ./... | grep -v /vendor/)

.SILENT: test
test:
	for pkg in $(PROJECT_PKGS); do \
		go test -race -cover ${VERBOSITY} -coverprofile $$(basename $$pkg).coverprof $$pkg || exit 1 ;\
	done
	echo "mode: set" > coverage.cov
	grep -h -v "mode: *" *.coverprof >> coverage.cov
	go tool cover -html=coverage.cov -o coverage.html
	rm *.coverprof
	rm *.cov
	echo "Profile coverage in coverage.html"

build:
	(cd cmd/slowly && go build)

run:
	touch .env
	(set -a && . ./.env.dist && . ./.env && ./cmd/slowly/slowly && set +a)
