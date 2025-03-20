FROM buildpack-deps:stable

RUN set -xe && \
    apt-get update && \
    apt-get install -y --no-install-recommends git libcap-dev && \
    rm -rf /var/lib/apt/lists/* && \
    git clone https://github.com/judge0/isolate.git /tmp/isolate && \
    cd /tmp/isolate && \
    make -j$(nproc) install && \
    rm -rf /tmp/*

RUN set -xe && \
    curl -OL https://golang.org/dl/go1.22.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz && \
    rm go1.22.1.linux-amd64.tar.gz 

ENV PATH=$PATH:/usr/local/go/bin:/build/bin:/usr/local/bin \
    GOPATH=/build

WORKDIR /build
COPY . .
RUN go build -o /usr/local/bin/app.out cmd/server/main.go
RUN mkdir /isolateBox
RUN mkdir ./temp

EXPOSE 8080
EXPOSE 50051

CMD ["app.out"]