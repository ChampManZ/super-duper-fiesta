FROM golang:1.22.6-alpine3.20

RUN apk add --no-cache bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x wait-for-it.sh

EXPOSE 1323

CMD ["./wait-for-it.sh", "db:3306", "--", "go", "run", "main.go"]