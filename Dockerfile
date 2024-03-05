# syntax = docker/dockerfile:experimental

ARG GOLANG_VER=1.20.8
ARG ALPINE_VER=3.18.6

# Builder image
FROM golang:${GOLANG_VER} as builder

ENV GOCACHE=/tmp/
ENV GOMODCACHE=/tmp/
ENV CGO_ENABLED=0

ADD . /app
WORKDIR /app
RUN make build


# actual image
FROM alpine:${ALPINE_VER}
ARG BINARY_NAME=dns-proxy

# ARG can't be used in CMD, but ENV can
ENV BINARY_NAME=${BINARY_NAME}

# add ca-certificates
# required for establishing SSL connections
RUN set -ex \
    && apk --no-cache add ca-certificates tzdata

ENV TZ=Europe/Berlin

# add a custom user for service execution
RUN set -ex \
    && addgroup -S service \
    && adduser -S service -G service

# add binaries
COPY --from=builder "/app/${BINARY_NAME}" /usr/local/bin/

# symlink binary for convenience
RUN set -ex \
    && if [ "${BINARY_NAME}" != "service" ]; then \
        ln -sv "/usr/local/bin/${BINARY_NAME}" /usr/local/bin/service; \
    fi

USER service:service

CMD ["sh","-c","exec /usr/local/bin/${BINARY_NAME} -tcp"]
