FROM golang:1.22-alpine AS dev
RUN apk update && apk upgrade && apk add --no-cache git
WORKDIR /app
COPY . .
ENTRYPOINT ["go", "run", "main.go"]
