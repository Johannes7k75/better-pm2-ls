@echo off
mkdir bin

setlocal enabledelayedexpansion

set "original_os=%GOOS%"
set "original_arch=%GOARCH%"

rem Define OS targets
set "OS_TARGETS[0]=windows"
set "OS_TARGETS[1]=linux"
set "OS_TARGETS[2]=darwin"

rem Define Architecture targets
set "ARCH_TARGETS_windows=386 amd64"
set "ARCH_TARGETS_linux=386 amd64 arm arm64"
set "ARCH_TARGETS_darwin=amd64 arm64"

for /L %%t in (0, 1, 2) do (
    rem Get OS target
    set "GOOS=!OS_TARGETS[%%t]!"

    rem Determine corresponding architectures for the current OS
    if "!GOOS!"=="windows" (
        set "ARCH_LIST=!ARCH_TARGETS_windows!"
    ) else if "!GOOS!"=="linux" (
        set "ARCH_LIST=!ARCH_TARGETS_linux!"
    ) else if "!GOOS!"=="darwin" (
        set "ARCH_LIST=!ARCH_TARGETS_darwin!"
    )

    rem Iterate over each architecture for the current OS
    for %%a in (!ARCH_LIST!) do (
        set "GOARCH=%%a"
        set "BINARY_NAME=pmls-!GOOS!-!GOARCH!"
        
        echo Building !BINARY_NAME!

        if "!GOOS!"=="windows" (
            set "BINARY_NAME=!BINARY_NAME!.exe"
        )

        go build -trimpath -o "bin/!BINARY_NAME!" main.go 
        echo "Built for !GOOS! (!GOARCH!) as bin/!BINARY_NAME!"
    )
)

set "GOOS=%original_os%"
set "GOARCH=%original_arch%"