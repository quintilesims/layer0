FROM nginx:1.10-alpine

RUN apk add --update \
    python \
    py-pip \
    make

RUN rm -rf /var/lib/apt/lists/*

WORKDIR /data

ADD ./Makefile .
RUN make deps

ADD ./nginx.default /usr/local/nginx/conf/sites-enabled/default
ADD ./nginx.conf /etc/nginx/nginx.conf

ADD ./run.sh .
RUN chmod +x ./run.sh

EXPOSE 80
CMD ["./run.sh"]
