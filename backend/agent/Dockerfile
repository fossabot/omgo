FROM golang:1.9.1
WORKDIR /go/src/
COPY . /go/src/github.com/master-g/omgo/backend/agent
RUN go install -v github.com/master-g/omgo/backend/agent
ENTRYPOINT ["/go/bin/agent", "-l", ":8888"]
EXPOSE 8888
