# Troubleshooting

## Camera does not connect

- Allow camera access for FacePointer in Windows privacy settings.
- Close other apps that may be using the camera, then reconnect it in FacePointer.
- The first face-tracking setup needs internet access to download MediaPipe's runtime and model.
- Install the Microsoft Edge WebView2 Runtime if the camera controls are unavailable.

## Cursor does not move or click

- Confirm **FACE / EYE STATUS** reports a detected face, then calibrate while looking straight ahead.
- Keep your face centered and evenly lit. Increase pointer range if you cannot reach desktop icons; increase the deadzone if the cursor drifts while your head is still.
- For desktop use, choose **Desktop mini mode** and leave its window visible. Minimizing the window can throttle camera processing.
- A single blink selects; two quick blinks open the item beneath the cursor by default. Check that the matching actions are selected in **Gesture mapping**.

## Windows blocks a locally built executable

Local builds are unsigned. Windows Smart App Control or an organization policy may block an executable whose publisher cannot be verified. Use a trusted, signed release or ask the device administrator about the policy. Do not run a build whose source you do not trust.
