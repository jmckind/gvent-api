# Create the builder image
FROM golang:alpine AS builder
RUN apk --update add --no-cache ca-certificates git

COPY ./ /go/src/github.com/jmckind/gvent-api

RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -a -ldflags="-w -s" -o /usr/local/bin/gvent-api github.com/jmckind/gvent-api/cmd/gvent-api \
    && ls -lh /usr/local/bin/gvent-api

# Create the runtime image
FROM scratch AS runtime
LABEL maintainer="John McKenzie <jmckind@gmail.com>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="jmckind/gvent-api"
LABEL org.label-schema.description="Event management system written in Go"
LABEL org.label-schema.url="https://github.com/jmckind/gvent-api"
LABEL org.label-schema.vcs-url="https://github.com/jmckind/gvent-api.git"
LABEL org.label-schema.docker.cmd="docker run -it jmckind/gvent-api"

# Run as "nobody" user
USER 65534:65534

# Copy root certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy pre-built statically linked golang binary and make it the entrypoint
COPY --from=builder /usr/local/bin/gvent-api /usr/local/bin/gvent-api

ENTRYPOINT ["/usr/local/bin/gvent-api"]
