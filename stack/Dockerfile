FROM golang:latest as build

MAINTAINER emanueljoivo@lsd.ufcg.edu.br

WORKDIR /go/src/arrebol

COPY ./ ./

RUN go get -v

RUN go install -v

FROM debian:stretch-slim as production

COPY --from=build /go/bin/arrebol /usr/local/bin

RUN echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc

EXPOSE 5000

CMD ["arrebol"]