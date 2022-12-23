FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY bot/ ./bot/
COPY commands/ ./commands/
COPY utils/ ./utils/

RUN go build -o caigobot-discord

EXPOSE 8080

CMD [ "/app/caigobot-discord" ]
