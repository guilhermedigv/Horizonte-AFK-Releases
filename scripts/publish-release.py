"""Publish validated artifacts before activating latest.json."""
import hashlib, json, pathlib, subprocess, tempfile
ROOT = pathlib.Path(__file__).resolve().parent.parent
META = json.loads((ROOT / 'release-candidate.json').read_text(encoding='utf-8-sig'))
TAG = 'v' + META['version']
BETA = META.get('channel') == 'beta'
REPO = 'guilhermedigv/Horizonte-AFK-Releases'
EXE = ROOT / 'Horizonte_AFK.exe'
def gh(*args, check=True):
    return subprocess.run(['gh', *args, '--repo', REPO], check=check, text=True, capture_output=True)
data = EXE.read_bytes()
assert data[:2] == b'MZ', 'Not a Windows executable'
assert len(data) == META['size_bytes'], 'Executable size mismatch'
assert hashlib.sha256(data).hexdigest() == META['sha256'], 'SHA-256 mismatch'
assert META['available'] is True and META['encoding'] == 'none'
assert META['download_url'] == f'https://github.com/{REPO}/releases/download/{TAG}/Horizonte_AFK.exe'
view = gh('release', 'view', TAG, '--json', 'isDraft', check=False)
published = False
if view.returncode == 0:
    published = not json.loads(view.stdout)['isDraft']
else:
    gh('release', 'create', TAG, '--draft', '--target', 'main', '--title', 'Horizonte AFK ' + META['version'], '--notes-file', str(ROOT / 'RELEASE_NOTES.md'), *(['--prerelease'] if BETA else []))
if not published:
    gh('release', 'upload', TAG, str(EXE), str(ROOT / 'SHA256SUMS.txt'), str(ROOT / 'Repair-HorizonteAFK.ps1'), '--clobber')
with tempfile.TemporaryDirectory() as folder:
    gh('release', 'download', TAG, '--pattern', 'Horizonte_AFK.exe', '--dir', folder)
    assert hashlib.sha256(pathlib.Path(folder, 'Horizonte_AFK.exe').read_bytes()).hexdigest() == META['sha256'], 'Uploaded asset mismatch'
if not published:
    gh('release', 'edit', TAG, '--draft=false', *(['--prerelease', '--latest=false'] if BETA else ['--latest']))
print(f'Official release verified: https://github.com/{REPO}/releases/tag/{TAG}')
