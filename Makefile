## oapi generate files
.PHONY: oapi

oapi: 
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=internal/app/config/oapi.yaml --package oapi /home/igor/Desktop/code/GophKeeper/api/api.yaml
