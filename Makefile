all: deps build

build: frontend
	fileb0x assets.json
	go build cmd/prismriver/prismriver.go

build-dev: frontend-dev
	fileb0x assets.json
	go build cmd/prismriver/prismriver.go

deps:
	dep ensure
	cd web && yarn

frontend:
	cd web && yarn run build
	fileb0x web.json

frontend-dev:
	cd web && yarn run dev
	fileb0x web.json

install:
	install -D -m755 "prismriver" "/usr/local/bin/prismriver"
	install -D -m644 "prismriver.service" "/usr/lib/systemd/system/prismriver.service"

run: build
	./prismriver

.PHONY: all build deps frontend install run