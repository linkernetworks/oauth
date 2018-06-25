FROM golang:1.10-alpine3.7

RUN apk add --no-cache curl
RUN apk add --no-cache docker
RUN apk add --no-cache git
RUN apk add --no-cache sudo
RUN apk add --no-cache bash
RUN apk add --no-cache mongodb
RUN apk add --no-cache make
RUN apk add --no-cache gcc

##
## Go tools
##
RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/jstemmer/go-junit-report
RUN go get -u github.com/t-yuki/gocover-cobertura

##
## bats
##
RUN git clone https://github.com/bats-core/bats-core.git /bats-core && \
    cd /bats-core && \
    ./install.sh /usr/local && \
    cd / && \
    rm -r /bats-core

##
## Make jenkins can do anything
##
RUN echo '%jenkins ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN chmod a+rw /etc/group
RUN chmod a+rw /etc/passwd
RUN echo 'jenkins:!::' >> /etc/gshadow
RUN mkdir -p /home/jenkins && chmod -R a+rw /home/jenkins
ENV HOME=/home/jenkins
COPY entry.sh /entry.sh
RUN chmod +x /entry.sh
ENTRYPOINT /entry.sh


