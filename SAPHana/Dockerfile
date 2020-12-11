FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN cd SAPHana && go build -o main

CMD ["./SAPHana/main"]

