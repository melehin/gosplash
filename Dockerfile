# STEP 1 build executable binary
FROM golang as builder
ENV SRC=$GOPATH/src/github.com/melehin/gosplash
COPY src/ ${SRC}
WORKDIR ${SRC}
# get dependancies
RUN go get -d -v
# build the binary
RUN go build -o /go/bin/gosplash

# STEP 2 build a work image
FROM ubuntu:latest

# Installs latest Chromium package.
RUN apt-get update && apt-get install -y chromium-browser xvfb && apt clean

# Add Chrome as a user
RUN mkdir -p /usr/src/app \
    && useradd --user-group --system --create-home --no-log-init chrome \
    && chown -R chrome:chrome /usr/src/app
# Run Chrome as non-privileged
USER chrome
WORKDIR /usr/src/app

ENV PATH=$PATH:/usr/bin/:/go/bin/
# Copy our static executable
COPY --from=builder /go/bin/gosplash /go/bin/gosplash
EXPOSE 8050
CMD gosplash