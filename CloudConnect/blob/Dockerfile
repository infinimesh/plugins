FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN cd CloudConnect/blob && go build -o main

CMD ["./CloudConnect/blob/main"]

