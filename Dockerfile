FROM golang:1.13.8 AS builder

COPY . /sif-example

WORKDIR /sif-example

RUN make dependencies

RUN make build

FROM scratch

COPY --from=builder /sif-example/bin/example.bin /example.bin

ENTRYPOINT [ "/example.bin" ]
