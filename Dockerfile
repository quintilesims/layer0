FROM golang:1.7.4
ENV DEBIAN_FRONTEND noninteractive

ENV LAYER0_PREFIX=$LAYER0_PREFIX
ENV TF_VAR_key_pair=$TF_VAR_key_pair
ENV TF_VAR_access_key=$TF_VAR_access_key
ENV TF_VAR_secret_key=$TF_VAR_secret_key
ENV AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
ENV AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY
ENV DOCKERHUB_USERNAME=$DOCKERHUB_USERNAME
ENV DOCKERHUB_PASSWORD=$DOCKERHUB_PASSWORD

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
ENV DEBIAN_FRONTEND interactive
ENTRYPOINT [ "./scripts/entrypoint.sh" ]

COPY . /go/src/$APP/
