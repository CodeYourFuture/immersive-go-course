FROM golang:1.22-bookworm as imagemagickbase

# Work out how to lock the version of libmagickwand-dev
RUN apt-get update && apt-get install -y libmagickwand-dev

FROM golang:1.22-bookworm as appbase

WORKDIR /batch-processing

COPY go.mod go.sum ./

RUN go mod download

COPY inputs inputs

RUN mkdir -p outputs

COPY *.go .

FROM golang:1.22-bookworm as combined

# Work out how to minimise how many libraries (and headers(?)) we are copying across
COPY --from=imagemagickbase /usr/lib/ /usr/lib/
COPY --from=imagemagickbase /usr/include /usr/include
COPY --from=appbase /batch-processing/ /batch-processing/

# root@917b05018df9:/usr/lib/aarch64-linux-gnu# pkg-config --cflags --libs MagickWand
# -fopenmp -DMAGICKCORE_HDRI_ENABLE=0 -DMAGICKCORE_QUANTUM_DEPTH=16 -fopenmp -DMAGICKCORE_HDRI_ENABLE=0 -DMAGICKCORE_QUANTUM_DEPTH=16 -I/usr/include/aarch64-linux-gnu//ImageMagick-6 -I/usr/include/ImageMagick-6 -I/usr/include/aarch64-linux-gnu//ImageMagick-6 -I/usr/include/ImageMagick-6 -lMagickWand-6.Q16 -lMagickCore-6.Q16
# root@917b05018df9:/usr/lib/aarch64-linux-gnu# 

WORKDIR /batch-processing

RUN go build -o out .

FROM combined as develop

ENTRYPOINT [ "/bin/bash" ]

FROM combined as test

ENTRYPOINT [ "go", "test", "-v" ]

FROM combined as run

RUN ls -la

ENTRYPOINT ["./out", "--input=inputs/gradient.jpg", "--output=outputs/gradient_bw.jpg"]