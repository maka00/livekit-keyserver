FROM golang:1.23-alpine AS builder

COPY cmd /opt/app/cmd
COPY internal /opt/app/internal
COPY go.mod /opt/app/go.mod
COPY go.sum /opt/app/go.sum
COPY main.go /opt/app/main.go

WORKDIR /opt/app

RUN go mod download
RUN go build -o /opt/app/bin/server .

FROM alpine

COPY --from=builder /opt/app/bin/server /opt/app/server

CMD ["/opt/app/server", "serveToken"]