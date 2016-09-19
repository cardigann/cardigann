FROM scratch
ARG BIN=release/cardigann-linux-amd64
ADD ./ca-cert.pem /etc/ssl/certs/ca-certificates.crt
COPY ${BIN} /cardigann
COPY definitions/ /definitions
EXPOSE 5060
VOLUME [ "/config.json" ]
ENTRYPOINT [ "/cardigann" ]
CMD [ "server" ]