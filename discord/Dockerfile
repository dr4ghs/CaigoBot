FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY bot/ ./bot/
COPY configs/ ./configs/
COPY config.yml ./

RUN go build -o /caigobot-discord

EXPOSE 8080

CMD [ "/caigobot-discord" ]
