Запуск postgres в docker
```console
docker run --name=postgres -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d --rm postgres
```
