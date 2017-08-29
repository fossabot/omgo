FROM golang:1.9
WORKDIR /go/src/
COPY . /go/src/github.com/master-g/omgo/backend/reception
RUN go install -v github.com/master-g/omgo/backend/reception
ENTRYPOINT ["/go/bin/reception", "-l", ":8080"]
EXPOSE 8080
