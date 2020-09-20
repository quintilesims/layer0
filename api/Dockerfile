FROM alpine
RUN apk add --no-cache ca-certificates
ADD ./l0-api /
ADD ./external /api/external/
CMD ["/l0-api"]
