# syntax=docker/dockerfile:1
##
## IMAGEMAGICK
##
FROM golang:1.22-bullseye as imagemagick

# Ignore APT warnings about not having a TTY
ENV DEBIAN_FRONTEND noninteractive

# Install build essentials
RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

# Install ImageMagick deps
RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev \
    libgif-dev libx11-dev --no-install-recommends

ENV IMAGEMAGICK_VERSION=6.9.10-11

# Install ImageMagick
RUN cd && \
    wget https://github.com/ImageMagick/ImageMagick6/archive/${IMAGEMAGICK_VERSION}.tar.gz

RUN cd && \
    tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
    cd ImageMagick* && \
    ./configure \
    --without-magick-plus-plus \
    --without-perl \
    --disable-openmp \
    --with-gvc=no \
    --disable-docs

RUN cd && cd ImageMagick* && \
    make -j$(nproc) && make install && \
    ldconfig /usr/local/lib

# DEBUG
# install bash so i can examine the filesystem
# RUN apt-get update && apt-get install -y bash
# docker build --target imagemagick -t imagemagick . && docker run -it --rm imagemagick /bin/bash

##
## BUILD
##
FROM golang:1.22-bullseye as build

COPY --from=imagemagick /usr/local/bin /usr/local/bin
COPY --from=imagemagick /usr/local/lib /usr/local/lib
COPY --from=imagemagick /usr/lib/x86_64-linux-gnu/ /usr/lib/x86_64-linux-gnu/
COPY --from=imagemagick /usr/local/include/ImageMagick-6 /usr/local/include/ImageMagick-6

# Build the app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

COPY inputs /inputs
RUN mkdir -p /outputs

# This is required for test and run, but for develop it ensures we have a build cache
RUN go build -o /out

# DEBUG
# install bash so i can examine the filesystem
# RUN apt-get update && apt-get install -y bash
# docker build --target gobuild -t gobuild . && docker run -it --rm gobuild /bin/bash

##
## COMBINED
##
# FROM golang:1.22-bullseye as combined

# COPY --from=gobuild /out /out

# COPY --from=imagemagick /usr/local/bin /usr/local/bin
# COPY --from=imagemagick /usr/local/lib /usr/local/lib
# COPY --from=imagemagick /usr/lib/x86_64-linux-gnu/ /usr/lib/x86_64-linux-gnu/
# COPY --from=imagemagick /usr/local/include/ImageMagick-6 /usr/local/include/ImageMagick-6

# this doesnt work because the go build stage needs access to the libraries

##
## TEST
##
FROM build as test

ENTRYPOINT [ "go", "test", "-v" ]

##
## DEVELOP
##
FROM build as develop

COPY .bash_history /root/.bash_history
ENTRYPOINT [ "/bin/bash" ]

##
## RUN
##
FROM build as run

WORKDIR /
ENTRYPOINT /out --input /inputs/gradient.jpg --output /outputs/gradient_bw.jpg