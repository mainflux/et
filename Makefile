PROGRAM = callhome
MF_DOCKER_IMAGE_NAME_PREFIX ?= mainflux
SOURCES = $(wildcard *.go) cmd/main.go
CGO_ENABLED ?= 0
GOARCH ?= amd64
VERSION ?= $(shell git describe --abbrev=0 --tags 2>/dev/null || echo "0.13.0")
COMMIT ?= $(shell git rev-parse HEAD)
TIME ?= $(shell date +%F_%T)
DOMAIN ?= callhome.mainflux.com

all: $(PROGRAM)

.PHONY: all clean $(PROGRAM)

define make_docker
	docker build \
		--no-cache \
		--build-arg SVC=$(PROGRAM) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg TIME=$(TIME) \
		--tag=$(MF_DOCKER_IMAGE_NAME_PREFIX)/$(PROGRAM) \
		-f docker/Dockerfile .
endef

define make_docker_ui
	docker build \
	--tag=$(MF_DOCKER_IMAGE_NAME_PREFIX)/$(PROGRAM)-ui \
	-f frontend/Dockerfile ./frontend
endef

define make_dev_cert
	sudo openssl req -x509 -out ./docker/certbot/conf/live/$(DOMAIN)/fullchain.pem \
	-keyout ./docker/certbot/conf/live/$(DOMAIN)/privkey.pem \
	-newkey rsa:2048 -nodes -sha256 \
	-subj '/CN=localhost'
endef

$(PROGRAM): $(SOURCES)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) \
	go build -mod=vendor -ldflags "-s -w \
	-X 'github.com/mainflux/mainflux.BuildTime=$(TIME)' \
	-X 'github.com/mainflux/mainflux.Version=$(VERSION)' \
	-X 'github.com/mainflux/mainflux.Commit=$(COMMIT)'" \
	-o ./build/$(PROGRAM)-$(PROGRAM) cmd/main.go

clean:
	rm -rf $(PROGRAM)

docker-image-server:
	$(call make_docker)
docker-image-ui:
	$(call make_docker_ui)
dev-cert:
	$(call make_dev_cert)

run:
	docker compose -f ./docker/docker-compose.yml up
