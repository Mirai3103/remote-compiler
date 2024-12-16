FROM buildpack-deps:stable
WORKDIR /build
RUN set -xe && \
    sudo apt-get update && \
    sudo apt-get install -y --no-install-recommends git libcap-dev && \
    sudo rm -rf /var/lib/apt/lists/* && \
    sudo git clone https://github.com/judge0/isolate.git /tmp/isolate && \
    cd /tmp/isolate && \
    sudo make -j$(nproc) install && \
    sudo rm -rf /tmp/*
RUN set -xe &&\
    sudo curl -OL https://golang.org/dl/go1.23.4.linux-amd64.tar.gz &&\
    sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz &&\
    sudo rm go1.23.4.linux-amd64.tar.gz 
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/build
ENV PATH=$PATH:$GOPATH/bin
ENV PATH=$PATH:/usr/local/bin
WORKDIR /build
COPY . .
RUN go build -o main.out cmd/server/main.go 
COPY  main.out /usr/local/bin/app.out
EXPOSE 8080
CMD ["app.out"]





