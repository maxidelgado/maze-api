FROM golang:1.14-stretch as builder
WORKDIR /github.com/maxidelgado/maze-api
COPY . .
RUN go mod vendor
RUN go test ./...
RUN CGO_ENABLED=0 go build -o app ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /github.com/maxidelgado/maze-api/app .
CMD ["./app"]
