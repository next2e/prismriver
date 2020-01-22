FROM golang:1.13.6-alpine3.11 as build

RUN apk add --no-cache gcc
RUN apk add --no-cache git
RUN apk add --no-cache libc-dev
RUN apk add --no-cache make
RUN apk add --no-cache vlc-dev
RUN apk add --no-cache yarn
RUN go get -u github.com/UnnoTed/fileb0x
RUN go get -u golang.org/x/lint/golint

COPY . /srv

WORKDIR /srv
RUN make

FROM alpine:3.11.3

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache ffmpeg
RUN apk add --no-cache vlc
RUN apk add --no-cache youtube-dl

COPY --from=build /srv/prismriver /usr/local/bin/prismriver

CMD ["/usr/local/bin/prismriver"]
