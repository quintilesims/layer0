FROM alpine
RUN apk add --no-cache ca-certificates
ADD ./static /static
ADD ./l0-api /
CMD ["/l0-api"]
