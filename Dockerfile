FROM alpine:latest

ARG CARDIGANN_BIN=cardigann-linux-amd64
ENV CONFIG_DIR=/etc/cardigann

COPY ${CARDIGANN_BIN} /bin/cardigann

RUN \
    # install ca certs
    apk update && apk add ca-certificates && \
    # make sure the executable can see the shared library
    mkdir /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 5060
VOLUME [ "/etc/cardigann" ]
ENTRYPOINT [ "/bin/cardigann" ]
CMD [ "server" ]
