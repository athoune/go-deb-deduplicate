FROM golang:1.25-alpine3.22 AS dev

RUN apk update && apk add --no-cache make
COPY . /src
WORKDIR /src
RUN go install
RUN make test bin

FROM alpine:3.22
COPY --from=dev /data/bin/* /usr/local/bin/
