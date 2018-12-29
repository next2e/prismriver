all: deps build

build: frontend
	go build cmd/prismriver/prismriver.go

deps:
	dep ensure
	cd web && yarn

frontend:
	cd web && yarn run build
	statik -src=web/dist -f

install:
	install -D -m755 "prismriver" "/usr/local/bin/prismriver"

run: build
	./prismriver

.PHONY: all build deps frontend install run