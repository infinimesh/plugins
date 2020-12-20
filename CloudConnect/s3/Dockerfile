FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN cd CloudConnect/s3 && go build -o main

CMD ["./CloudConnect/s3/main"]

