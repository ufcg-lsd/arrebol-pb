FROM golang:latest as build-stage

MAINTAINER emanueljoivo@lsd.ufcg.edu.br

WORKDIR /go/src/arrebol

COPY ./ ./

RUN go get -v

RUN go install -v

FROM debian:stretch-slim as production-stage

COPY --from=build-stage /go/bin/arrebol /usr/local/bin
COPY --from=build-stage /go/src/.env /usr/local/bin

RUN echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc

EXPOSE 5000

CMD ["arrebol"]