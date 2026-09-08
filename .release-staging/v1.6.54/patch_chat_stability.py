from __future__ import annotations

import hashlib
import struct
import sys
from pathlib import Path

BASE_SHA256 = "91efd565c086c9ed7cd40b741e936a2a903e7b9b091e3b94054899020ae16c03"
OLD_VERSION = b"1.6.53-beta.1"
NEW_VERSION = b"1.6.54-beta.1"

NEW_RESOLVER = r'''function Get-HZPageCloudSession($page,[bool]$force=$false){
if($null-eq$page){return $null};$a=Get-HZPageAccount $page;if($null-eq$a){return $null};$sid=([string](Get-HZAccountCloudSessionId $a)).Trim();$n=([string]$a.Nickname).Trim();$h=([string]$a.Host).Trim();$p=Get-HZScalarInt $a.Port 7777
$cur=([string]$script:cloudSessionId).Trim();if($cur-and[string]$script:activeAccountId-eq[string]$a.Id){try{$r=Invoke-HZCloudApi 'GET' ('/v1/client/sessions/'+$cur) $null 4;if($r){return $r}}catch{}}
$pin='';try{$pin=([string]$a.CloudSessionId).Trim()}catch{};if($pin){try{$r=Invoke-HZCloudApi 'GET' ('/v1/client/sessions/'+$pin) $null 4;$rn='';try{$rn=([string]$r.configured_nickname).Trim()}catch{};if(!$rn){try{$rn=([string]$r.nickname).Trim()}catch{}};$rh='';try{$rh=([string]$r.host).Trim()}catch{};$rp=7777;try{$rp=Get-HZScalarInt $r.port 7777}catch{};if($rn.Equals($n,[StringComparison]::OrdinalIgnoreCase)-and$rh.Equals($h,[StringComparison]::OrdinalIgnoreCase)-and$rp-eq$p){try{Register-HZDeviceCloudSession $a $r}catch{};return $r}}catch{}}
$bad=$false;foreach($r in @(Get-HZCloudSessionsSafe $force)){$rs='';$rn='';$rh='';$rp=7777;try{$rs=([string]$r.session_id).Trim()}catch{};try{$rn=([string]$r.configured_nickname).Trim()}catch{};if(!$rn){try{$rn=([string]$r.nickname).Trim()}catch{}};try{$rh=([string]$r.host).Trim()}catch{};try{$rp=Get-HZScalarInt $r.port 7777}catch{};$same=($rn.Equals($n,[StringComparison]::OrdinalIgnoreCase)-and$rh.Equals($h,[StringComparison]::OrdinalIgnoreCase)-and$rp-eq$p);if($sid-and$rs.Equals($sid,[StringComparison]::OrdinalIgnoreCase)-and!$same){$bad=$true;continue};if($same){try{Register-HZDeviceCloudSession $a $r}catch{};return $r}}
if($bad){try{Set-HZAccountCloudProperty $a 'CloudSessionId' '';Save-HZAccounts;Save-HZDeviceRegistry}catch{}};return $null
}
'''.encode("utf-8")


def pe_map(data: bytes):
    pe = struct.unpack_from("<I", data, 0x3C)[0]
    count = struct.unpack_from("<H", data, pe + 6)[0]
    opt_size = struct.unpack_from("<H", data, pe + 20)[0]
    section = pe + 24 + opt_size
    return [
        (
            data[section + i * 40 : section + i * 40 + 8],
            struct.unpack_from("<II", data, section + i * 40 + 16),
        )
        for i in range(count)
    ]


def main(src: str, dst: str) -> None:
    original = Path(src).read_bytes()
    digest = hashlib.sha256(original).hexdigest()
    if digest != BASE_SHA256:
        raise SystemExit(f"unexpected base hash {digest}")

    patched = bytearray(original)
    root = original.find(b"# Horizonte AFK Launcher")
    start = original.find(b"function Get-HZPageCloudSession", root)
    end = original.find(b"\nfunction ", start + 10)
    if start < 0 or end < 0:
        raise SystemExit("Get-HZPageCloudSession range missing")
    end += 1
    old_len = end - start
    raw = NEW_RESOLVER
    if len(raw) > old_len:
        raise SystemExit(f"new resolver too large {len(raw)} > {old_len}")
    gap = old_len - len(raw)
    if gap == 1:
        raw += b"\n"
    elif gap >= 2:
        raw += b"#" + b" " * (gap - 2) + b"\n"
    if len(raw) != old_len:
        raise SystemExit("fixed-size resolver mismatch")
    patched[start:end] = raw

    pos = 0
    replacements = 0
    while True:
        idx = patched.find(OLD_VERSION, pos)
        if idx < 0:
            break
        patched[idx : idx + len(OLD_VERSION)] = NEW_VERSION
        pos = idx + len(NEW_VERSION)
        replacements += 1
    if replacements < 1:
        raise SystemExit("version marker missing")

    result = bytes(patched)
    if len(result) != len(original):
        raise SystemExit("PE size changed")
    if pe_map(result) != pe_map(original):
        raise SystemExit("PE section map changed")
    for marker in (
        b"ServerClosed",
        b"learned",
        b"Open-HZUcpConnection",
        b"Get-HZActivePageCloudSessionId",
        b"reconciliation-only migration",
    ):
        if marker not in result:
            raise SystemExit(f"missing preserved marker {marker!r}")

    Path(dst).write_bytes(result)
    print("resolver_old", old_len)
    print("resolver_new", old_len - gap)
    print("padding", gap)
    print("version_replacements", replacements)
    print("size", len(result))
    print("sha256", hashlib.sha256(result).hexdigest())


if __name__ == "__main__":
    if len(sys.argv) != 3:
        raise SystemExit("usage: patch_chat_stability.py base.exe Horizonte_AFK.exe")
    main(sys.argv[1], sys.argv[2])
