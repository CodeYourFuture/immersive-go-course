# syntax=docker/dockerfile:1
##
## IMAGEMAGICK
##
FROM golang:1.19-bullseye as imagemagick

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

# Clean up the image
RUN cd && rm -rf ${IMAGEMAGICK_VERSION}.tar.gz ImageMagick6-${IMAGEMAGICK_VERSION}

##
## GOBUILD
##
FROM golang:1.19-bullseye as gobuild

# Build the app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY inputs /inputs
RUN mkdir -p /outputs
COPY *.go ./

# This is required for test and run, but for develop it ensures we have a build cache
RUN go build -o /out

##
## COMBINED
##
FROM golang:1.19-bullseye as combined

# copy ImageMagick
COPY --from=imagemagick /usr/local /usr/local

# copy the go app
COPY --from=gobuild /out /out
COPY --from=gobuild /inputs /inputs
COPY --from=gobuild /outputs /outputs

##
## TEST
##
FROM combined as test

ENTRYPOINT [ "go", "test", "-v" ]

##
## DEVELOP
##
FROM combined as develop

COPY .bash_history /root/.bash_history
ENTRYPOINT [ "/bin/bash" ]

##
## RUN
##
FROM combined as run

WORKDIR /
ENTRYPOINT /out --input /inputs/gradient.jpg --output /outputs/gradient_bw.jpg