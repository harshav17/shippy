# user-cli/Dockerfile
FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app/user-cli

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o user-cli -a -installsuffix cgo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/user-cli .

ENTRYPOINT ["./user-cli"]

CMD ["./user-cli"]