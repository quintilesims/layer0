FROM alpine:3.12.0
RUN apk add --no-cache ca-certificates
ADD ./l0-api /
ADD ./external /api/external/
CMD ["/l0-api"]
