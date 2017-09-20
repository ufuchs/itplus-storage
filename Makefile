GO ?= go
GIT ?= git
INSTALL := install

export GOOS=linux
export GOARCH=arm
export GOARM=7

ifeq ($(OS),Windows_NT)
    GO = /d/Go/bin/go
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CCFLAGS += -D LINUX
    endif
    ifeq ($(UNAME_S),Darwin)
		SHA1 = shasum
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif
#-------------------------------------------------------------------------------
APP := itplus-hubng
VERSION := 0.10
REVL := $(shell git rev-parse HEAD | tail -c 8)
REVH := $(shell git rev-parse HEAD | head -c 7)
REV := $(REVH)..$(REVL)
VER_STRING := $(VERSION) (build $(REV))
#-------------------------------------------------------------------------------
#GOPATH := $(CURDIR)/_vendor:$(GOPATH)
#-------------------------------------------------------------------------------
all: build

version:
	@echo 'package main' > version.go
	@echo 'var (' >> version.go
	@echo '    version = "$(VER_STRING)"' >> version.go
	@echo '    buildDate = "$(shell date)";' >> version.go
	@echo '    builder = "$(LOGNAME)@$(shell hostname)"' >> version.go
	@echo ')' >> version.go

build: version
#	golint ./  -ldflags "-s -w"
	$(GO) build  -o $(APP)

deploy: build
	ansible-playbook ./ansible/provisioning/site.yml

install:
	mount -o remount,rw /mnt/root-ro
	$(INSTALL) -D -m 755 /home/pi/itplus-hub/$(APP) $(DESTDIR)/opt/itplus/hub/$(APP)
	mount -o remount,ro /mnt/root-ro

.PHONY: all version build deploy install
