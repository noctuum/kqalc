#!/bin/bash
set -e
CGO_ENABLED=0 go build -ldflags="-s -w" -o kqalc .
install -Dm755 kqalc "$HOME/.local/bin/kqalc"
install -Dm644 dist/org.kde.krunner1.kqalc.desktop \
  "$HOME/.local/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop"

# DBus service file for auto-activation (substitutes actual home path)
sed "s|Exec=.*|Exec=$HOME/.local/bin/kqalc|" dist/org.kde.krunner1.kqalc.service \
  | install -Dm644 /dev/stdin "$HOME/.local/share/dbus-1/services/org.kde.krunner1.kqalc.service"

echo "Installed. Restart KRunner: kquitapp6 krunner && kstart6 krunner"
