FROM golang:1.9
WORKDIR /go/src/
COPY . /go/src/github.com/master-g/omgo/backend/game
RUN go install -v github.com/master-g/omgo/backend/game
ENTRYPOINT ["/go/bin/game", "-l", ":10000"]
EXPOSE 10000
