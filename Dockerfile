FROM golang:1.20.3-alpine3.17 as build
COPY . .
ENV GO111MODULE=on
ENV GOPATH=$PWD
ENV GOOS=linux
RUN go build -o /test .


FROM alpine:3.17
RUN apk add  --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community hugo
COPY --from=build /test  .
COPY dev_utils .
CMD ["/test"]
COPY /web /web