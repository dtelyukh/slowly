Сборка и запуск в ОС Linux

	make build
	make run

Проверка работоспособности

	curl -i \
    -H "Accept: application/json" \
    -H "Content-type: application/json" \
    -X POST \
    -d '{"timeout":1000}' \
    http://localhost:8080/api/slow

Прогон тестов

	make test

Параметры приложения по умолчанию в .env.dist переопределяются в .env
