## oapi generate files
.PHONY: oapi

oapi: 
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=internal/client/oapi/cfg.yaml --package oapi internal/client/oapi/keepercl.yaml
