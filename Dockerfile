FROM alpine:latest

ADD law /law

ENTRYPOINT ["/law"]

