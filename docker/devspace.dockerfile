FROM golang:1.20-alpine AS dev
RUN apk update && apk upgrade && apk add --no-cache git
WORKDIR /app
COPY . .
ENTRYPOINT ["go", "run", "main.go"]
