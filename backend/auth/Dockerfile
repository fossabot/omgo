FROM golang:1.8
MAINTAINER mg <masterg@yeah.net>
WORKDIR /go/src/
COPY . /go/src/github.com/master-g/omgo/backend/auth
RUN go install -v github.com/master-g/omgo/backend/auth
ENTRYPOINT ["/go/bin/auth", "-l", ":40000"]
EXPOSE 40000
