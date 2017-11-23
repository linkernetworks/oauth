#Build stage
FROM golang:1.8.3
MAINTAINER Evan Lin <evanlin@linkernetworks.com>
COPY . /go/src/bitbucket.org/linkernetworks/oauth 
WORKDIR /go/src/bitbucket.org/linkernetworks/oauth/cmd/lnk-auth
RUN go get github.com/stretchr/testify/mock
RUN CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o ../../http_server


#final stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0  /go/src/bitbucket.org/linkernetworks/oauth .
CMD ["./http_server"]  
