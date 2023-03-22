.PHONY: build linux-build windows-build mac-build download clean test generate doc proto format help

APP_BIN_PATH = bin/app
APP_MAIN_DIR = cmd
API_SWAGGER_SCAN_DIR = internal/app/adapter/server/http
API_SWAGGER_SCAN_ENTRY = http.go
API_SWAGGER_OUT_DIR = internal/app/adapter/server/http/api/docs
API_PROTO_FILES=$(shell find internal/app/adapter/server/grpc/api -name *.proto)
API_PROTO_PB_FILES=$(shell find internal/app/adapter/server/grpc/api -name *.pb.go)

build:
	@make generate
ifeq (${OS}, Windows_NT)
	set CGO_ENABLED=0
	set GOOS=windows
	go build ${BUILD_FLAGS} -o ${APP_BIN_PATH}.exe ${APP_MAIN_DIR}/main.go
else
	CGO_ENABLED=0 go build ${BUILD_FLAGS} -o ${APP_BIN_PATH} ${APP_MAIN_DIR}/main.go
endif

linux-build:
	@make generate
	CGO_ENABLED=0 GOOS=linux go build ${BUILD_FLAGS} -o ${APP_BIN_PATH}_linux ${APP_MAIN_DIR}/main.go

windows-build:
	@make generate
	set CGO_ENABLED=0
	set GOOS=windows
	go build ${BUILD_FLAGS} -o ${APP_BIN_PATH}_windows.exe ${APP_MAIN_DIR}/main.go

mac-build:
	@make generate
	CGO_ENABLED=0 GOOS=darwin go build ${BUILD_FLAGS} -o ${APP_BIN_PATH}_mac ${APP_MAIN_DIR}/main.go

download:
	@go env -w GOPROXY=https://goproxy.cn,direct; \
	go mod download; \
	go get -u github.com/davecgh/go-spew/spew; \
	go get github.com/google/wire/cmd/wire@main; \
	go install github.com/google/wire/cmd/wire@main; \
	go install github.com/cosmtrek/air@latest; \
	go install github.com/swaggo/swag/cmd/swag@latest; \
	go install github.com/golang/mock/mockgen@latest; \
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest; \
	go install github.com/envoyproxy/protoc-gen-validate@latest; \
	go install github.com/favadi/protoc-go-inject-tag@latest; \
	go install entgo.io/ent/cmd/ent@latest; \
	go install github.com/incu6us/goimports-reviser/v3@latest;

clean:
	@if [ -f ${APP_BIN_PATH} ] ; then rm ${APP_BIN_PATH} ; fi

test:
	go test -gcflags=-l -v ${TEST_FLAGS} ./...

generate:
	go generate ./...

doc:
	swag fmt -d ${API_SWAGGER_SCAN_DIR} -g ${API_SWAGGER_SCAN_ENTRY}
	swag init -d ${API_SWAGGER_SCAN_DIR} -g ${API_SWAGGER_SCAN_ENTRY} -o ${API_SWAGGER_OUT_DIR} --parseInternal

proto:
	@$(foreach f, ${API_PROTO_FILES}, kratos proto client --proto_path=./proto $(f);)
	@$(foreach f, ${API_PROTO_PB_FILES}, protoc-go-inject-tag -input=$(f);)

format:
	goimports-reviser -imports-order std,project,company,general -format -recursive ./...

help:
	@printf "%-30s %-100s\n" "make" "automatically compile binaries according to the platform"
	@printf "%-30s %-100s\n" "make build" "automatically compile binaries according to the platform"
	@printf "%-30s %-100s\n" "make linux-build" "compile the binaries for the linux platform"
	@printf "%-30s %-100s\n" "make windows-build" "compile the binaries for the windows platform"
	@printf "%-30s %-100s\n" "make mac-build" "compile the binaries for the mac platform"
	@printf "%-30s %-100s\n" "make download" "download the dependencies required for compilation"
	@printf "%-30s %-100s\n" "make clean" "clean up the binaries generated by the compilation"
	@printf "%-30s %-100s\n" "make test" "unit tests"
	@printf "%-30s %-100s\n" "make generate" "build the files required for the application"
	@printf "%-30s %-100s\n" "make doc" "generate documentation"
	@printf "%-30s %-100s\n" "make proto" "generate the proto file"
	@printf "%-30s %-100s\n" "make format" "format all go files"
