FROM --platform=$BUILDPLATFORM golang:1.22.3-alpine AS build

ARG VERSION
ARG TARGETOS
ARG TARGETARCH
ARG TAILWIND_VERSION

WORKDIR /app

COPY . .

# Install curl
RUN apk --no-cache add curl

# add Tailwind
RUN curl -sOL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$TARGETOS-$TARGETARCH
RUN chmod +x tailwindcss-$TARGETOS-$TARGETARCH
RUN mv tailwindcss-$TARGETOS-$TARGETARCH tailwindcss

RUN go mod download
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-X github.com/damongolding/eventsforce/cmd/version.version=${VERSION} -X github.com/damongolding/eventsforce/cmd/version.tailwindVersion=${TAILWIND_VERSION}" -o eventsforce .



FROM  zenika/alpine-chrome:latest

WORKDIR /

COPY --from=build /app/eventsforce .
COPY --from=build /app/tailwindcss ./tailwindcss

EXPOSE 3000

EXPOSE 35729

ENTRYPOINT ["/eventsforce"]
