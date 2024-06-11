FROM golang:1.20.3-alpine3.17 as build

WORKDIR /lp_app
COPY . .
ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN go build -o app .

FROM alpine:3.17
RUN apk add --no-cache --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community hugo
RUN wget https://github.com/CloudCannon/pagefind/releases/download/v1.1.0/pagefind-v1.1.0-x86_64-unknown-linux-musl.tar.gz && \
    tar -xvf pagefind-v1.1.0-x86_64-unknown-linux-musl.tar.gz && \
    cp pagefind /usr/bin
COPY --from=build /lp_app/app .
COPY --from=build /lp_app/web web/
CMD ["./app"]
