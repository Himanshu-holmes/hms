FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest 
COPY wait-for-it.sh .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hms .

FROM alpine:latest
RUN apk --no-cache add ca-certificates netcat-openbsd 
WORKDIR /root/
COPY --from=builder /app/hms .
COPY --from=builder /go/bin/goose /usr/local/bin/goose 
RUN chmod +x /usr/local/bin/goose 
COPY --from=builder /app/wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh 
EXPOSE 3000
CMD ["./hms"]