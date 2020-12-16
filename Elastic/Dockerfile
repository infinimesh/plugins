FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN cd Elastic && go build -o main

CMD ["./Elastic/main"]

