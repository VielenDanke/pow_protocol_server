FROM golang:latest AS builder
WORKDIR /go/src/github.com/vielendanke/pow_protocol_server
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
ENV SERVER_ADDRESS='127.0.0.1:8080'
ENV SERVER_NETWORK_TYPE='tcp'
WORKDIR /root/
COPY wisdom_words.txt ./
COPY users.txt ./
COPY --from=builder /go/src/github.com/vielendanke/pow_protocol_server/app ./
CMD ["./app"]