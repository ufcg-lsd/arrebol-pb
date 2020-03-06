FROM golang:latest as build

MAINTAINER emanueljoivo@lsd.ufcg.edu.br

WORKDIR /go/src/arrebol

COPY ./ ./

RUN go get -v

RUN go install -v

FROM debian:stretch-slim as production

COPY --from=build /go/bin/arrebol /usr/local/bin

ARG DB_ADDR
ARG DB_USER
ARG DB_PASSWD
ARG DB_NAME
ARG DB_PORT
ARG WORKER_AMOUNT

ENV DATABASE_ADDRESS=$DB_ADDR
ENV DATABASE_USER=$DB_USER
ENV DATABASE_PASSWORD=$DB_PASSWD
ENV DATABASE_PORT=$DB_PORT
ENV DATABASE_NAME=$DB_NAME
ENV STATIC_WORKER_POOL=$WORKER_AMOUNT

RUN echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc

EXPOSE 5000

CMD ["arrebol"]