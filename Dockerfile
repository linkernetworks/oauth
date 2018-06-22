##
## Build stage
##
FROM golang:1.10-alpine3.7
RUN apk --no-cache add make
WORKDIR /go/src/github.com/linkernetworks/oauth
COPY . ./

ENV CGO_ENABLED 0
ENV GOOS linux
RUN make src.build

##
## final image
##
FROM alpine:3.7
WORKDIR /root/
RUN apk --no-cache add ca-certificates
COPY --from=0  /go/src/github.com/linkernetworks/oauth/build/src/cmd/oauth_server/oauth_server .
ENTRYPOINT ["/root/oauth_server"]
