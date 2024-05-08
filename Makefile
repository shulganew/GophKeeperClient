## oapi generate files
.PHONY: oapi

# install xcode
# https://stackoverflow.com/questions/36303013/can-i-install-xcode-in-ubuntu 
oapi: 
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=internal/app/config/oapi.yaml --package oapi /home/igor/Desktop/code/GophKeeper/api/api.yaml

.PHONY: build_win build_linux build_mac
build_win: export GOOS=windows
build_win: export GOARCH=amd64
build_win: 
	go build -ldflags "-X main.buildVersion=$(git tag -l) -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" -o cmd/gkeepercl/gcl_$(GOOS)_$(GOARCH) cmd/gkeepercl/main.go
build_mac: export GOOS=ios
build_mac: export GOARCH=arm64
build_mac: export CGO_ENABLED=1
build_mac: export OSXCROSS_NO_INCLUDE_PATH_WARNINGS=1
build_mac: export CC=clang
build_mac: 
	go build -ldflags "-X main.buildVersion=$(git tag -l) -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" -o cmd/gkeepercl/gcl_$(GOOS)_$(GOARCH) cmd/gkeepercl/main.go
build_linux: export GOOS=linux
build_linux: export GOARCH=amd64
build_linux: 
	go build -ldflags "-X main.buildVersion=$(git tag -l) -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" -o cmd/gkeepercl/gcl_$(GOOS)_$(GOARCH) cmd/gkeepercl/main.go