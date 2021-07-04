export VERSION = localnet
BINARY := go-spacemesh
BIN_DIR := build
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=_ -X main.branch=_"

$(info "OS=$(OS)")
ifeq ($(OS),Windows_NT)
  platform := windows
else
  OS := $(shell uname)
  $(info "uname OS=$(OS)")
  ifeq ($(OS),Darwin)
	  platform := macos
  else
      platform := linux
  endif
endif

GPU_SETUP_REV = 0.1.9
GPU_SETUP_ZIP = libgpu-setup-$(platform)-$(GPU_SETUP_REV).zip
GPU_SETUP_URL_ZIP = https://github.com/spacemeshos/gpu-post/releases/download/v0.1.19/$(GPU_SETUP_ZIP)
ifeq ($(platform), windows)
  GPU_SETUP_LIBS = gpu-setup.lib gpu-setup.dll
else
  ifeq ($(platform), macos)
  	GPU_SETUP_LIBS = libgpu-setup.dylib libMoltenVK.dylib libvulkan.1.dylib MoltenVK_icd.json
  else
  	GPU_SETUP_LIBS = libgpu-setup.so
  endif
endif

test: get-gpu-setup
	go test ./gpu -v

BINDIR_GPU_SETUP_LIBS = $(foreach X,$(GPU_SETUP_LIBS),$(BIN_DIR)/$(X))
get-gpu-setup: $(BIN_DIR) $(BIN_DIR)/$(GPU_SETUP_ZIP) $(BINDIR_GPU_SETUP_LIBS)
$(BINDIR_GPU_SETUP_LIBS):
	cd $(BIN_DIR) && unzip $(GPU_SETUP_ZIP) $(notdir $@)
$(BIN_DIR)/$(GPU_SETUP_ZIP):
	curl -L $(GPU_SETUP_URL_ZIP) -o $(BIN_DIR)/$(GPU_SETUP_ZIP)
$(BIN_DIR):
	mkdir -p $(BIN_DIR)
$(BIN_DIR)/$(BINARY):
	go build ${LDFLAGS} -o $(BIN_DIR)/$(BINARY)
