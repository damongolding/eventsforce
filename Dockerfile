FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o eventsforce-docker .

EXPOSE 3000

ENTRYPOINT ["./eventsforce-docker"]
