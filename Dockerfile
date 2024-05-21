FROM --platform=$BUILDPLATFORM golang:1.22.3-alpine AS build

ARG version
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN apk --no-cache add curl

# add Dart SASS
# ADD https://github.com/sass/dart-sass/releases/download/1.77.2/dart-sass-1.77.2-linux-$TARGETARCH-musl.tar.gz dart-sass-1.77.2-linux-$TARGETARCH-musl.tar.gz
RUN if [ "$TARGETARCH" = "amd64" ] ; then \
    curl -O -L https://github.com/sass/dart-sass/releases/download/1.77.2/dart-sass-1.77.2-linux-x64-musl.tar.gz ; else \
    curl -O -L https://github.com/sass/dart-sass/releases/download/1.77.2/dart-sass-1.77.2-linux-$TARGETARCH-musl.tar.gz ; fi

RUN if [ "$TARGETARCH" = "amd64" ] ; then \
    tar -zxvf dart-sass-1.77.2-linux-x64-musl.tar.gz ; else \
    tar -zxvf dart-sass-1.77.2-linux-$TARGETARCH-musl.tar.gz ; fi


RUN go mod download
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-X github.com/damongolding/eventsforce/cmd/version.version=${version}" -o eventsforce-docker .



FROM  zenika/alpine-chrome:latest

WORKDIR /app

COPY --from=build /app/eventsforce-docker .
COPY --from=build /app/dart-sass ./dart-sass

EXPOSE 3000

EXPOSE 35729

ENTRYPOINT ["/app/eventsforce-docker"]
