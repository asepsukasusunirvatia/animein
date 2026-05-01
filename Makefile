PKGNAME  := animein
BUILDDIR := build
PREFIX 	 := "$(HOME)/.local"
BINARY   := $(BUILDDIR)/$(PKGNAME)
QUERY 	 := "Keikenzumi na Kimi to, Keiken Zero na Ore ga, Otsukiai suru Hanashi"
LDFLAGS  := -s -w


# Default target
.PHONY: all
all: build

.PHONY: build
build:
	@mkdir -p $(BUILDDIR)
	@echo "Building for Current OS..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) ./main.go
	@echo "Done..."

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
	@mkdir -p "$(PREFIX)/bin"
	@echo "Installing to $(HOME)/local/bin..."
	cp $(BINARY) "$(PREFIX)/bin/$(PKGNAME)"

.PHONY: uninstall
uninstall:
	rm -f "$(PREFIX)/bin/$(PKGNAME)"

.PHONY: search
search: build
	@./$(BINARY) "$(QUERY)"

.PHONY: fmt
fmt:
	@echo 'Ngerapiin format kode...'
	go fmt ./...

# vim: ft=make

