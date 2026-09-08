//go:build windows

package main

import (
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	appVersion = "1.6.58-beta.1"
	panelURL   = "https://hzrp.discloud.app/panel"
	noticeText = "O menu desktop do Horizonte AFK foi descontinuado.\n\nTodo o AFK agora funciona pelo painel web da Discloud. Suas contas 24/7 continuam na Cloud mesmo com este EXE fechado.\n\nDeseja abrir o painel web agora?"
)

const (
	mbYesNo           = 0x00000004
	mbIconInformation = 0x00000040
	mbSetForeground   = 0x00010000
	mbOK              = 0x00000000
	mbIconError       = 0x00000010
	idYes             = 6
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

func messageBox(title, text string, flags uintptr) uintptr {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	textPtr, _ := syscall.UTF16PtrFromString(text)
	result, _, _ := proc.Call(
		0,
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		flags,
	)
	return result
}

func main() {
	title := "Horizonte AFK " + appVersion + " • Menu descontinuado"
	choice := messageBox(title, noticeText, mbYesNo|mbIconInformation|mbSetForeground)
	if choice != idYes {
		return
	}
	if err := openPanel(); err != nil {
		messageBox(
			"Horizonte AFK",
			"Não foi possível abrir o painel web automaticamente.\n\nAcesse: "+panelURL,
			mbOK|mbIconError|mbSetForeground,
		)
	}
}
