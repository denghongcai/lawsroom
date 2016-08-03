FROM alpine:latest

ADD law /law
ADD public/index.html /public/index.html

ENTRYPOINT ["/law"]

