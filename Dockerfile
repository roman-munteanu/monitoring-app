FROM golang:1.19.2

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /monitoring-app

EXPOSE 3001

CMD [ "/monitoring-app" ]
