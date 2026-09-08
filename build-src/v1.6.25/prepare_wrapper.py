from __future__ import annotations

import hashlib
import re
import sys
from pathlib import Path

root = Path(sys.argv[1])
main = root / "main.go"
payload = root / "payload"
s = main.read_text(encoding="utf-8")

core_hash = hashlib.sha256((payload / "RakSAMPClient.exe").read_bytes()).hexdigest()
lua_hash = hashlib.sha256((payload / "lua51.dll").read_bytes()).hexdigest()

old_embed = "//go:embed payload/Horizonte-AFK.ps1 payload/RakSAMPClient.exe payload/hero_bg.jpg payload/hz_icon.png"
new_embed = "//go:embed payload/Horizonte-AFK.ps1 payload/RakSAMPClient.exe payload/lua51.dll payload/hero_bg.jpg payload/hz_icon.png"
if s.count(old_embed) != 1:
    raise SystemExit(f"wrapper embed anchor count={s.count(old_embed)}")
s = s.replace(old_embed, new_embed, 1)

old_version = 'appVersion           = "1.5.0"'
if s.count(old_version) != 1:
    raise SystemExit(f"wrapper version anchor count={s.count(old_version)}")
s = s.replace(old_version, 'appVersion           = "1.6.25-beta.1"', 1)

s, n = re.subn(
    r'expectedEngineSHA256 = "[0-9a-f]{64}"',
    f'expectedEngineSHA256 = "{core_hash}"',
    s,
    count=1,
)
if n != 1:
    raise SystemExit("engine hash constant not found")

hash_anchor = f'\texpectedEngineSHA256 = "{core_hash}"\n\tmanifestURL'
if hash_anchor not in s:
    raise SystemExit("lua hash insertion anchor missing")
s = s.replace(
    hash_anchor,
    f'\texpectedEngineSHA256 = "{core_hash}"\n\texpectedLuaSHA256    = "{lua_hash}"\n\tmanifestURL',
    1,
)

# The extracted 1.6.24 PowerShell is already the complete runtime. Reapplying the
# old v1.5 feature text patch here would duplicate handlers/functions.
feature_call = "\tps = applyFeaturePatch(ps)\n"
if s.count(feature_call) != 1:
    raise SystemExit(f"applyFeaturePatch call count={s.count(feature_call)}")
s = s.replace(feature_call, "", 1)

old_root = 'root := filepath.Join(local, "HorizonteAFK", "runtime-v9-rpc-dialogs")'
if s.count(old_root) != 1:
    raise SystemExit(f"runtime root anchor count={s.count(old_root)}")
s = s.replace(
    old_root,
    'root := filepath.Join(local, "HorizonteAFK", "runtime-v13-cloud-reconcile")',
    1,
)

engine_block = '''\tif fileSHA256(enginePath) != expectedEngineSHA256 {\n\t\tif e := os.WriteFile(enginePath, engine, 0755); e != nil {\n\t\t\tmessageBox("Horizonte AFK", "Não foi possível preparar a conexão.")\n\t\t\treturn\n\t\t}\n\t}\n'''
if s.count(engine_block) != 1:
    raise SystemExit(f"engine materialization anchor count={s.count(engine_block)}")
lua_block = engine_block + '''\tluaPath := filepath.Join(root, "lua51.dll")\n\tlua51, e := payload.ReadFile("payload/lua51.dll")\n\tif e != nil {\n\t\tmessageBox("Horizonte AFK", "O runtime auxiliar da conexão não foi encontrado.")\n\t\treturn\n\t}\n\tluaSum := sha256.Sum256(lua51)\n\tif hex.EncodeToString(luaSum[:]) != expectedLuaSHA256 {\n\t\tmessageBox("Horizonte AFK", "A verificação do runtime auxiliar falhou.")\n\t\treturn\n\t}\n\tif fileSHA256(luaPath) != expectedLuaSHA256 {\n\t\tif e := os.WriteFile(luaPath, lua51, 0644); e != nil {\n\t\t\tmessageBox("Horizonte AFK", "Não foi possível preparar o runtime auxiliar.")\n\t\t\treturn\n\t\t}\n\t}\n'''
s = s.replace(engine_block, lua_block, 1)

required = [
    'appVersion           = "1.6.25-beta.1"',
    "runtime-v13-cloud-reconcile",
    "payload/lua51.dll",
    "expectedLuaSHA256",
    core_hash,
    lua_hash,
    "runUpdater(selfPath)",
]
for needle in required:
    if needle not in s:
        raise SystemExit(f"wrapper invariant missing: {needle}")

# The function itself remains available in features.go for build compatibility,
# but this wrapper must never call it against the already-complete 1.6.24 runtime.
if "ps = applyFeaturePatch(ps)" in s:
    raise SystemExit("wrapper would reapply legacy feature patch")

main.write_text(s, encoding="utf-8")
print(f"CORE_SHA256={core_hash}")
print(f"LUA_SHA256={lua_hash}")
