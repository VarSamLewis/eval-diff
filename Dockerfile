FROM debian:bookworm-slim

RUN apk add --no-cache ca-certificates git 2>/dev/null || \
    (apt-get update && apt-get install -y --no-install-recommends ca-certificates git && rm -rf /var/lib/apt/lists/*)

COPY bin/eval-diff /usr/local/bin/eval-diff
RUN chmod +x /usr/local/bin/eval-diff

ENTRYPOINT ["/usr/local/bin/eval-diff"]
