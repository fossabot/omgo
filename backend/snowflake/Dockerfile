FROM golang:1.8
MAINTAINER mg <masterg@yeah.net>
WORKDIR /go/src/
COPY . /go/src/github.com/master-g/omgo/backend/snowflake
RUN go install -v github.com/master-g/omgo/backend/snowflake
ENTRYPOINT ["/go/bin/snowflake", "-p", "40001", "-e", "http://127.0.0.1:2379"]
EXPOSE 40001
