FROM alpine:3.4
RUN apk add --update ca-certificates
COPY cardigann-linux-amd64 /cardigann
EXPOSE 5060
ENV CONFIG_DIR=/.config/cardigann
VOLUME [ "/.config/cardigann" ]
ENTRYPOINT [ "/cardigann" ]
CMD [ "server" ]