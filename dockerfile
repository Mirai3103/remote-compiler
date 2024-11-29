FROM debian:bookworm-slim AS build
WORKDIR /build

RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libcap-dev \
    libsystemd-dev \
    asciidoc-base \
    unzip \
    curl \
    && rm -rf /var/lib/apt/lists/*

COPY  ./ /build/
WORKDIR /build/third_party/isolate
RUN make

# build golang
RUN curl -O https://dl.google.com/go/go1.23.3.linux-amd64.tar.gz
RUN tar xvf go1.23.3.linux-amd64.tar.gz
RUN chown -R root:root ./go
RUN mv go /usr/local
# Set go path
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/build
ENV PATH=$PATH:$GOPATH/bin

WORKDIR /build
RUN go build -o main.out main.go



FROM debian:bookworm-slim
RUN curl -sL https://deb.nodesource.com/setup_18.x | bash

RUN apt-get upgrade -y && apt-get update && apt-get install -y \
    g++ \
    gcc \
    python3 \
    openjdk-17-jdk \
    curl\
    nodejs\
    && rm -rf /var/lib/apt/lists/*
RUN npm install -g typescript

RUN curl -O https://dl.google.com/go/go1.23.3.linux-amd64.tar.gz
RUN tar xvf go1.23.3.linux-amd64.tar.gz
RUN chown -R root:root ./go
RUN mv go /usr/local
# Set go path
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/build
ENV PATH=$PATH:$GOPATH/bin
ENV PATH=$PATH:/usr/local/bin


# copy go binary
COPY --from=build /build/third_party/isolate/isolate /usr/local/bin/isolate
COPY --from=build /build/main.out /usr/local/bin/main.out
EXPOSE 8080
CMD ["main.out"]





