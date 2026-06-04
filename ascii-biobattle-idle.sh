#!/bin/bash
# ascii-biobattle-idle.sh
# Replaces xscreensaver with xautolock + ascii-biobattle

BINARY="$(dirname "$(realpath "$0")")/ascii-biobattle"
TIMEOUT_MINS=9  # Change this to how many minutes before screensaver starts

# Kill any stray xscreensaver/xautolock instances
pkill -x xscreensaver 2>/dev/null
pkill -x xautolock   2>/dev/null

echo "Starting ascii-biobattle idle watcher (timeout: ${TIMEOUT_MINS} min)"
echo "Binary: $BINARY"

# xautolock: after TIMEOUT_MINS of idle, run ascii-biobattle fullscreen.
# -locker is the screensaver program. -killtime 60 kills it after 60s if no input.
exec xautolock \
    -time "$TIMEOUT_MINS" \
    -locker "$BINARY --theme night" \
    -detectsleep \
    -secure \
    -notify 10 \
    -notifier "echo about to lock"
