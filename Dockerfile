FROM scratch
ARG BIN=release/cardigann-linux-amd64
COPY ${BIN} /cardigann
COPY definitions/ /definitions
EXPOSE 5060
ENTRYPOINT [ "/cardigann", "server" ]