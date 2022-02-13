FROM golang:1.17-bullseye

RUN apt -y update && apt -y upgrade
RUN apt -y install php7.4 php7.4-common php7.4-pdo php7.4-mysql
RUN which php
RUN php -v

RUN go get github.com/cespare/reflex
COPY reflex.conf /
ENTRYPOINT ["reflex", "-c", "/reflex.conf"]