FROM golang:1.23.6-alpine AS build

RUN apk add --no-cache ca-certificates

ENV GOPATH=/go
RUN go version

WORKDIR /build
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \ 
    go mod download && \
    CGO_ENABLED=0 go build -ldflags='-s -w' -o mobile-device-controller .


FROM alpine:latest AS runtime

RUN apk add --no-cache curl ca-certificates

WORKDIR /root

COPY --from=build /build/mobile-device-controller /usr/bin/mobile-device-controller

CMD ["/usr/bin/mobile-device-controller", "-logtostderr=true", "-v=0"]

