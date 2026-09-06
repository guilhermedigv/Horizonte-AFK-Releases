from __future__ import annotations

import re
import sys
from pathlib import Path

p = Path(sys.argv[1])
s = p.read_text(encoding="utf-8-sig")

old_version = "$appVersion = '1.6.24-beta.1'"
new_version = "$appVersion = '1.6.25-beta.1'"
if s.count(old_version) != 1:
    raise SystemExit(f"expected exact 1.6.24 runtime version once, found {s.count(old_version)}")
s = s.replace(old_version, new_version, 1)

pattern = re.compile(
    r"(?ms)^function Invoke-HZOneTimeCloudCleanup1619 \{.*?^\}\r?\n\r?\n(?=# Horizonte AFK 1\.6\.7 - Centro de Contas\.)"
)
m = pattern.search(s)
if not m:
    raise SystemExit("Invoke-HZOneTimeCloudCleanup1619 not found")

replacement = r'''function Invoke-HZOneTimeCloudCleanup1619 {
    if ($env:HZ_MULTI_CHILD -eq '1' -or (Test-Path -LiteralPath $script:cloudCleanup1619Marker)) { return }
    try {
        if ([string]::IsNullOrWhiteSpace($script:cloudToken)) { Load-HZCloudClient }
        if ([string]::IsNullOrWhiteSpace($script:cloudToken)) { return }

        # 1.6.25: reconciliation-only migration. Older launchers stopped every
        # visible Cloud session when this local marker was absent. That could kill
        # legitimate 24/7 sessions after reinstall/profile cleanup. Never stop a
        # remote session from a one-time local migration.
        $status = Invoke-HZCloudApi 'GET' '/v1/client/status' $null 8
        $sessions = @($status.sessions)
        if ($sessions.Count -gt 0) {
            Sync-HZAccountsFromCloud $sessions
            Save-HZDeviceRegistry
        }
        $script:cloudSessionsCache = @($sessions)
        $script:cloudSessionsCacheAt = Get-Date
        [IO.File]::WriteAllText(
            $script:cloudCleanup1619Marker,
            ('reconciled|' + [DateTimeOffset]::UtcNow.ToUnixTimeSeconds() + '|' + $sessions.Count),
            (New-Object Text.UTF8Encoding($false))
        )
        if ($sessions.Count -gt 0) {
            Add-SafeActivity ($sessions.Count.ToString() + ' sessão(ões) 24/7 reconciliada(s) sem interromper a Discloud.')
        }
    } catch {
        # A real Cloud failure must not consume the migration marker. Retry safely
        # on the next startup instead of changing remote state.
    }
}

'''
s = s[:m.start()] + replacement + s[m.end():]

health_anchor = """        if ($h.explicit_stop_v2 -ne $true -or $h.detach_safe -ne $true) {
            $CloudStatusText.Text = 'Atualize a Discloud para Cloud v0.5.1 antes de conectar'
            $CloudStatusText.Foreground = '#E5C06E'
            return $false
        }
"""
if s.count(health_anchor) != 1:
    raise SystemExit(f"Cloud capability anchor count={s.count(health_anchor)}")
s = s.replace(
    health_anchor,
    health_anchor + """        if ($h.legacy_owner_reconcile_v1 -ne $true -or $h.dialog_response_ack_v1 -ne $true) {
            $CloudStatusText.Text = 'Atualize a Discloud para Cloud v0.5.8 antes de conectar'
            $CloudStatusText.Foreground = '#E5C06E'
            return $false
        }
""",
    1,
)

# Cross-feature invariants from the exact public 1.6.24 build. These intentionally
# fail the build if a future source extraction silently loses an already-fixed path.
required = [
    "$appVersion = '1.6.25-beta.1'",
    "Decode-HZCoreDialogHex",
    "Source = 'corefile'",
    "Wait-HZCoreDialogTransition",
    "Invoke-HZCoreDialogResponse",
    "$psi.RedirectStandardInput = $true",
    "device-registry-v1.json",
    "accounts-v1.json",
    "HZLocalProcessJob",
    "Add_Closing",
    "Hide-HZToTray",
    "Stop-HZLocalProcessesUnderDataRoot",
    "Add_PreviewKeyDown",
    "Send-HZServerText",
    "Chat do servidor aparece somente em CONEXÃO E CHAT.",
    "multiQueue",
    "Start-NextHZMultiSession",
    "Read-MultiSessionState",
    "legacy_owner_reconcile_v1",
    "dialog_response_ack_v1",
    "reconciled|",
]
for needle in required:
    if needle not in s:
        raise SystemExit(f"runtime regression invariant missing: {needle}")

for forbidden in [
    "ReadMemory($pid32, 0x004E26C8",
    "ReadMemory($pid32, 0x0050CF00",
    "Get-Random",
]:
    if forbidden in s:
        raise SystemExit(f"forbidden regression present: {forbidden}")

start = s.index("function Invoke-HZOneTimeCloudCleanup1619")
end = s.index("# Horizonte AFK 1.6.7 - Centro de Contas.", start)
migration = s[start:end]
for forbidden in ["/stop", "Stop-HZCloudSessionIdSafe", "Set-HZRegistrySessionsOffline"]:
    if forbidden in migration:
        raise SystemExit(f"destructive Cloud migration remains: {forbidden}")
if "Sync-HZAccountsFromCloud $sessions" not in migration:
    raise SystemExit("Cloud migration no longer reconciles account state")

p.write_text(s, encoding="utf-8")
print("patched exact public 1.6.24 runtime -> 1.6.25-beta.1")
