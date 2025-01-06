FROM golang:1.23.4-alpine3.21

RUN apk update && apk add --no-cache iputils curl

RUN go install github.com/air-verse/air@latest

WORKDIR /workspaces

COPY . .

CMD ["go", "run", "./app/main.go"]