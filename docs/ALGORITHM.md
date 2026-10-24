# FacePointer tracking algorithm

FacePointer uses MediaPipe FaceLandmarker in the WebView to locate one face in each new camera frame. Video is processed on the local device.

## Cursor movement

1. The user looks straight ahead and calibrates. The nose landmark becomes the neutral reference point.
2. Each frame compares the current nose position with that reference. The horizontal coordinate is mirrored to match the preview; both axes are scaled by the pointer-range setting and clamped to the screen.
3. An exponential moving average smooths the target. Its update factor is `1 - smoothing / 100`, so lower smoothing values react faster.
4. The Go controller converts normalized coordinates to Windows screen coordinates and calls `SetCursorPos`.

## Blink actions

MediaPipe's `eyeBlinkLeft` and `eyeBlinkRight` blendshape scores detect a two-eye blink with separate close and reopen thresholds. If blendshape scores are unavailable, the app uses the eye aspect ratio (EAR): vertical eyelid distance divided by horizontal eye width. The user-configurable sensitivity adjusts the close threshold.

A blink is counted after both eyes reopen, provided they were closed for a valid duration. By default, one blink clicks immediately. A second blink inside the Windows double-click window completes the native double-click; a slower pair inside the gesture window sends a deliberate double-click. The cursor pauses while the blink action is processed so it stays on the chosen item.

## Performance

Frames are checked on a timer and face inference runs only when the video has advanced. Pointer calls are single-flight in the UI and serialized in Go, preventing stale movement requests from piling up ahead of clicks.
