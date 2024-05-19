FROM golang:latest

RUN useradd -u 1001 -m iamuser

WORKDIR /app

COPY . .

USER 1001

RUN go mod download

RUN go build -o eventsforce-docker .

EXPOSE 3000

ENTRYPOINT ["./eventsforce-docker"]
