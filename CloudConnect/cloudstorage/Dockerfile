FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN cd CloudConnect/cloudstorage && go build -o main

CMD ["./CloudConnect/cloudstorage/main"]

