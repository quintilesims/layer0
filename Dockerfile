FROM golang:1.7.4

RUN apt-get update && apt-get install -y \
	wget \
	curl \ 
	zip \ 
	python \
	python-pip

RUN curl -sSL https://get.docker.com/ | sed 's/docker-engine/docker-engine=1.9.1-0~jessie/' |  sh
RUN pip install awscli

ENV APP github.com/quintilesims/layer0
RUN mkdir -p /go/src/$APP
WORKDIR /go/src/$APP
ENTRYPOINT [ "./scripts/entrypoint.sh" ]

COPY . /go/src/$APP/
