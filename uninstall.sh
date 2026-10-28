#!/bin/bash
rm -f "$HOME/.local/bin/kqalc"
rm -f "$HOME/.local/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop"
rm -f "$HOME/.local/share/dbus-1/services/org.kde.krunner1.kqalc.service"

# The running copy is stopped last. While the service file was in place, the
# bus would have started a new one on the next query.
pkill -x -u "$USER" kqalc

echo "Uninstalled. Restart KRunner: systemctl --user restart plasma-krunner.service"
