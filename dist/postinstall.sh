#!/bin/sh
# Packages drop the service file into /usr/share/dbus-1/services/, which the
# message bus already watches — only KRunner needs a nudge.
echo "kqalc installed. Restart KRunner to pick it up:"
echo "  kquitapp6 krunner && kstart6 krunner"
