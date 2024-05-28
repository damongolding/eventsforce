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
# RUN if [ "$TARGETARCH" = "amd64" ] ; then \
#     curl -sOL https://github.com/sass/dart-sass/releases/download/1.77.2/dart-sass-$SASS_VERSION-linux-x64-musl.tar.gz ; else \
#     curl -sOL https://github.com/sass/dart-sass/releases/download/1.77.2/dart-sass-$SASS_VERSION-linux-$TARGETARCH-musl.tar.gz ; fi

# RUN if [ "$TARGETARCH" = "amd64" ] ; then \
#     tar -zxvf dart-sass-$SASS_VERSION-linux-x64-musl.tar.gz ; else \
#     tar -zxvf dart-sass-$SASS_VERSION-linux-$TARGETARCH-musl.tar.gz ; fi


# add Tailwind
RUN curl -sOL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$TARGETOS-$TARGETARCH
RUN chmod +x tailwindcss-$TARGETOS-$TARGETARCH
RUN mv tailwindcss-$TARGETOS-$TARGETARCH tailwindcss

RUN go mod download
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-X github.com/damongolding/eventsforce/cmd/version.version=${VERSION}" -o eventsforce .



FROM  zenika/alpine-chrome:latest

WORKDIR /

COPY --from=build /app/eventsforce .
# COPY --from=build /app/dart-sass ./dart-sass
COPY --from=build /app/tailwindcss ./tailwindcss

EXPOSE 3000

EXPOSE 35729

ENTRYPOINT ["/eventsforce"]
