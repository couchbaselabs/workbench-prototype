ARG PROD_VERSION=1.0.0
ARG PROD_BUILD=999
FROM golang:1.16 as builder
# Builder image so ignore pinning and we do want recommended packages
# hadolint ignore=DL3008,DL3015
RUN apt-get update && apt-get install -y sqlcipher libssl-dev openssl openssh-client && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /src
# Not best practice but means we cache the Go download during build
COPY ./go.* /src/
RUN go mod download
COPY ./ /src/
# Statically build it to allow reuse on Alpine Linux
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -a -ldflags '-linkmode external -extldflags "-static"' -o /bin/cbmultimanager ./cluster-monitor/cmd/cbmultimanager && \
    CGO_ENABLED=1 GOOS=linux go build -trimpath -a -ldflags '-linkmode external -extldflags "-static"' -o /bin/cbeventlog ./cluster-monitor/cmd/cbeventlog

FROM node:16-alpine3.14 as ui-builder
WORKDIR /src
COPY ui/ /src/
RUN npm install && npm run build

FROM alpine:3.14
ARG PROD_VERSION
ARG PROD_BUILD

# Use latest rather than pin
# hadolint ignore=DL3018
RUN apk add --upgrade --no-cache ca-certificates bash openssl wget

# Ensure we include licensing information
COPY ./LICENSE /licenses/couchbase.txt
# NOTE: this creates a chicken-and-egg effect - if we update notices, we need to run two builds to ensure it gets picked up
RUN wget -q -O /licenses/notices.txt https://raw.githubusercontent.com/couchbase/product-metadata/master/couchbase-cluster-monitor/blackduck/${PROD_VERSION}/notices.txt

# Generate a certificate
RUN mkdir -p /data /priv && \
    openssl req -x509 -nodes -days 365 \
    -subj "/C=CA/ST=QC/O=Couchbase, Inc./CN=couchbase.com" \
    -addext "subjectAltName=DNS:couchbase.com" \
    -newkey rsa:2048 -keyout /priv/server.key -out /priv/server.crt

# Copy in all the executables we need plus a launch script
COPY --from=builder /bin/cbmultimanager /bin/cbmultimanager
COPY --from=builder /bin/cbeventlog /bin/cbeventlog
COPY --from=ui-builder /src/dist/app/ /ui/
COPY ./entrypoint.sh /entrypoint.sh

# Provide a configurable user to run as and for ownership
ARG CB_UID="8453"
ARG CB_GID="8453"

RUN addgroup -g $CB_GID -S couchbase && \
    adduser -u $CB_UID -S couchbase -G couchbase

# Ensure we have everything we need set up with correct permissions
RUN chmod a+x /bin/cb* && \
    chmod a+x /entrypoint.sh && \
    chown -R couchbase:couchbase /ui /data /priv && \
    chmod -R a+r /licenses

EXPOSE 7196 7197

# Run as non-root as per best practices and requirements for some platforms
USER $CB_UID
CMD [ "/entrypoint.sh" ]