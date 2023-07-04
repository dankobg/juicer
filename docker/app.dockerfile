FROM golang:1.20-alpine AS dev
RUN apk update && apk upgrade && apk add --no-cache git dpkg gcc musl-dev
RUN go install github.com/cespare/reflex@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -gcflags="all=-N -l" -o juicer_dev .
RUN CGO_ENABLED=0 GOOS=linux go build -o juicer .
CMD ["reflex", "-d", "none", "-r", "\\.go$$|go\\.mod", "-s", "--", "sh", "-c", "go run ./"]



FROM alpine:latest AS debug
RUN apk update && apk upgrade && apk add --no-cache git
WORKDIR /dbg
COPY --from=dev /app .
COPY --from=dev /go/bin/dlv .
CMD ["./dlv", "--api-version=2", "--listen=:2345", "--headless", "--accept-multiclient", "--log", "exec", "./juicer_dev"]



FROM alpine:latest AS prod
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
WORKDIR /final
COPY --from=dev /app/juicer .
CMD ["./juicer"]
