BINARY_NAME=server

run: build
	 @[ "${config_file_path}" ] || ( echo "config_file_path is not set"; exit 1 )
	 ./build/${BINARY_NAME} -config ${config_file_path}

.PHONY: build
build:
	@ $(MAKE) swagger
	GOARCH=amd64 GOOS=linux go build -o build/$(BINARY_NAME) .

.PHONY: clean
clean:
	@go clean
	@rm -rf build/

.PHONY: swagger
swagger:
	@rm -rf docs/
	@swag init

.PHONY: fmt
fmt:
	@go fmt
	@swag fmt