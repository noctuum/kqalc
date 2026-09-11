#!/bin/bash
rm -f "$HOME/.local/bin/kqalc"
rm -f "$HOME/.local/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop"
rm -f "$HOME/.local/share/dbus-1/services/org.kde.krunner1.kqalc.service"

# Only now stop a running instance: while the service file was still in place
# the bus would have activated a fresh one on the next query.
pkill -x -u "$USER" kqalc

echo "Uninstalled. Restart KRunner: kquitapp6 krunner && kstart6 krunner"
