# ASCII Biobattle Screensaver

An animated, procedural ASCII art battle between futuristic AI robots and medieval knights. Don't leave your workstation without entertainment at night!

---

## 🚀 Installation & Setup (Using Releases)

### 1. Download the Binary
Visit the [GitHub Releases](https://github.com/aminkvh/ascii-biobattle/releases) page and download the appropriate file for your platform:

### 2. Configure as a Screensaver

<details>
<summary><b>🏁 Windows Setup</b></summary>

Windows natively supports `.scr` screensavers:
1. Download `ascii-biobattle_amd64.scr`.
2. Right-click the `.scr` file and select **Install**. Alternatively, copy it directly into `C:\Windows\System32\`.
3. The Screen Saver Settings dialog will open automatically, allowing you to select it and configure the idle timeout.

</details>

<details>
<summary><b>🐧 Linux Setup (Recommended: xscreensaver)</b></summary>

For Linux systems, using `xscreensaver` is the recommended method. It automatically handles launching the screensaver on idle, locking the screen upon movement, and supports multiple monitors out of the box (no custom keyboard shortcuts needed).

1. Install `xscreensaver`:
   ```bash
   sudo apt install xscreensaver xscreensaver-gl
   ```
2. Download the binary `ascii-biobattle_linux_amd64` (or build it from source) and make it executable:
   ```bash
   chmod +x ascii-biobattle_linux_amd64
   sudo mv ascii-biobattle_linux_amd64 /usr/local/bin/ascii-biobattle
   ```
3. Initialize the configuration:
   - Run `xscreensaver-demo` in your terminal to start the daemon and generate the configuration file.
4. Edit the configuration file `~/.xscreensaver`. Find the `programs:` section and add this line to register the screensaver:
   ```text
   "ASCII Biobattle"  /usr/local/bin/ascii-biobattle --theme night \n\
   ```
5. Choose the screensaver:
   - Run `xscreensaver-demo` again and select **ASCII Biobattle** from the list.
   - (Optional) Under the **Advanced** tab, check the Multi-Monitor settings to display the screensaver on all screens.
6. Auto-start on login:
   - Open **Startup Applications** in Ubuntu, add a new entry with the command `xscreensaver -nosplash`.

*(Note: If you are using Wayland/swayidle, you can configure your idle daemon to launch a fullscreen terminal running the binary instead: `alacritty -e ascii-biobattle --fullscreen`)*

</details>

<details>
<summary><b>🍎 macOS Setup</b></summary>

macOS requires a wrapper utility to run terminal programs as screensavers:
1. Download `ascii-biobattle_mac_arm64` (or `amd64`) and make it executable:
   ```bash
   chmod +x ascii-biobattle_mac_arm64
   ```
2. Download a free screensaver wrapper like **ScriptSaver** or **SaveScreenie**.
3. Configure the wrapper tool in your macOS Desktop & Screen Saver settings to point to your downloaded binary.

</details>
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

<details>
<summary><b>🖥️ Multi-Monitor Support</b></summary>

Since `ascii-biobattle` runs inside a terminal emulator, it relies on the operating system's screensaver manager to handle multi-monitor layouts:

*   **Linux (`xscreensaver`)**: Supports multiple monitors out of the box. Open `xscreensaver-demo`, go to the **Advanced** tab, and under **Multi-Monitor Settings**, select **"Display screensaver on all screens"**.
*   **Windows**: Windows natively runs screensavers on your primary monitor. To span across multiple monitors or run separate instances, you can use a screensaver wrapper utility (like *ScreenSaver Commander*).
*   **macOS**: Screensaver wrappers (like *ScriptSaver* or *SaveScreenie*) include configuration options to run the screensaver script/command on all active displays.

</details>
