FROM golang:1.22.3-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o eventsforce-docker .

#FROM alpine:latest
#COPY --from=build /app/eventsforce-docker .

EXPOSE 3000
EXPOSE 35729
ENTRYPOINT ["./eventsforce-docker"]
