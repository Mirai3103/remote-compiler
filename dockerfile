FROM debian:bookworm-slim AS build
WORKDIR /build

# To compile Isolate, you need:

#     - pkg-config
  
#     - headers for the libcap library (usually available in a libcap-dev package)
  
#     - headers for the libsystemd library (libsystemd-dev package) for compilation
#       of isolate-cg-keeper
# Cài  đặt các gói cần thiết
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libcap-dev \
    libsystemd-dev \
    asciidoc-base \
    unzip \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Tải về mã nguồn của Isolate và build
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
COPY --from=build /build/third_party/isolate/isolate /usr/local/bin/isolate

RUN apt-get update && apt-get install -y \
    g++ 

RUN apt-get install -y \
    g++ \
    gcc \
    python3 \
    openjdk-17-jdk \
    curl\
    && rm -rf /var/lib/apt/lists/*
RUN curl -sL https://deb.nodesource.com/setup_18.x | bash
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y nodejs
RUN npm install -g typescript

#  cài golang 
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
COPY --from=build /build/main.out /usr/local/bin/main.out
EXPOSE 8080
CMD ["main.out"]





