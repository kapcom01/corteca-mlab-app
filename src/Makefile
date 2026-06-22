export BUILD_PATH  ?= /app/build
RM                  = rm -rf
DEBUG_MODE    	   ?= false
BIN_PATH            = $(BUILD_PATH)/bin
LIB_PATH  			= $(BUILD_PATH)/lib
PATHS               = $(BIN_PATH) $(LIB_PATH)
export GO_LDFLAGS   = -s -w

ifeq ($(DEBUG_MODE),true)
	GO_EXTRA_FLAGS = 
else 
	GO_EXTRA_FLAGS = -ldflags="$(GO_LDFLAGS)"  
endif

all: ndt7-client mlabs-nokia

ndt7-client: $(BIN_PATH)/ndt7-client
mlabs-nokia: $(BIN_PATH)/mlabs-nokia

$(PATHS):
	@echo "Create target path $@"
	mkdir -p $@

$(BIN_PATH)/ndt7-client: $(PATHS)
	@echo "Building $@"
	@echo $(BUILD_PATH)
	go env
	cd ndt7-client-go && go build -buildvcs=false $(GO_EXTRA_FLAGS) -o $@ ./cmd/ndt7-client

$(BIN_PATH)/mlabs-nokia: $(PATHS)
	@echo "Building $@..."
	SYSROOT=/dist/usr $(MAKE) -C mlabs-nokia
	@echo "Running unit tests $@..."
	$(MAKE) -C mlabs-nokia test

clean:
	$(RM) $(BIN_PATH)/ndt7-client $(BIN_PATH)/mlabs-nokia
	$(MAKE) -C mlabs-nokia clean
	$(RMDIR) $(BIN_PATH) $(LIB_PATH) || true

.PHONY: all clean ndt7-client mlabs-nokia
