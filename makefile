DIR 		= ./build
EXECUTABLE  	= KeilUtils

GOARCH		= amd64
GOOSWIN		= windows
GOOSX		= darwin
GOOSLINUX	= linux
GOMOD		= on
CGO_ENABLED 	= 0

WINBIN 		= $(DIR)/$(EXECUTABLE)-win-$(GOARCH).exe
OSXBIN 		= $(DIR)/$(EXECUTABLE)-darwin-$(GOARCH)
LINUXBIN 	= $(DIR)/$(EXECUTABLE)-linux-$(GOARCH)

CC 		= go build
CFLAGS		= 
LDFLAGS		= all=-w -s
GCFLAGS 	= all=
ASMFLAGS 	= all=

.PHONY: all
all: win64

.PHONY: win64
win64: $(WINBIN)

.PHONY: $(WINBIN)
$(WINBIN):
	GO111MODULE=$(GOMOD) GOARCH=$(GOARCH) GOOS=$(GOOSWIN) CGO_ENABLED=$(CGO_ENABLED) $(CC) $(CFLAGS) -o $(WINBIN) -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -asmflags="$(ASMFLAGS)"
#	Using a compression shell such as upx can compress the binary to about one-third of its original size, but it can be easily misreported as a Trojan by antivirus software.
# 	upx --best --lzma $(WINBIN)

.PHONY: clean
clean:
	rm -rf $(DIR)/*
