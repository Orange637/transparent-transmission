FROM golang
MAINTAINER  introcc

ENV PORT 5555

WORKDIR /go/src/transparent-transmission
COPY . .

RUN go build .

ENTRYPOINT  ["./transparent-transmission"]

#CMD ["/bin/bash", "build.sh"]