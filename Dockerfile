FROM alpine:latest

ENV CARDIGANN_PKG cardigann-stable-linux-amd64.tgz
ENV CONFIG_DIR /etc/cardigann

ADD https://bin.equinox.io/c/3u8U4iwUn6o/${CARDIGANN_PKG} /

# install ca certs
RUN apk update && apk add ca-certificates
# make sure the executable can see the shared library
RUN mkdir /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# unpack
RUN tar xfz /${CARDIGANN_PKG} -C /bin/
# cleanup
RUN rm -fr /${CARDIGANN_PKG} /root/.cache /tmp/*

EXPOSE 5060
VOLUME [ "/etc/cardigann" ]
ENTRYPOINT [ "/bin/cardigann" ]
CMD [ "server" ]
