FROM golang:1.23.1 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

WORKDIR /app/pkg

COPY pkg/ ./

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/my_app

FROM alpine:latest AS production

WORKDIR /app/web

COPY web/ ./

WORKDIR /app

COPY --from=build /app/my_app .

CMD ["/app/my_app"]