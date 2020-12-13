FROM golang:1.15.6-alpine

WORKDIR /go/src/app
COPY . .

RUN cd Snowflake && go build -o main

CMD ["./Snowflake/main"]

