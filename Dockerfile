FROM golang
MAINTAINER  introcc

WORKDIR /go/src/transparent-transmission
COPY . .

CMD ["/bin/bash", "build.sh"]