FROM golang:1.27-alpine3.24 AS build

RUN adduser -D lpuser
USER lpuser
WORKDIR /lp_app
COPY . .
ENV GO111MODULE=on
ENV CGO_ENABLED=0 
RUN go build -o app .

FROM alpine:3.24
RUN apk add --no-cache --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community hugo
COPY --from=build /lp_app/app .
COPY --from=build /lp_app/web web/
CMD ["./app"]
