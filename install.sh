#!/bin/bash
set -e
CGO_ENABLED=0 go build -ldflags="-s -w" -o kqalc .
install -Dm755 kqalc "$HOME/.local/bin/kqalc"
install -Dm644 dist/org.kde.krunner1.kqalc.desktop \
  "$HOME/.local/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop"

# The service file lets the bus start kqalc on demand from the binary above.
sed "s|Exec=.*|Exec=$HOME/.local/bin/kqalc|" dist/org.kde.krunner1.kqalc.service \
  | install -Dm644 /dev/stdin "$HOME/.local/share/dbus-1/services/org.kde.krunner1.kqalc.service"

# A copy from an earlier install holds the bus name and keeps answering from
# the old binary.
pkill -x -u "$USER" kqalc || true

# The bus watches only the service directories that existed when it started,
# so it cannot start kqalc while ~/.local/share/dbus-1/services/ is new to it.
# A running kqalc needs no activation, and it stays up until the bus can start
# it on demand.
setsid "$HOME/.local/bin/kqalc" >/dev/null 2>&1 &

echo "Installed and running. Try it: Alt+Space, then 'qc 2+2'."
