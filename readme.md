# Shorten URL Service

## How to run it

1 Copy environment file

```sh
cp .env.example .env
```

2 Start docker container

```sh
docker-compose up -d
```

3 Import posman collection from file `docs/Shorten.postman_collection.json`

4 Tes call API from postmant

** You can see all log with command

```sh
docker logs -f shorten-url-service_go_1
```
