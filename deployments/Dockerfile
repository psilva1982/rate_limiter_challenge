FROM golang:1.21-alpine

RUN apk --no-cache add curl


WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
COPY .env /app/cmd/server/

EXPOSE 8080

WORKDIR /app/cmd/server/

RUN go build -o run

CMD ["/app/cmd/server/run"]