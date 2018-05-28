#Build stage
FROM golang:1.8.3
MAINTAINER Evan Lin <evanlin@linkernetworks.com>
COPY . /go/src/github.com/linkernetworks/oauth 
WORKDIR /go/src/github.com/linkernetworks/oauth/cmd/lnk-auth
RUN go get github.com/stretchr/testify/mock
RUN CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o ../../http_server


#final stage
FROM alpine:3.7  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0  /go/src/github.com/linkernetworks/oauth .
CMD ["./http_server"]  
