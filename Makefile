export BUILD_PATH ?= /app/build
RM                := rm -rf

build:
	@echo "Building application"
	$(MAKE) -C src/

distclean:
	$(RM) build/ dist/

clean:
	$(MAKE) -C src/ clean
	
precompiled:
	@echo "You can add here the commands to download precompiled binaries"
	@echo "Put architecture agnostic files in './target/noarch' (scripts, configuration, etc)"
	@echo "put architecture specific in the respective folder eg:'./target/armv8' (precompiled binaries ot libraries, etc)"
	@echo "in order for them to be included in the target root fs"

docs:
	@echo "You can add here the commands to generate the documentation if needed"

.PHONY: build clean precompiled docs distclean
