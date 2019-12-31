# consignment-cli/Dockerfile
FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app/consignment-cli

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o consignment-cli -a -installsuffix cgo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/consignment-cli .

CMD ["./consignment-cli"]