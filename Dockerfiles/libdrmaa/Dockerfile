FROM drmaa/gridengine

RUN yum install -y wget tar git gcc

RUN export VERSION=1.21.2 OS=linux ARCH=amd64 && \
  wget https://go.dev/dl/go$VERSION.$OS-$ARCH.tar.gz && \
  tar -C /usr/local -xzvf go$VERSION.$OS-$ARCH.tar.gz && \
  rm go$VERSION.$OS-$ARCH.tar.gz

ENV GOPATH /go
ENV PATH /usr/local/go/bin:${PATH}:${GOPATH}/bin
ENV PATH ${PATH}:/opt/sge/include

RUN mkdir -p /go/src/github.com/dgruber/drmaa2os

ADD . /go/src/github.com/dgruber/drmaa2os
ADD ./Dockerfiles/libdrmaa/entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]

