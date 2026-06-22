BUILD_PATH    ?= /app/build
RM             = rm -rf
DEBUG_MODE    ?= false
SYSROOT       ?= $(abspath /)
BIN_PATH       = $(BUILD_PATH)/bin
PATHS          = $(BIN_PATH) $(LIB_PATH)
GO_LDFLAGS     = -s -w
LIBRARIES_DIR  = $(SYSROOT)/lib
LIBRARIES      =

ifeq ($(DEBUG_MODE),true)
	GO_EXTRA_FLAGS = 
else 
	GO_EXTRA_FLAGS = -ldflags="$(GO_LDFLAGS)"  
endif

$(BIN_PATH)/mlabs-nokia: | $(PATHS)
	env && go env
	go build -buildvcs=false ${GO_EXTRA_FLAGS} -o $@ .
	
$(PATHS):
	@echo "Create target path $@"
	mkdir -p $@

test:
	go test -tags test

clean:
	$(RM) $(BIN_PATH)
