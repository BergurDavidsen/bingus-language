# Makefile for Bingus

OUTDIR := bin
BINARY := $(OUTDIR)/bingus
SRC := ./cmd/bingus

# Detect host OS/ARCH
HOST_OS := $(shell go env GOOS)
HOST_ARCH := $(shell go env GOARCH)

# Target OS/ARCH, defaults to host
ifeq ($(LINUX),1)
    TARGET_OS := linux
    TARGET_ARCH := amd64
else ifeq ($(MAC),1)
    TARGET_OS := darwin
    TARGET_ARCH := arm64
else
    TARGET_OS := $(HOST_OS)
    TARGET_ARCH := $(HOST_ARCH)
endif

.PHONY: all build clean

all: build

build:
	@mkdir -p $(OUTDIR)
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o $(BINARY) $(SRC)
	@echo "Built $(BINARY) for $(TARGET_OS)/$(TARGET_ARCH)"

clean:
	rm -rf $(OUTDIR)
