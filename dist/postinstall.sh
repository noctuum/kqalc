#!/bin/sh
# The service file goes to /usr/share/dbus-1/services/, which the bus already
# watches. Versions up to 0.2.2 never exit on their own, so after an upgrade
# from one of them the old binary keeps answering until it is stopped.
echo "kqalc installed. Restart KRunner to pick it up:"
echo "  systemctl --user restart plasma-krunner.service"
echo "After an upgrade from 0.2.2 or earlier, stop the old copy once:"
echo "  pkill -x kqalc"
