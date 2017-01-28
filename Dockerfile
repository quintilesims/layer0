FROM golang:1.7.4
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y \
	wget \
	curl \ 
	zip \ 
	python \
    mysql-server \
    jq \
	python-pip

RUN curl -sSL https://get.docker.com/ | sh
RUN pip install awscli

ENV APP github.com/quintilesims/layer0
RUN mkdir -p /go/src/$APP
WORKDIR /go/src/$APP
ENTRYPOINT [ "./scripts/entrypoint.sh" ]

COPY . /go/src/$APP/
