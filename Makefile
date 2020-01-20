all: deps build-prod

build:
	fileb0x assets.json
	go build cmd/prismriver/prismriver.go

build-dev: frontend-dev build validate

build-prod: frontend build validate

deps:
	cd web && yarn

frontend:
	cd web && yarn run build
	fileb0x web.json

frontend-dev:
	cd web && yarn run dev
	fileb0x web.json

install:
	install -b -D -m644 "prismriver.yml" "/etc/prismriver/prismriver.yml"
	install -D -m755 "prismriver" "/usr/local/bin/prismriver"
	install -D -m644 "prismriver.service" "/usr/lib/systemd/system/prismriver.service"
	install -D -m644 "prismriver-user.service" "/usr/lib/systemd/user/prismriver.service"

run: build-prod
	./prismriver

validate:
	./scripts/validate.sh

.PHONY: all build deps frontend install run
