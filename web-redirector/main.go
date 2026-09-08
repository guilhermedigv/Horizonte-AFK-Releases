//go:build windows

package main

import (
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	appVersion = "1.6.57-beta.1"
	panelURL   = "https://hzrp.discloud.app/panel"
)

func hidden(cmd *exec.Cmd) *exec.Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return cmd
}

func openPanel() error {
	cmd := hidden(exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", panelURL))
	if err := cmd.Start(); err == nil {
		return nil
	}
	fallback := hidden(exec.Command("cmd.exe", "/d", "/s", "/c", "start", "", panelURL))
	return fallback.Start()
}

func messageBox(title, text string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	textPtr, _ := syscall.UTF16PtrFromString(text)
	_, _, _ = proc.Call(
		0,
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		0x00000010,
	)
}

func main() {
	if err := openPanel(); err != nil {
		messageBox(
			"Horizonte AFK",
			"Não foi possível abrir o painel web. Acesse manualmente: https://hzrp.discloud.app/panel",
		)
	}
}
