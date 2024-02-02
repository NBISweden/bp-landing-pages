FROM golang:1.20.3-alpine3.17 as build
WORKDIR /lp_app
COPY . .
ENV GO111MODULE=on
ENV GOPATH=$PWD
ENV GOOS=linux
RUN go build -o /test .


FROM alpine:3.17
WORKDIR /gen_app
RUN apk add  --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community hugo
COPY --from=build /test /gen_app/
COPY dev_utils .
CMD ["/gen_app/test"]
COPY /web /web