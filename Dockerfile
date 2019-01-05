FROM golang:1.10.5-alpine3.8 as build

RUN apk add --no-cache dep
RUN apk add --no-cache gcc
RUN apk add --no-cache git
RUN apk add --no-cache libc-dev
RUN apk add --no-cache make
RUN apk add --no-cache vlc-dev
RUN apk add --no-cache yarn
RUN go get github.com/rakyll/statik

COPY . /go/src/gitlab.com/ttpcodes/prismriver

WORKDIR /go/src/gitlab.com/ttpcodes/prismriver
RUN make

FROM alpine:3.8

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache ffmpeg
RUN apk add --no-cache vlc
RUN apk add --no-cache youtube-dl

COPY --from=build /go/src/gitlab.com/ttpcodes/prismriver/prismriver /usr/local/bin/prismriver

CMD ["/usr/local/bin/prismriver"]