FROM golang:1.14 as go

# The port your service will listen on
EXPOSE 8080

# The command to run
CMD ["/roteiro"]

ARG BUILD_TAG=unknown
LABEL BUILD_TAG=$BUILD_TAG


ENV ROTEIRO_SERVER_PORT 8080
ENV ROTEIRO_SERVER_MODE release

COPY bin/roteiro /roteiro
