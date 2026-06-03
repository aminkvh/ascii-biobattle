# ASCII Battle Screensaver

An animated, procedural ASCII art battle between futuristic AI robots and medieval knights. Features a dynamic starry night sky, swaying grass, reflective water lakes, and smooth AI-controlled combat.

## Installation & Setup

This screensaver is a terminal application. Setting it as a native system screensaver depends on your operating system.

### Windows
Windows natively supports `.scr` screensavers.
1. Build the executable as a `.scr` file:
   ```bash
   GOOS=windows GOARCH=amd64 go build -o screensaver.scr .
   ```
2. Right-click the compiled `screensaver.scr` file and select **Install**. Alternatively, copy it directly into `C:\Windows\System32\`.
3. The Screen Saver Settings dialog will open automatically, allowing you to select it and configure the timeout.

### Linux (X11 / xscreensaver)
For Linux systems using X11, `xscreensaver` is the most common way to run custom executables.
1. Build the executable:
   ```bash
   go build -o screensaver .
   ```
2. Move the binary to a secure location (e.g., `/usr/local/bin/screensaver`).
3. Edit your `~/.xscreensaver` configuration file. Find the `programs:` section and add this line:
   ```text
   "ASCII Battle"  /usr/local/bin/screensaver \n\
   ```
4. Run `xscreensaver-demo` and select "ASCII Battle" from the list of screensavers.
*(Note: If using Wayland or swayidle, you can configure your idle daemon to launch a fullscreen terminal running this binary, e.g., `alacritty -e screensaver --fullscreen`)*

### macOS
macOS natively requires `.saver` bundles for screensavers. To use a terminal application:
1. Build the macOS binary:
   ```bash
   GOOS=darwin GOARCH=arm64 go build -o screensaver_mac .
   ```
   *(Use `GOARCH=amd64` for Intel Macs)*
2. Use a free wrapper like **ScriptSaver** or **SaveScreenie**. These wrappers allow you to run any script or terminal executable as a screensaver.
3. Configure the wrapper in your macOS Desktop & Screen Saver settings to point to your compiled `screensaver_mac` binary.
