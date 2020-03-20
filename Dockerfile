FROM golang:1.13

WORKDIR /go/src/github.com/eagraf/cloudworker
COPY . .

RUN go get github.com/githubnemo/CompileDaemon

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build" -command="./cloudworker"

#CMD ["cloudworker"]