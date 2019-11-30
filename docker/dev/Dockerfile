FROM golang:1.13

ARG ostype=Linux

RUN apt-get update && apt-get install -y \
    git \
    bash \
    zip

# Mock creator
RUN go get -u github.com/vektra/mockery/... && \
    go get -u github.com/golangci/golangci-lint/...

# Create user.
ARG uid=1000
ARG gid=1000

RUN bash -c 'if [ ${ostype} == Linux ]; then addgroup -gid $gid app; else addgroup app; fi && \
    adduser --disabled-password -uid $uid --ingroup app --gecos "" app && \
    chown app:app -R /go'

# Fill go mod cache.
RUN mkdir /tmp/cache
COPY go.mod /tmp/cache
COPY go.sum /tmp/cache
RUN chown app:app -R /tmp/cache
USER app
RUN cd /tmp/cache && \
    go mod download

WORKDIR /src
