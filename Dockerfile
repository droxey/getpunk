FROM golang:1.14

RUN mkdir -p $HOME/go/src
COPY . $HOME/go/src
WORKDIR $HOME/go/src
RUN go mod download
EXPOSE 8888
RUN go build
CMD ["./getpunk"]
