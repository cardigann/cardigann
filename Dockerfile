FROM scratch
ARG BIN=release/cardigann-linux-amd64
ADD cacert.pem https://curl.haxx.se/ca/cacert.pem
COPY ${BIN} /cardigann
COPY definitions/ /definitions
EXPOSE 5060
VOLUME [ "/config.json" ]
ENTRYPOINT [ "/cardigann" ]
CMD [ "server" ]