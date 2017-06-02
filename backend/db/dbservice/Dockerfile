FROM golang:1.8.1
WORKDIR /go/src/
COPY . /go/src/github.com/master-g/omgo/backend/db/dbservice
RUN go install -v github.com/master-g/omgo/backend/db/dbservice
EXPOSE 60001
ENTRYPOINT ["/go/bin/dbservice", "-p", "60001", "-e", "http://127.0.0.1:2379"]
