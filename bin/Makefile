LDFLAGS += -X "main.BuildTime=$(shell date "+%F %T %Z")"
LDFLAGS += -X "main.Version=$(shell git rev-parse --short HEAD)"
LDFLAGS += -s -w -extldflags "-static"

GOBUILD := CGO_ENABLED=0 GO111MODULE=auto go build -a -ldflags '$(LDFLAGS)' -trimpath

NO_COLOR=\e[0m
OK_COLOR=\e[32;01m

# for mac arm only
LIPO := $(shell command -v lipo 2> /dev/null)

.PHONY: release generate
.DEFAULT_GOAL := release

generate: # generate bin data
	@printf "$(OK_COLOR)--> generate$(NO_COLOR)\n"
	(cd helpers && go run ./generate/generate.go)

release: # build main bin
	@printf "$(OK_COLOR)--> build linux x86_64$(NO_COLOR)\n"
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o lp_linux_x86_64 && chmod +x lp_linux_x86_64
	@printf "$(OK_COLOR)--> build linux$(NO_COLOR)\n"
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o lp_linux_arm64 && chmod +x lp_linux_arm64
	@printf "$(OK_COLOR)--> build win x86_64$(NO_COLOR)\n"
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o lp_win_x86_64.exe
	@printf "$(OK_COLOR)--> build win arm64$(NO_COLOR)\n"
	GOOS=windows GOARCH=arm64 $(GOBUILD) -o lp_win_arm64.exe
	@printf "$(OK_COLOR)--> build mac x86_64$(NO_COLOR)\n"
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o lp_mac_x86_64
	@printf "$(OK_COLOR)--> build mac arm$(NO_COLOR)\n"
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o lp_mac_arm64
ifdef LIPO
	@printf "$(OK_COLOR)--> pack mac bins$(NO_COLOR)\n"
	lipo -create -output lp_mac_universal lp_mac_x86_64 lp_mac_arm64
endif
