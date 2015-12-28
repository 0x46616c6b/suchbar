FROM       alpine:latest
MAINTAINER louis <louis@systemli.org>

ENV  GOPATH /go
ENV APPPATH $GOPATH/src/github.com/0x46616c6b/suchbar
COPY . $APPPATH
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc \
    && cd $APPPATH && go get -d && go build -o /bin/suchbar \
    && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT [ "/bin/suchbar" ]
