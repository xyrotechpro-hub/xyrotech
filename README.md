# NethraMouse AI

NethraMouse AI is a hands-free desktop control project. Its Windows desktop app, FacePointer, uses head movement to guide the normal system cursor and blink gestures for mouse actions. It opens in its own native application window; it does not launch a browser tab.

## Requirements

- Windows 10 or 11
- Go 1.22 or newer
- Microsoft Edge WebView2 Runtime
- Internet access the first time face tracking is loaded

## Run

1. Install Go and the WebView2 Runtime if they are not already installed.
2. Double-click **Start FacePointer.bat**. It downloads/checks Go modules, builds the native app with Wails' production build tags, and opens it. To run from PowerShell instead:

   ```powershell
   go mod tidy
   go run -tags "desktop,production" .
   ```

3. In the FacePointer desktop window, connect your camera and select **Start tracking**. Allow camera access if Windows asks.

The first camera connection downloads MediaPipe's face-tracking runtime and model from Google's and jsDelivr's public CDNs. Face frames and landmarks are processed on your device. MediaPipe may send service usage metrics to Google; see its [privacy notice](https://developers.google.com/edge/mediapipe/solutions/tasks#mediapipe_tasks_privacy_notice).

See [the tracking algorithm](docs/ALGORITHM.md) for details on cursor movement, blink detection, and smoothing.
If setup or Windows security blocks the app, see [troubleshooting](docs/TROUBLESHOOTING.md).

## Gestures

- Move your head to move the normal Windows pointer; the physical mouse remains usable too. Use **Calibrate neutral position** while facing forward to reset the center.
- By default, one deliberate two-eye blink sends a left click immediately to select the item under the pointer. Cursor movement pauses during the blink so the click stays on the item. A second blink within the Windows double-click interval opens it; the app also honors a longer interval configured in Windows. Point at This PC, a file, or a folder first.
- The one-blink menu can also send a right click, hold/release drag, scroll up/down, or Tab. For drag, select **Start/stop drag**, blink once to hold the left button while moving your head, then blink again to release it. Pausing tracking or closing the app releases a held drag.
- The window opens at a size based on the current display, its minimum size is constrained to fit smaller screens, and the layout stacks on narrower windows. Mini mode also stays within the display bounds. Use Tab or Ctrl+Tab from the two-blink menu if you want keyboard navigation instead.
- Pointer movement and clicks work across the Windows desktop, including File Explorer and the combined virtual area of connected monitors. Adjust pointer range and smoothing in the control panel.
- Gesture choices, blink sensitivity, pointer range, and smoothing are saved locally for the next launch.
- Use **Neutral movement deadzone** to ignore tiny involuntary head movements near the calibrated center.
- The most recently selected camera is remembered when its device is still available.
- If face inference stops unexpectedly, the app releases the camera and asks you to reconnect it.
- Select **Restore default controls** to clear saved settings and return to the default click and pointer values.
- Use **Desktop mini mode** to keep a small tracking window in a screen corner while working with other desktop apps. Keep this mini window visible instead of minimizing it so camera tracking continues. The frame loop uses a timer so it does not depend on visible-window animation frames.
- Press **Space** to pause or resume tracking, or **Escape** to pause.

## Build a standalone Windows executable

```powershell
go build -tags "desktop,production" -trimpath -ldflags "-H=windowsgui -s -w" -o FacePointer-Updated.exe .
```

Double-click `FacePointer-Updated.exe` to launch the desktop window without a console. The target PC needs the WebView2 Runtime.
