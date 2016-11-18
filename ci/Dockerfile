FROM gitlab/gitlab-runner:alpine-v1.4.1

ADD entrypoint /
RUN chmod +x /entrypoint

RUN apk --update add --no-cache git make && rm -rf /var/cache/apk/*

ADD config.json /root/.docker/config.json
