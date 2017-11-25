FROM golang:1.9

RUN mkdir -p /go/src/app

RUN apt-get update
RUN apt-get install -y netcat

COPY . /go/src/app
WORKDIR /go/src/app/cmd



RUN go-wrapper download   # "go get -d ."
RUN go-wrapper install    # "go install -v ./..."

EXPOSE 8080

CMD ["go-wrapper", "run"] # ["main"]