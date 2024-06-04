FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/


COPY ./static/ ./static/
COPY ./README.md ./static/
COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]
