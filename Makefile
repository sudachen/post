CGO_LDFLAGS_EXT = -Wl,-rpath,$(PROJ_DIR)build
include Makefile.Inc
test: get-gpu-setup
	go test ./gpu -v
