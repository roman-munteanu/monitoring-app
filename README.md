monitoring-app
-----

## Init

```
docker-compose up -d

winpty docker exec -it <container_id> bash
```

## Start

```
go mod tidy

go mod vendor

go run main.go
```

GET 
http://localhost:3001/heroes
