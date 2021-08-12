#### Start of system configuration section. ####

# The directory is project's point of entry.
entry_dir := ../cmd/local_library

# The directory to install private app code or library code.
internal_dir :=  ../internal/app

# The directory to install common app code.
service_dir := ../pkg/service

#### End of system configuration section. ####

.PHONY: all
all: format lint vet

format:
	goimports -w $(entry_dir)/*.go $(internal_dir)/*.go $(service_dir)/handlers/*.go $(service_dir)/models/*.go

lint:
	golint $(entry_dir)/
	golint	$(internal_dir)/
	golint	$(service_dir)/...
	golint

vet:
	go vet $(entry_dir)/ $(internal_dir)/ $(service_dir)/...

.PHONY: install
install: all
	go install $(entry_dir)
	go build -o ../ $(entry_dir)
	../local_library
