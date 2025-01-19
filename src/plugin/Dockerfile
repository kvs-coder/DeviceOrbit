FROM golang:1.23.6-alpine AS build

RUN apk update && apk upgrade && apk add tar ca-certificates build-base libusb-dev pkgconfig

ENV GOPATH=/go
RUN go version

WORKDIR /build
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \ 
    go mod download && \
    go build -ldflags='-s -w -extldflags="-static"' -o mobile-device-plugin && go test

    
FROM alpine:latest AS runtime

LABEL service.name="mobile-device-plugin"

WORKDIR /root

COPY --from=build /build/mobile-device-plugin /usr/bin/mobile-device-plugin

CMD ["/usr/bin/mobile-device-plugin","-logtostderr=true","-v=0"]
