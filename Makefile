PKGNAME  := animein
BUILDDIR := build
BINARY   := $(BUILDDIR)/$(PKGNAME)
QUERY 	 := "Keikenzumi na Kimi to, Keiken Zero na Ore ga, Otsukiai suru Hanashi"

# Default target
.PHONY: all
all: build

.PHONY: build
build:
	@mkdir -p $(BUILDDIR)
	@echo "Building for Current OS..."
	go build -o $(BINARY) ./main.go

.PHONY: run
run: build
	@echo "Running $(PKGNAME)..."
	@./$(BINARY)

.PHONY: clean
clean:
	@echo "Cleaning build folder..."
	rm -rf $(BUILDDIR)

.PHONY: install
install: build
	@echo "Installing to $(HOME)/local/bin..."
	cp $(BINARY) $(HOME)/.local/bin/$(PKGNAME)

.PHONY: search
search: build
	@./$(BINARY) "$(QUERY)"

.PHONY: fmt
fmt:
	@echo 'Ngerapiin format kode...'
	go fmt ./...

# vim: ft=make

