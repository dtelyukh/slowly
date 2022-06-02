Build and run on Linux OS

	make build
	make run

Usage

	curl -i \
    -H "Accept: application/json" \
    -H "Content-type: application/json" \
    -X POST \
    -d '{"timeout":1000}' \
    http://localhost:8080/api/slow

Test run

	make test

The default application parameters in .env.dist are overridden in .env
