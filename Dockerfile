FROM ubuntu:20.04

COPY . /app

RUN apt update && apt install -y \
 build-essential \
 git \
 curl \
 wget \
 gcc

RUN /app/install.sh

RUN wget -P /tmp https://dl.google.com/go/go1.19.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf /tmp/go1.19.linux-amd64.tar.gz
RUN rm /tmp/go1.19.linux-amd64.tar.gz

RUN apt install -y libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev liblz4-dev libzstd-dev

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR "/app" 

RUN go mod download

RUN go build -tags builtin_static ./cmd/http