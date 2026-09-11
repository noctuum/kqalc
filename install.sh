#!/bin/bash
set -e
CGO_ENABLED=0 go build -ldflags="-s -w" -o kqalc .
install -Dm755 kqalc "$HOME/.local/bin/kqalc"
install -Dm644 dist/org.kde.krunner1.kqalc.desktop \
  "$HOME/.local/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop"

# DBus service file for auto-activation (substitutes actual home path)
sed "s|Exec=.*|Exec=$HOME/.local/bin/kqalc|" dist/org.kde.krunner1.kqalc.service \
  | install -Dm644 /dev/stdin "$HOME/.local/share/dbus-1/services/org.kde.krunner1.kqalc.service"

echo "Installed."
echo
echo "Restart KRunner:  kquitapp6 krunner && kstart6 krunner"
echo
echo "If 'qc 2+2' still finds nothing, the session bus has not picked up the new"
echo "service file. It only watches directories that existed when it started, and"
echo "this install may have just created ~/.local/share/dbus-1/services/. Logging"
echo "out will not help — the user bus outlives the session. Reboot once; later"
echo "installs and updates are picked up immediately."
