FROM iron/go:dev

WORKDIR /app

ENV SRC_DIR=/go/src/github.com/syn-inc/server/

ADD . $SRC_DIR

EXPOSE 8000

RUN cd $SRC_DIR; go get github.com/tools/godep; godep restore; go build -o server; cp server /app/

ENTRYPOINT ["./server"]