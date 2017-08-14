FROM golang:1.8
WORKDIR /go/src/
COPY . /go/src/github.com/master-g/omgo/backend/db/dbservice
RUN go install -v github.com/master-g/omgo/backend/db/dbservice
EXPOSE 60001
ENTRYPOINT ["/go/bin/dbservice", "-l", "60001", "-r", "127.0.0.1:6379", "-m", "127.0.0.1:37017", "-b", "master", "-u", "admin", "-p", "admin"]
