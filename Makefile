DIST=dist
TARGET=goscheduler
IMAGE_NAME=jobscheduler
GOLANG_IMAGE_TAG=1.17.6-alpine3.15
PORT=20001

all:  cleandarwin macos

docker: dockerbuild dockerimage dockerrun

macos:
	GOPROXY="https://goproxy.cn,direct" CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -o $(DIST)/$(TARGET)-darwin .

linux:
	GOPROXY="https://goproxy.cn,direct" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(DIST)/$(TARGET)-linux .

dockerbuild:
	docker run --rm --env GOPROXY="https://goproxy.cn,direct" --env CGO_ENABLED=0 --env GOOS=linux --env GOARCH=amd64 -v "$(PWD)":/usr/src/myapp -w /usr/src/myapp golang:$(GOLANG_IMAGE_TAG) go build -v -o $(DIST)/$(TARGET)-docker .

dockerimage:
	docker build -t $(IMAGE_NAME) .

dockerrun:
	docker run --rm -p $(PORT):$(PORT) $(IMAGE_NAME) /dist/$(TARGET)-docker

dockerclean:
	rm -rf dist/$(TARGET)-docker
	docker rmi $(IMAGE_NAME)

cleandarwin:
	rm -rf dist/$(TARGET)-darwin

cleanlinux:
	rm -rf dist/$(TARGET)-linux
