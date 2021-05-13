LDFLAGS += -X "main.BuildTime=$(shell date "+%F %T %Z")"
LDFLAGS += -X "main.Version=$(shell git rev-parse --short HEAD)"
LDFLAGS += -s -w -extldflags "-static"

GOBUILD := CGO_ENABLED=0 GO111MODULE=auto go build -a -ldflags '$(LDFLAGS)' -trimpath

NO_COLOR=\e[0m
OK_COLOR=\e[32;01m

.PHONY: release generate
.DEFAULT_GOAL := release

generate: # generate bin data
	@printf "$(OK_COLOR)--> generate$(NO_COLOR)\n"
	(cd helpers && go run ./generate/generate.go)

release: # build main bin
	@printf "$(OK_COLOR)--> build linux$(NO_COLOR)\n"
	GOOS=linux $(GOBUILD) -o lp_linux && chmod +x lp_linux
	@printf "$(OK_COLOR)--> build win$(NO_COLOR)\n"
	GOOS=windows $(GOBUILD) -o lp_win.exe
	@printf "$(OK_COLOR)--> build mac intel$(NO_COLOR)\n"
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o lp_mac
	@printf "$(OK_COLOR)--> build mac arm$(NO_COLOR)\n"
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o lp_mac_arm
	@printf "$(OK_COLOR)--> pack bin$(NO_COLOR)\n"
	upx lp_linux lp_mac lp_mac_arm
