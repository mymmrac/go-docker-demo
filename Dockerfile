# Source code + build dependencies base image
FROM golang:1.20-alpine AS source

WORKDIR /demo

RUN go env -w CGO_ENABLED="0"

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Dev image with live reloading + debugger
FROM source AS dev

RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/cosmtrek/air@latest

ENTRYPOINT air

# Test image
FROM source AS test

RUN go test ./...

# Build image + compress binary
FROM source AS build

RUN apk --update add ca-certificates upx && update-ca-certificates

RUN go build -ldflags="-s -w" -o /bin/demo . && upx --best --lzma /bin/demo

# Release image + run dependencies + helath check
FROM scratch AS release

COPY --from=mymmrac/mini-health:latest /mini-health /mini-health
HEALTHCHECK CMD ["/mini-health", "/health"]

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/demo /demo

ENTRYPOINT ["/demo"]
