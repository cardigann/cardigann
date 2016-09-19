FROM scratch
COPY ./cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY definitions/ /definitions
COPY cardigann-linux-amd64 /cardigann
EXPOSE 5060
VOLUME [ "/.config/cardigann/cardigann" ]
ENTRYPOINT [ "/cardigann" ]
CMD [ "server" ]