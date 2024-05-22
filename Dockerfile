FROM --platform=$BUILDPLATFORM golang:1.22.3-alpine AS build

ARG VERSION
ARG TARGETOS
ARG TARGETARCH
ARG SASS_VERSION

WORKDIR /app

COPY . .

# Install curl
RUN apk --no-cache add curl

# add Dart SASS
RUN if [ "$TARGETARCH" = "amd64" ] ; then \
    curl -O -L https://github.com/sass/dart-sass/releases/download/1.77.2/dart-sass-$SASS_VERSION-linux-x64-musl.tar.gz ; else \
    curl -O -L https://github.com/sass/dart-sass/releases/download/1.77.2/dart-sass-$SASS_VERSION-linux-$TARGETARCH-musl.tar.gz ; fi

RUN if [ "$TARGETARCH" = "amd64" ] ; then \
    tar -zxvf dart-sass-$SASS_VERSION-linux-x64-musl.tar.gz ; else \
    tar -zxvf dart-sass-$SASS_VERSION-linux-$TARGETARCH-musl.tar.gz ; fi


RUN go mod download
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-X github.com/damongolding/eventsforce/cmd/version.version=${VERSION}" -o eventsforce .



FROM  zenika/alpine-chrome:latest

WORKDIR /

COPY --from=build /app/eventsforce .
COPY --from=build /app/dart-sass ./dart-sass

EXPOSE 3000

EXPOSE 35729

ENTRYPOINT ["/eventsforce"]
