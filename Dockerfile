FROM golang:1.20-alpine AS source

WORKDIR /demo

RUN go env -w CGO_ENABLED="0"

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

FROM source AS dev

RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/cespare/reflex@latest

CMD reflex --decoration="none" -R "bin/" -s -- sh -c \
    "dlv debug --output ./bin/demo --headless --continue --accept-multiclient --listen :2345 --api-version=2 --log ./"

FROM source AS test

RUN go test ./...

FROM source AS build

RUN apk --update add ca-certificates upx && update-ca-certificates

RUN go build -ldflags="-s -w" -o /bin/demo . && upx --best --lzma /bin/demo

FROM scratch AS release

COPY --from=mymmrac/mini-health:latest /mini-health /mini-health
HEALTHCHECK CMD ["/mini-health", "/health"]

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/demo /demo

ENTRYPOINT ["/demo"]
