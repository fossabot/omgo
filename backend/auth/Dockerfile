FROM golang:1.8
WORKDIR /go/src/
COPY . /go/src/github.com/master-g/omgo/backend/auth
RUN go install -v github.com/master-g/omgo/backend/auth
ENTRYPOINT ["/go/bin/auth", "-l", ":40000", "-p", ":8080"]
EXPOSE 8080 40000
