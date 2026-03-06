PREFIX ?= /usr
DESTDIR ?=

.PHONY: build install uninstall clean

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

clean:
	rm -f kqalc
