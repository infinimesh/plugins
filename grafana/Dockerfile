FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN go build -o main .

CMD ["./main"]

