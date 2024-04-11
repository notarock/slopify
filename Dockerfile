FROM golang:1.22-bookworm
MAINTAINER Roch Damour <roch.damour@gmail.com>

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
        ffmpeg \
        make \
        build-essential \
        automake

RUN apt-get install -y --no-install-recommends imagemagick

WORKDIR /src

CMD ["go build . -o slop"]
