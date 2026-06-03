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

<details>
<summary><b>🏁 Windows Setup</b></summary>

Windows natively supports `.scr` screensavers:
1. Download `ascii-biobattle_amd64.scr`.
2. Right-click the `.scr` file and select **Install**. Alternatively, copy it directly into `C:\Windows\System32\`.
3. The Screen Saver Settings dialog will open automatically, allowing you to select it and configure the idle timeout.

</details>

<details>
<summary><b>🐧 Linux Setup (X11 / xscreensaver / Custom Shortcut)</b></summary>

#### A. Classic Setup via `xscreensaver` (X11)
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

#### B. GNOME / Ubuntu Lock Screen Setup
Since default Ubuntu/GNOME does not support custom third-party screensavers out of the box, you can choose one of the following setups to lock your screen:

##### 1. Lock screen with xscreensaver
If you installed `xscreensaver` (via the steps above), open `xscreensaver-demo` and check the **Lock Screen After** checkbox. Moving the mouse will automatically stop the battle and prompt for your password.

##### 2. Direct Keyboard Shortcut with Auto-Lock (e.g. `Ctrl + Alt + K`)
To start the screensaver instantly and lock the workstation when you exit it:
1. Copy the binary:
   ```bash
   sudo cp ./ascii-biobattle /usr/local/bin/ascii-biobattle
   ```
2. Go to **Settings** -> **Keyboard** -> **Keyboard Shortcuts** -> **View and Customise Shortcuts** -> **Custom Shortcuts** -> **Add Shortcut (+)**.
3. Fill in:
   *   **Name**: `ASCII Biobattle Screensaver`
   *   **Command**: `bash -c "gnome-terminal --full-screen -- /usr/local/bin/ascii-biobattle; dbus-send --type=method_call --dest=org.gnome.ScreenSaver /org/gnome/ScreenSaver org.gnome.ScreenSaver.Lock"`
   *   **Shortcut**: Press a shortcut key combination (e.g., `Ctrl + Alt + K` or `Super + Alt + S`). *Note: Standard `Ctrl + K` is often reserved by GNOME or web browsers for internal actions and might not trigger globally.*
4. Press your shortcut anytime to trigger it. When you press `Ctrl + C` or close the terminal, GNOME will immediately lock your screen.

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
