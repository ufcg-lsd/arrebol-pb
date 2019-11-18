FROM debian:stretch-slim

MAINTAINER emanueljoivo@lsd.ufcg.edu.br

ARG OBJECT_FILE

WORKDIR /service

EXPOSE 8080

COPY $OBJECT_FILE /service

RUN ["/bin/sh", "-c", "chmod +x $OBJECT_FILE"]

ENTRYPOINT ["./$OBJECT_FILE"]