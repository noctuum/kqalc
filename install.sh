#!/bin/bash
set -e
CGO_ENABLED=0 go build -ldflags="-s -w" -o kqalc .
install -Dm755 kqalc "$HOME/.local/bin/kqalc"
install -Dm644 dist/org.kde.krunner1.kqalc.desktop \
  "$HOME/.local/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop"

# DBus service file for auto-activation (substitutes actual home path)
sed "s|Exec=.*|Exec=$HOME/.local/bin/kqalc|" dist/org.kde.krunner1.kqalc.service \
  | install -Dm644 /dev/stdin "$HOME/.local/share/dbus-1/services/org.kde.krunner1.kqalc.service"

# Start it now instead of waiting for D-Bus to activate it. The bus only watches
# service directories that existed when it started, so on a machine whose
# ~/.local/share/dbus-1/services/ we just created it would not find kqalc at all
# — and the user bus outlives a logout, so signing out would not help either.
# A process that already owns the name needs no activation. From the next login
# on the directory is there and activation takes over.
setsid "$HOME/.local/bin/kqalc" >/dev/null 2>&1 &

echo "Installed and running. Try it: Alt+Space, then 'qc 2+2'."
