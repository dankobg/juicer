ARG GO_VERSION=1.22-alpine
FROM golang:${GO_VERSION} AS dev
RUN apk update && apk upgrade \ 
  && apk add --no-cache git dpkg gcc musl-dev \
  && go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "main.go"]
