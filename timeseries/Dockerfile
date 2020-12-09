FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN cd timeseries && go build -o main

CMD ["./timeseries/main"]

