FROM golang:1.9-alpine
RUN apk add --update ca-certificates
WORKDIR /go/src/github.com/cardigann/cardigann
COPY . /go/src/github.com/cardigann/cardigann
RUN go build -o /bin/cardigann
EXPOSE 5060
ENV CONFIG_DIR=/.config/cardigann
ENTRYPOINT [ "/bin/cardigann" ]
CMD [ "server" ]
