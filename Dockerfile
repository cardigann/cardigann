FROM scratch
COPY release/cardigann-linux-amd64 /cardigann
COPY definitions/ /definitions
EXPOSE 5060
ENTRYPOINT [ "/cardigann", "server" ]