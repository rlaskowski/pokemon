build_name := pokemon
port := 8080
os := darwin

ifeq ($(os),windows)
	build_name := $(build_name).exe
endif

build:
	@-${MAKE} clean
	GOOS=$(os) go build -o dist/$(build_name) cmd/pokemon/main.go

clean:
	@rm -Rf dist test_results

run:
	go run cmd/pokemon/main.go 