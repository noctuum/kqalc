PREFIX ?= /usr
DESTDIR ?=

.PHONY: build install uninstall store-bundle clean

STORE_DIR = build/store/kqalc-kde-store
STORE_FILES = kqalc kqalc-amd64 kqalc-arm64 krunner-plugininstallerrc \
	org.kde.krunner1.kqalc.desktop README.md LICENSE

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o kqalc .

install: build
	install -Dm755 kqalc $(DESTDIR)$(PREFIX)/bin/kqalc
	install -Dm644 dist/org.kde.krunner1.kqalc.desktop \
		$(DESTDIR)$(PREFIX)/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop
	install -Dm644 dist/org.kde.krunner1.kqalc.service \
		$(DESTDIR)$(PREFIX)/share/dbus-1/services/org.kde.krunner1.kqalc.service

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/kqalc
	rm -f $(DESTDIR)$(PREFIX)/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop
	rm -f $(DESTDIR)$(PREFIX)/share/dbus-1/services/org.kde.krunner1.kqalc.service

# Bundle for store.kde.org / opendesktop.org, installable by
# krunner-plugininstaller without root. Both architectures ship in one archive:
# KNewStuff resolves updates by version alone, so a second download slot would
# hand arm64 users the amd64 build on their next update.
# The archive stays flat — KNewStuff only wraps it in its own subdirectory when
# the top level holds more than one entry — and its name carries no dot before
# the extension, since that directory is named with QFileInfo::baseName().
store-bundle:
	rm -rf $(STORE_DIR)
	mkdir -p $(STORE_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(STORE_DIR)/kqalc-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(STORE_DIR)/kqalc-arm64 .
	install -m755 dist/kqalc-dispatch.sh $(STORE_DIR)/kqalc
	install -m644 dist/krunner-plugininstallerrc dist/org.kde.krunner1.kqalc.desktop \
		README.md LICENSE $(STORE_DIR)/
	tar -czf build/store/kqalc-kde-store.tar.gz -C $(STORE_DIR) $(STORE_FILES)

clean:
	rm -f kqalc
	rm -rf build/store
