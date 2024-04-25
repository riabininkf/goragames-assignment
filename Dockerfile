FROM golang:1.19-alpine as builder

WORKDIR /build
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM scratch
COPY --from=builder /build/app /app
COPY --from=builder /build/configs/config.docker.yaml /config.yaml

ENTRYPOINT ["/app"]