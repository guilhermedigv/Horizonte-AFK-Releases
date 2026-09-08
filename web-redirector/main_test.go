//go:build windows

package main

import "testing"

func TestPanelURL(t *testing.T) {
	if panelURL != "https://hzrp.discloud.app/panel" {
		t.Fatalf("unexpected panel URL: %s", panelURL)
	}
}

func TestVersion(t *testing.T) {
	if appVersion != "1.6.57-beta.1" {
		t.Fatalf("unexpected version: %s", appVersion)
	}
}
