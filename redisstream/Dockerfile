FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN cd redisstream && go build -o main

CMD ["./redisstream/main"]

