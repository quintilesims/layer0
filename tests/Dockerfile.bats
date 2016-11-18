from alpine:3.3

RUN apk add --update \
    git \
    bash \
    jq \
    openssl \
  rm -rf /var/lib/apt/lists/*

WORKDIR /tmp
RUN git clone https://github.com/sstephenson/bats.git
RUN bats/install.sh /usr/local

ENV APP smoketest
RUN mkdir -p /tmp/$APP

WORKDIR /tmp/$APP
COPY . /tmp/$APP
RUN cp l0 /usr/local/bin

CMD [ "bash" ]
