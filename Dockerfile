FROM alpine:3.20

RUN apk add --no-cache git ca-certificates

COPY bin/eval-diff /usr/local/bin/eval-diff
RUN chmod +x /usr/local/bin/eval-diff

ENTRYPOINT ["/usr/local/bin/eval-diff"]

