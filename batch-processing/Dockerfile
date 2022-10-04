# syntax=docker/dockerfile:1
##
## BASE
##
FROM golang:1.19-bullseye as base

# Ignore APT warnings about not having a TTY
ENV DEBIAN_FRONTEND noninteractive

# install build essentials
RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

# Install ImageMagick deps
RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev \
    libgif-dev libx11-dev --no-install-recommends

ENV IMAGEMAGICK_VERSION=6.9.10-11

# Install ImageMagick
RUN cd && \
    wget https://github.com/ImageMagick/ImageMagick6/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
    tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
    cd ImageMagick* && \
    ./configure \
    --without-magick-plus-plus \
    --without-perl \
    --disable-openmp \
    --with-gvc=no \
    --disable-docs && \
    make -j$(nproc) && make install && \
    ldconfig /usr/local/lib

# Build the app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
RUN mkdir -p /inputs /outputs

# This is required for test and run, but for develop it ensures we have a build cache
RUN go build -o /out

# Set up environment
ENV AWS_REGION ""
ENV AWS_ROLE_ARN ""
ENV S3_BUCKET ""

##
## TEST
##
FROM base as test

ENTRYPOINT [ "go", "test", "-v" ]

##
## DEVELOP
##
FROM base as develop

RUN mkdir -p /root/.aws /root/.cache/go-build
COPY .bash_history /root/.bash_history
ENTRYPOINT [ "/bin/bash" ]

##
## RUN
##
FROM base as run

WORKDIR /
ENTRYPOINT ["/out"]
CMD ["--input",  "/inputs/example.csv",  "--output",  "/outputs/example-result.csv"]