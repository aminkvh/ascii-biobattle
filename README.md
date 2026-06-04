# ASCII Biobattle Screensaver

Your monitor has been standing there doing nothing for too long. It deserves conflict.

ASCII Biobattle Screensaver turns your idle screen into a chaotic procedural battlefield where futuristic AI robots clash with medieval knights in glorious ASCII combat. While you are away pretending to be productive, tiny warriors and machines fight for control of your pixels across one monitor or many.

---

## 🏁 Setup & Installation

<details>
<summary><b>🏁 Windows Setup</b></summary>

Windows natively supports `.scr` screensavers:
1. Copy `ascii-biobattle.scr` (or rename `ascii-biobattle.exe` to `ascii-biobattle.scr`) from the `build` folder.
2. Right-click the `.scr` file and select **Install**. Alternatively, copy it directly into `C:\Windows\System32\`.
3. The Screen Saver Settings dialog will open automatically. Here, you can select it, configure the idle timeout, and preview it.

</details>

<details>
<summary><b>🐧 Linux Setup (Multi-Monitor & Single Screen)</b></summary>

Because modern graphical engines clash with `xscreensaver`'s antiquated X11 window-ID embedding, the recommended approach on Linux is to use an idle locker like `xautolock` combined with our provided multi-monitor script.

1. **Install xautolock**:
   ```bash
   sudo apt install xautolock
   ```
2. **Configure the lock script**:
   * Open `ascii-biobattle-lock.sh` and make sure `BINARY` points to your compiled executable (e.g. `/media/amin/10TB_2/WORK/screensaver/build/ascii-biobattle`).
3. **Test the script**:
   ```bash
   chmod +x ascii-biobattle-lock.sh
   ./ascii-biobattle-lock.sh
   ```
   *The screensaver will launch across all monitors sorted left-to-right. Move the mouse or press any key to exit and trigger the system lock screen.*
4. **Auto-start on Login**:
   * Search for and open **Startup Applications** from your application menu (or run `gnome-session-properties` in the terminal).
   * Click **Add**.
   * Fill out the fields:
     * **Name**: `ASCII Biobattle Screensaver`
     * **Command**: `xautolock -time 1 -locker "/media/amin/10TB_2/WORK/screensaver/ascii-biobattle-lock.sh" -detectsleep`
     * **Comment**: `Launch screensaver after 1 minute of idle time`
   * Click **Save**.

</details>

<details>
<summary><b>🍎 macOS Setup</b></summary>

macOS strictly requires screensavers to be compiled `.saver` bundles.
1. Download a free screensaver wrapper like [SaverRunner](https://github.com/marnen/saverrunner).
2. Configure it to point to your compiled `ascii-biobattle` binary.
3. It will execute the graphical window fullscreen when the Mac goes idle.

</details>

---

## 💻 Running from Command Line

You can run the application directly:
```bash
./ascii-biobattle [flags]
```

<details>
<summary><b>⚙️ Command Line Flags & Multi-Monitor Configuration</b></summary>

### Available Flags
*   `--theme`: Color theme (`night` [default], `day`)
*   `--density`: Density of combat units on screen (`1` to `80`, default `3`)
*   `--speed`: Animation speed in FPS (`1` to `60`, default `30`)
*   `--intensity`: Color intensity percent (`10` to `100`, default `100`)
*   `--windowed`: Run in a window instead of fullscreen (recommended for development)
*   `--list-monitors`: List connected monitors (useful for debugging multi-display setups)

### Multi-Monitor Viewport Flags (Used by lock script)
*   `--monitor`: Target monitor index (Ebitengine index) to launch fullscreen on
*   `--seed`: Shared random seed so all processes run the identical battle
*   `--world-cols`: Combined width of all monitor viewports (in columns)
*   `--viewport-x`: Horizontal offset (in columns) for this monitor's viewport slice

</details>

---

## 🛠️ How to Build from Source

<details>
<summary><b>🛠️ Build Instructions</b></summary>

Ensure you have Go installed (v1.24+) and the necessary system development libraries for Ebitengine (on Linux: `xorg-dev`, `libgl1-mesa-dev`, `libxcursor-dev`, `libxrandr-dev`, `libxinerama-dev`, `libxi-dev`, `pkg-config`).

1. Build for the host system:
   ```bash
   go build -o build/ascii-biobattle .
   ```
2. Cross-compile for Windows (Standalone and Screen Saver):
   ```bash
   GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/ascii-biobattle.exe .
   cp build/ascii-biobattle.exe build/ascii-biobattle.scr
   ```
3. Cross-compile for macOS (Apple Silicon):
   ```bash
   GOOS=darwin GOARCH=arm64 go build -o build/ascii-biobattle_mac .
   ```

</details>
