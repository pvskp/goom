# Goom - Screen Magnifier for X11

Goom is a lightweight screen magnification tool for X11 environments that works by capturing screenshots of the current monitor and zooming into the resulting image.

## Features

- Captures screenshots of the monitor where the cursor is located
- Supports multi-monitor setups via Xinerama
- Automatic detection of the correct monitor in multi-display configurations

## Current Status

This project is in early development. Currently, it can:
- Detect which monitor contains the mouse cursor
- Capture a screenshot of the current monitor
- Save the screenshot as a PNG file

## Planned Features

- Real-time screen magnification
- Customizable zoom levels
- Focus tracking to follow text cursor or mouse pointer
- Hotkeys for zoom control

## Dependencies

- [XGB](https://github.com/jezek/xgb) - X protocol Go language Binding
- X11 with Xinerama extension

## Building

```bash
make build
```

## Usage

Currently, running the program will simply capture a screenshot of the monitor where your cursor is located and save it as "screenshot.png":

```bash
./goom
```

## License

[MIT License](LICENSE)

## Acknowledgments
[tsoding/boomer](https://github.com/tsoding/boomer): main inspiration for this project
