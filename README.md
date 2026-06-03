# ASCII Biobattle Screensaver

An animated, procedural ASCII art battle between futuristic AI robots and medieval knights. Don't leave your workstation without entertainment at night!

---

## 🚀 Installation & Setup (Using Releases)

### 1. Download the Binary
Visit the [GitHub Releases](https://github.com/aminkvh/ascii-biobattle/releases) page and download the appropriate file for your platform:
*   **Windows**: `ascii-biobattle_amd64.scr` (native screensaver) or `ascii-biobattle_windows_amd64.exe` (run in terminal)
*   **Linux**: `ascii-biobattle_linux_amd64` (Intel/AMD) or `ascii-biobattle_linux_arm64` (ARM64)
*   **macOS**: `ascii-biobattle_mac_amd64` (Intel) or `ascii-biobattle_mac_arm64` (Apple Silicon)

### 2. Configure as a Screensaver

#### Windows
Windows natively supports `.scr` screensavers:
1. Download `ascii-biobattle_amd64.scr`.
2. Right-click the `.scr` file and select **Install**. Alternatively, copy it directly into `C:\Windows\System32\`.
3. The Screen Saver Settings dialog will open automatically, allowing you to select it and configure the idle timeout.

#### Linux (X11 / xscreensaver)
For Linux systems using X11, `xscreensaver` is the most common way to run custom screensavers:
1. Download `ascii-biobattle_linux_amd64` and make it executable:
   ```bash
   chmod +x ascii-biobattle_linux_amd64
   sudo mv ascii-biobattle_linux_amd64 /usr/local/bin/ascii-biobattle
   ```
2. Edit your `~/.xscreensaver` configuration file. Find the `programs:` section and add this line:
   ```text
   "ASCII Biobattle"  /usr/local/bin/ascii-biobattle --theme night \n\
   ```
3. Run `xscreensaver-demo` and select "ASCII Biobattle" from the list.
*(Note: If you are using Wayland/swayidle, you can configure your idle daemon to launch a fullscreen terminal running the binary, e.g., `alacritty -e ascii-biobattle --fullscreen`)*

#### macOS
macOS requires a wrapper utility to run terminal programs as screensavers:
1. Download `ascii-biobattle_mac_arm64` (or `amd64`) and make it executable:
   ```bash
   chmod +x ascii-biobattle_mac_arm64
   ```
2. Download a free screensaver wrapper like **ScriptSaver** or **SaveScreenie**.
3. Configure the wrapper tool in your macOS Desktop & Screen Saver settings to point to your downloaded binary.

---

## 💻 Running from Command Line

You can run the application directly in any terminal:
```bash
./ascii-biobattle_linux_amd64 [flags]
```

### Available Configuration Flags
*   `--theme`: Color theme (`night` [default], `sunset`, `matrix`, `classic`)
*   `--density`: Density of combat units on screen (`1` to `5`, default `3`)
*   `--speed`: Simulation update ticks per second (default `30`)

---

<details>
<summary><b>🛠️ How to Build from Source</b></summary>

If you have Go installed (v1.24+) and want to build the executable yourself:

1. Clone the repository:
   ```bash
   git clone https://github.com/aminkvh/ascii-biobattle.git
   cd ascii-biobattle
   ```

2. Build for your operating system:
   *   **Linux**:
       ```bash
       go build -o ascii-biobattle .
       ```
   *   **macOS (Apple Silicon)**:
       ```bash
       GOOS=darwin GOARCH=arm64 go build -o ascii-biobattle_mac .
       ```
   *   **Windows (Standalone Executable)**:
       ```bash
       GOOS=windows GOARCH=amd64 go build -o ascii-biobattle.exe .
       ```
   *   **Windows (Screensaver binary)**:
       ```bash
       GOOS=windows GOARCH=amd64 go build -o ascii-biobattle.scr .
       ```

</details>
