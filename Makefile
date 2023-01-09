TARGET_OS := linux darwin
TARGET_ARCH := amd64 arm64
PROGRAMS := $(wildcard *.go)
DEST := build

define BuildTask
build-$(1)-$(2): $(PROGRAMS)
	@echo build-$(1)-$(2)
	@GOOS=$(1) GOARCH=$(2) go build -o $(DEST)/$(1)/$(2)/httpq $(PROGRAMS)

endef

define CleanTask
clean-$(1)-$(2): $(DEST)/$(1)/$(2)/httpq
	@echo clean-$(1)-$(2)
	@rm $(DEST)/$(1)/$(2)/httpq

endef

$(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),$(eval $(call BuildTask,$(os),$(arch)))))
$(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),$(eval $(call CleanTask,$(os),$(arch)))))

.PHONY:show
show:
	@$(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),echo $(os)-$(arch);))

.PHONY:build
build: $(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),build-$(os)-$(arch)))
	@echo $@ done

.PHONY: clean
clean: $(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),clean-$(os)-$(arch)))
	@rm -rf build
	@echo $@ done
