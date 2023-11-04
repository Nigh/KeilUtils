DIR 		= ./build
EXECUTABLE  	= KeilUtils

GOARCH		= amd64
GOOSWIN		= windows
GO111MODULE		= on
CGO_ENABLED 	= 0

WINBIN 		= $(DIR)/$(EXECUTABLE)-win-$(GOARCH).exe

CC 		= go build
CFLAGS		= 
LDFLAGS		= all=-w -s
GCFLAGS 	= all=
ASMFLAGS 	= all=

.PHONY: all
all: clean win64

.PHONY: win64
win64: $(WINBIN)

.PHONY: $(WINBIN)
$(WINBIN):
	$(CC) $(CFLAGS) -o $(WINBIN) -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -asmflags="$(ASMFLAGS)"
#	Using a compression shell such as upx can compress the binary to about one-third of its original size, but it can be easily misreported as a Trojan by antivirus software.
	./compile-toolset/mpress.exe $(WINBIN)

.PHONY: clean
clean:
	@echo clean
	@- rmdir /Q /S "$(DIR)"
