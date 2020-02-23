FROM scratch

ADD bin/example.bin /example.bin

ENTRYPOINT [ "/example.bin" ]
