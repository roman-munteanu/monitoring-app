monitoring-app
-----

## Init

Build the app image:
```
docker build --tag monitoring-app .

docker image ls
```

Run all the services:
```
docker-compose up -d

winpty docker exec -it <container_id> bash
```

## Database

```
docker exec -it mariadb /bin/bash

mysql -u root -p

USE monitoringdb;
```
Execute the queries from `schema.sql`

## API

GET 
http://localhost:3001/heroes


## Datadog agent

Check the status:
```
docker exec -it dd-agent agent status
```


## Local start
```
go mod tidy

go mod vendor

go run main.go
```



## Resources

https://docs.docker.com/language/golang/build-images/

https://docs.datadoghq.com/containers/docker/apm/?tab=windows

https://docs.datadoghq.com/agent/guide/compose-and-the-datadog-agent/

https://docs.datadoghq.com/agent/guide/agent-commands/?tab=agentv6v7
