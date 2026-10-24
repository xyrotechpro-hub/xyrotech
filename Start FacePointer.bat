@echo off
setlocal
title FacePointer
cd /d "%~dp0"

where go >nul 2>nul
if errorlevel 1 (
  echo Go 1.22 or newer is required to start FacePointer.
  echo Install Go from https://go.dev/dl/ and then run this file again.
  pause
  exit /b 1
)

echo Building FacePointer desktop app...
go mod tidy
if errorlevel 1 (
  echo Dependency setup failed. Check your internet connection and Go installation.
  pause
  exit /b 1
)

go build -tags "desktop,production" -trimpath -ldflags "-H=windowsgui -s -w" -o FacePointer-Updated.exe .
if errorlevel 1 (
  echo Build failed. Check your Go installation and internet connection.
  pause
  exit /b 1
)

start "" "FacePointer-Updated.exe"
