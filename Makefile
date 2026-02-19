BINARY=yt
INSTALL_PATH=$(HOME)/.local/bin/$(BINARY)

.PHONY: build install uninstall clean

build:
	go build -ldflags="-s -w" -o $(BINARY) .

install: build
	cp $(BINARY) $(INSTALL_PATH)
	@echo "Installed to $(INSTALL_PATH)"

uninstall:
	rm -f $(INSTALL_PATH)
	@echo "Removed $(INSTALL_PATH)"

clean:
	rm -f $(BINARY)
