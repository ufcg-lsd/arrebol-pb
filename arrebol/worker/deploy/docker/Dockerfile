FROM golang:latest

WORKDIR /go/src/github.com/ufcg-lsd/arrebol-pb/arrebol

RUN echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
RUN apt install jq

COPY ../.. .
COPY ../../../../go.mod ./worker
COPY ../../../../go.sum ./worker

RUN go get -v
RUN go install -v

RUN go build -o main ./worker

EXPOSE 8080

CMD ["./main"]
