### START ----- KICS ###
FROM golang:1.25-alpine3.23 AS kics_builder

ARG KICS_VERSION=v2.1.17

WORKDIR /tmp

RUN apk add --no-cache git make

RUN git clone --depth 1 --branch ${KICS_VERSION} https://github.com/Checkmarx/kics.git

WORKDIR /tmp/kics

RUN go mod vendor && make build
### END ----- KICS ###

### START ----- OPENGREP ###
FROM alpine:3.23.0 AS opengrep_builder

ARG OPENGREP_VERSION=v1.13.1

WORKDIR /tmp

RUN apk add --no-cache git bash curl

RUN git clone --depth 1 --branch ${OPENGREP_VERSION} https://github.com/opengrep/opengrep.git

WORKDIR /tmp/opengrep

RUN mkdir -p /tmp/build

ENV HOME=/tmp/build

RUN ./install.sh -v ${OPENGREP_VERSION}

RUN cp /tmp/build/.opengrep/cli/${OPENGREP_VERSION}/opengrep /tmp/build
### END ----- OPENGREP ###

### START ----- GRYPE ###
FROM golang:1.25-alpine3.23 AS grype_builder

ARG GRYPE_VERSION=v0.104.2

WORKDIR /tmp

RUN apk add --no-cache git

RUN git clone --depth 1 --branch ${GRYPE_VERSION} https://github.com/anchore/grype.git

WORKDIR /tmp/grype

RUN ./install.sh
### END ----- GRYPE ###

### START ----- SYFT ###
FROM golang:1.25-alpine3.23 AS syft_builder

ARG SYFT_VERSION=v1.38.2

WORKDIR /tmp

RUN apk add --no-cache git

RUN git clone --depth 1 --branch ${SYFT_VERSION} https://github.com/anchore/syft.git

WORKDIR /tmp/syft

RUN ./install.sh
### END ----- SYFT ###

FROM golang:1.25-alpine3.23

COPY --from=kics_builder /tmp/kics/bin/kics /opt/kics/bin/kics
COPY --from=kics_builder /tmp/kics/assets/queries /opt/kics/assets/queries
COPY --from=kics_builder /tmp/kics/assets/cwe_csv /opt/kics/assets/cwe_csv
COPY --from=kics_builder /tmp/kics/assets/similarityID_transition /opt/kics/assets/similarityID_transition
COPY --from=kics_builder /tmp/kics/assets/libraries/* /opt/kics/assets/libraries/

COPY --from=opengrep_builder /tmp/build/opengrep /opt/opengrep/bin/opengrep
COPY --from=grype_builder /tmp/grype/bin/grype /opt/grype/bin/grype
COPY --from=syft_builder /tmp/syft/bin/syft /opt/syft/bin/syft

WORKDIR /tmp/data

RUN apk add git && git clone https://github.com/Nitr4x/WebGoat.git

ENV SCAN_DIR=/tmp/data

WORKDIR /go/src/