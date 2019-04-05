FROM golang:latest

WORKDIR /app

ENV SRC_DIR=/go/src/github.com/syn-inc/server/

ADD . $SRC_DIR

EXPOSE 8000

RUN cd $SRC_DIR; go build -o server; cp server /app/

ENTRYPOINT ["./server"]