#!/bin/bash
# ascii-biobattle-lock.sh
# Spawns one instance per monitor, all sharing the same world via seed + viewport.

# Resolve binary relative to this script — no hardcoded paths
SCRIPT_DIR="$(cd "$(dirname "$(realpath "$0")")" && pwd)"
BINARY="$SCRIPT_DIR/ascii-biobattle"

# Auto-detect DISPLAY if not set (common when called from xautolock)
if [ -z "$DISPLAY" ]; then
    export DISPLAY=$(who | grep -oP '\(:\d+\)' | head -1 | tr -d '()' || echo ":0")
fi

# Generate a shared random seed so all monitors simulate the same world
SEED=$RANDOM$RANDOM

# Get the Ebitengine monitors map: Name -> Index
declare -A EB_MONITORS
while read -r line; do
    if [[ $line =~ Monitor\ ([0-9]+):\ Name=\"([^\"]+)\" ]]; then
        idx="${BASH_REMATCH[1]}"
        name="${BASH_REMATCH[2]}"
        EB_MONITORS["$name"]=$idx
    fi
done < <("$BINARY" --list-monitors 2>/dev/null)

# Get all connected monitors from xrandr: Name, Width, X_offset
# Store as "x_offset:width:name"
RAW_MONITORS=()
while read -r line; do
    if [[ $line =~ ^([^[:space:]]+)[[:space:]]+connected[[:space:]]+(primary[[:space:]]+)?([0-9]+)x[0-9]+\+([0-9]+)\+ ]]; then
        name="${BASH_REMATCH[1]}"
        width="${BASH_REMATCH[3]}"
        x_off="${BASH_REMATCH[4]}"
        RAW_MONITORS+=("$x_off:$width:$name")
    fi
done < <(xrandr --query)

# Sort the monitors by X offset numerically
IFS=$'\n' SORTED_MONITORS=($(sort -n <<<"${RAW_MONITORS[*]}"))
unset IFS

CHAR_WIDTH=10
TOTAL_COLS=0
MON_COLS=()
MON_OFFSETS=()
MON_EB_INDICES=()

for entry in "${SORTED_MONITORS[@]}"; do
    if [ -z "$entry" ]; then continue; fi
    IFS=':' read -r x_off width name <<< "$entry"
    
    eb_idx="${EB_MONITORS[$name]}"
    if [ -z "$eb_idx" ]; then
        echo "Warning: Monitor $name not recognized by Ebitengine"
        continue
    fi
    
    COLS=$((width / CHAR_WIDTH))
    MON_OFFSETS+=($TOTAL_COLS)
    MON_COLS+=($COLS)
    MON_EB_INDICES+=($eb_idx)
    TOTAL_COLS=$((TOTAL_COLS + COLS))
done

NUM_ACTIVE=${#MON_EB_INDICES[@]}
if [ "$NUM_ACTIVE" -lt 1 ]; then
    echo "No matching active monitors found. Launching fallback fullscreen..."
    "$BINARY" --theme night
    loginctl lock-session 2>/dev/null
    exit 0
fi

echo "Seed: $SEED | Total columns: $TOTAL_COLS | Active Monitors: $NUM_ACTIVE"

PIDS=()
for ((i=0; i<NUM_ACTIVE; i++)); do
    OFFSET=${MON_OFFSETS[$i]}
    EB_IDX=${MON_EB_INDICES[$i]}
    echo "  Launch physical monitor $i (Ebitengine index $EB_IDX): viewport-x=$OFFSET"
    "$BINARY" --theme night --monitor "$EB_IDX" \
        --seed "$SEED" \
        --world-cols "$TOTAL_COLS" \
        --viewport-x "$OFFSET" &
    PIDS+=($!)
done

# Wait for first process to exit (user moved mouse)
wait -n "${PIDS[@]}" 2>/dev/null || wait "${PIDS[0]}"

# Kill all remaining
for PID in "${PIDS[@]}"; do
    kill "$PID" 2>/dev/null
done
sleep 0.3

# Lock screen after screensaver exits
loginctl lock-session 2>/dev/null || \
  gnome-screensaver-command --lock 2>/dev/null || \
  xdg-screensaver lock 2>/dev/null
