//go:build windows

package main

import (
	"strings"
	"testing"
)

func TestFinalVersion(t *testing.T) {
	if appVersion != "1.6.58-beta.1" {
		t.Fatalf("unexpected version: %s", appVersion)
	}
}

func TestPanelURL(t *testing.T) {
	if panelURL != "https://hzrp.discloud.app/panel" {
		t.Fatalf("unexpected panel URL: %s", panelURL)
	}
}

func TestNoticeExplainsRetirementAndCloudSafety(t *testing.T) {
	for _, expected := range []string{"descontinuado", "painel web", "continuam na Cloud"} {
		if !strings.Contains(noticeText, expected) {
			t.Fatalf("notice missing %q", expected)
		}
	}
}
