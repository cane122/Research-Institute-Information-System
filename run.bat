@echo off
echo Starting Research Institute Information System...
echo.

REM Set PATH to include Go and Wails
set PATH=%PATH%;C:\Program Files\Go\bin;%USERPROFILE%\go\bin

REM Check if executable exists
if exist "build\bin\research-institute-system.exe" (
    echo Running application...
    start "" "build\bin\research-institute-system.exe"
) else (
    echo Building application...
    echo This may take a few minutes...
    wails build
    if exist "build\bin\research-institute-system.exe" (
        echo Build successful! Starting application...
        start "" "build\bin\research-institute-system.exe"
    ) else (
        echo Build failed. Please check the error messages above.
        pause
    )
)
