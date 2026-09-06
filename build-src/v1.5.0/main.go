package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// Network/dialog payload remains based on the validated Todos_Dialogos runtime.
//
//go:embed payload/Horizonte-AFK.ps1 payload/RakSAMPClient.exe payload/hero_bg.jpg payload/hz_icon.png
var payload embed.FS

const (
	appVersion           = "1.5.0"
	expectedEngineSHA256 = "b63de2af9dff1fbea4f79897054576788c57427637385ab92397999cdac9320c"
	manifestURL          = "https://raw.githubusercontent.com/guilhermedigv/Horizonte-AFK-Releases/main/latest.json"
)

func psQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

// Runs the full updater in a separate PowerShell process. It only shows the
// update screen after a real newer/hash-different build has been found.
// Exit code 42 means a self-replacement has been scheduled and this EXE must exit.
func runUpdater(selfPath string) bool {
	tempRoot, err := os.MkdirTemp("", "HorizonteAFK-Updater-")
	if err != nil {
		return false
	}
	script := makeUpdaterScript(selfPath, tempRoot, os.Getpid())
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-STA", "-WindowStyle", "Hidden", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	err = cmd.Run()
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 42 {
		return true
	}
	_ = os.RemoveAll(tempRoot)
	return false
}

func makeUpdaterScript(selfPath, tempRoot string, pid int) string {
	return fmt.Sprintf(`
$ErrorActionPreference='Stop'
Import-Module (Join-Path $PSHOME 'Modules\Microsoft.PowerShell.Utility\Microsoft.PowerShell.Utility.psd1') -Force
Add-Type -AssemblyName PresentationFramework
Add-Type -AssemblyName PresentationCore
Add-Type -AssemblyName WindowsBase
$manifestUrl=%s
$currentVersion=%s
$selfPath=%s
$tempRoot=%s
$goPid=%d
function V([string]$v){$c=($v-replace'^[vV]','').Split('-')[0];try{return [version]$c}catch{return [version]'0.0.0'}}
function Pump { try { [Windows.Threading.Dispatcher]::CurrentDispatcher.Invoke([action]{},[Windows.Threading.DispatcherPriority]::Background) } catch {} }
function Set-P([int]$p,[string]$s,[string]$d){if($script:bar){$script:bar.Value=$p;$script:pct.Text=($p.ToString()+'%%');$script:status.Text=$s;$script:detail.Text=$d;Pump}}
try {
  $headers=@{'Cache-Control'='no-cache';'User-Agent'='Horizonte-AFK/%s'}
  $m=Invoke-RestMethod -Uri ($manifestUrl+'?t='+[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()) -UseBasicParsing -TimeoutSec 7 -Headers $headers
  if(-not $m.available){exit 0}
  if(([string]$m.sha256)-notmatch'^[a-fA-F0-9]{64}$'){exit 0}
  if([string]::IsNullOrWhiteSpace([string]$m.version)){exit 0}
  $rv=V([string]$m.version);$lv=V($currentVersion)
  if($rv-lt$lv){exit 0}
  $lh=(Get-FileHash -LiteralPath $selfPath -Algorithm SHA256).Hash.ToLowerInvariant();$rh=([string]$m.sha256).ToLowerInvariant()
  if($rv-eq$lv-and$lh-eq$rh){exit 0}

  $x=@'
<Window xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation" xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml" Width="560" Height="300" WindowStartupLocation="CenterScreen" WindowStyle="None" ResizeMode="NoResize" AllowsTransparency="True" Background="Transparent" FontFamily="Segoe UI" Topmost="True">
 <Border CornerRadius="18" Background="#0B1119" BorderBrush="#243247" BorderThickness="1"><Grid><Grid.RowDefinitions><RowDefinition Height="58"/><RowDefinition Height="*"/></Grid.RowDefinitions>
  <Border Grid.Row="0" Background="#101722" CornerRadius="18,18,0,0"><StackPanel Orientation="Horizontal" Margin="18,0" VerticalAlignment="Center"><Border Width="36" Height="36" CornerRadius="10" Background="#129EFF"><TextBlock Text="HZ" Foreground="White" FontWeight="Bold" FontSize="15" HorizontalAlignment="Center" VerticalAlignment="Center"/></Border><StackPanel Margin="11,0,0,0" VerticalAlignment="Center"><TextBlock Text="HORIZONTE AFK" Foreground="White" FontWeight="Bold" FontSize="14"/><TextBlock Text="Atualização segura" Foreground="#728198" FontSize="9"/></StackPanel></StackPanel></Border>
  <StackPanel Grid.Row="1" Margin="28,27,28,25"><TextBlock Text="Atualizando o launcher" Foreground="White" FontWeight="Bold" FontSize="24"/><TextBlock x:Name="Status" Text="Preparando atualização" Foreground="#B7C4D5" FontSize="12" Margin="0,11,0,0"/><TextBlock x:Name="Detail" Text="Aguarde um instante..." Foreground="#74849A" FontSize="10" Margin="0,4,0,19" TextWrapping="Wrap"/><ProgressBar x:Name="Bar" Minimum="0" Maximum="100" Height="12" BorderThickness="0" Background="#111A26" Foreground="#18BFFF" Value="2"/><Grid Margin="0,11,0,0"><TextBlock Text="O mesmo Horizonte_AFK.exe será substituído com segurança." Foreground="#65758C" FontSize="9"/><TextBlock x:Name="Pct" Text="2%%" Foreground="#35C9FF" FontSize="10" FontWeight="Bold" HorizontalAlignment="Right"/></Grid></StackPanel>
 </Grid></Border>
</Window>
'@
  $r=New-Object System.Xml.XmlNodeReader([xml]$x);$script:w=[Windows.Markup.XamlReader]::Load($r);$script:bar=$w.FindName('Bar');$script:status=$w.FindName('Status');$script:detail=$w.FindName('Detail');$script:pct=$w.FindName('Pct');$w.Show();Pump
  Set-P 7 'Preparando atualização' ('Versão '+[string]$m.version+' encontrada.')

  New-Item -ItemType Directory -Path $tempRoot -Force|Out-Null
  $packed=Join-Path $tempRoot 'update.pkg';$new=Join-Path $tempRoot 'Horizonte_AFK.new.exe'
  $allowed=@('raw.githubusercontent.com','github.com','objects.githubusercontent.com')
  Set-P 12 'Baixando atualização' 'Recebendo os novos arquivos...'
  if(-not [string]::IsNullOrWhiteSpace([string]$m.download_url)){
    $u=[Uri]([string]$m.download_url);if($u.Scheme-ne'https'-or$u.Host-notin$allowed){throw'URL inválida'}
    Invoke-WebRequest -Uri ([string]$m.download_url) -UseBasicParsing -TimeoutSec 60 -OutFile $packed -Headers $headers
    Set-P 70 'Baixando atualização' 'Download concluído.'
  }elseif($m.parts-and @($m.parts).Count-gt0){
    $out=[IO.File]::Create($packed);try{$i=0;$total=@($m.parts).Count;foreach($pu0 in @($m.parts)){$i++;$u=[Uri]([string]$pu0);if($u.Scheme-ne'https'-or$u.Host-notin$allowed){throw'URL inválida'};$tmp=Join-Path $tempRoot ('part-'+$i+'.tmp');Invoke-WebRequest -Uri ([string]$pu0) -UseBasicParsing -TimeoutSec 60 -OutFile $tmp -Headers $headers;try{if(([string]$m.part_encoding).ToLowerInvariant()-eq'base64'){$bytes=[Convert]::FromBase64String(([IO.File]::ReadAllText($tmp)).Trim())}else{$bytes=[IO.File]::ReadAllBytes($tmp)};$out.Write($bytes,0,$bytes.Length)}finally{Remove-Item $tmp -Force -ErrorAction SilentlyContinue};Set-P (12+[int](58*$i/$total)) 'Baixando atualização' ('Recebendo arquivos... '+$i+'/'+$total)}}finally{$out.Dispose()}
  }else{throw'Manifesto sem download'}

  Set-P 76 'Preparando arquivos' 'Descompactando a nova versão...'
  if(([string]$m.encoding).ToLowerInvariant()-eq'gzip'){$inputStream=[IO.File]::OpenRead($packed);try{$gz=New-Object IO.Compression.GZipStream($inputStream,[IO.Compression.CompressionMode]::Decompress);try{$o=[IO.File]::Create($new);try{$gz.CopyTo($o)}finally{$o.Dispose()}}finally{$gz.Dispose()}}finally{$inputStream.Dispose()}}else{Copy-Item -LiteralPath $packed -Destination $new -Force}
  Set-P 86 'Validando atualização' 'Conferindo SHA-256 e integridade...'
  $h=(Get-FileHash -LiteralPath $new -Algorithm SHA256).Hash.ToLowerInvariant();if($h-ne$rh){throw'Hash inválido'}
  if($m.size_bytes-and (Get-Item -LiteralPath $new).Length-ne[int64]$m.size_bytes){throw'Tamanho inválido'}
  $fs=[IO.File]::OpenRead($new);try{$b0=$fs.ReadByte();$b1=$fs.ReadByte()}finally{$fs.Dispose()};if($b0-ne77-or$b1-ne90){throw'Executável inválido'}
  Set-P 94 'Aplicando atualização' 'Fechando a versão atual e preparando a substituição...'
  $rep=@'
$ErrorActionPreference='Stop'
Import-Module (Join-Path $PSHOME 'Modules\Microsoft.PowerShell.Utility\Microsoft.PowerShell.Utility.psd1') -Force
Wait-Process -Id __PID__ -ErrorAction SilentlyContinue
$new='__NEW__';$self='__SELF__';$temp='__TEMP__';$hash='__HASH__'
$stage=$self+'.hz-update';$backup=$self+'.hz-backup'
for($i=0;$i-lt60;$i++){
  try{
    if((Get-FileHash -LiteralPath $new -Algorithm SHA256).Hash.ToLowerInvariant()-ne$hash){throw 'Hash inválido'}
    Copy-Item -LiteralPath $new -Destination $stage -Force -ErrorAction Stop
    if((Get-FileHash -LiteralPath $stage -Algorithm SHA256).Hash.ToLowerInvariant()-ne$hash){throw 'Cópia inválida'}
    [IO.File]::Replace($stage,$self,$backup)
    break
  }catch{if($i-eq59){Remove-Item -LiteralPath $stage -Force -ErrorAction SilentlyContinue;exit 1};Start-Sleep -Milliseconds 250}
}
try { Start-Process -FilePath $self -ErrorAction Stop | Out-Null }
catch {
  if(Test-Path -LiteralPath $backup){[IO.File]::Replace($backup,$self,[NullString]::Value)}
  exit 1
}
Remove-Item -LiteralPath $backup -Force -ErrorAction SilentlyContinue
Remove-Item -LiteralPath $temp -Recurse -Force -ErrorAction SilentlyContinue
'@
  $rep=$rep.Replace('__HASH__',$rh).Replace('__PID__',[string]$goPid).Replace('__NEW__',($new-replace"'","''")).Replace('__SELF__',($selfPath-replace"'","''")).Replace('__TEMP__',($tempRoot-replace"'","''"))
  $repPath=Join-Path $tempRoot 'replace.ps1'
  [IO.File]::WriteAllText($repPath,$rep,(New-Object Text.UTF8Encoding($true)))
  Start-Process powershell.exe -WindowStyle Hidden -ArgumentList @('-NoProfile','-ExecutionPolicy','Bypass','-WindowStyle','Hidden','-File',('"'+$repPath+'"'))|Out-Null
  Set-P 100 'Atualização concluída' 'Reiniciando o Horizonte AFK...';Start-Sleep -Milliseconds 650;$w.Close();exit 42
}catch{
  try{if($script:w){Set-P 100 'Não foi possível atualizar' 'A versão instalada foi mantida. Tente novamente mais tarde.';Start-Sleep -Milliseconds 2200;$script:w.Close()}}catch{}
  exit 0
}
`, psQuote(manifestURL), psQuote(appVersion), psQuote(selfPath), psQuote(tempRoot), pid, appVersion)
}

func main() {
	selfPath, _ := os.Executable()
	selfPath, _ = filepath.Abs(selfPath)
	if selfPath != "" && runUpdater(selfPath) {
		return
	}

	for {
		marker := filepath.Join(os.TempDir(), fmt.Sprintf("HorizonteAFK-update-request-%d.txt", os.Getpid()))
		_ = os.Remove(marker)
		runPanel(selfPath, marker)
		if _, err := os.Stat(marker); err == nil {
			_ = os.Remove(marker)
			if runUpdater(selfPath) {
				return
			}
			continue
		}
		return
	}
}

func runPanel(selfPath, marker string) {
	local := os.Getenv("LOCALAPPDATA")
	if local == "" {
		if c, e := os.UserConfigDir(); e == nil {
			local = c
		} else {
			local = os.TempDir()
		}
	}
	root := filepath.Join(local, "HorizonteAFK", "runtime-v9-rpc-dialogs")
	if e := os.MkdirAll(root, 0755); e != nil {
		messageBox("Horizonte AFK", "Não foi possível preparar os arquivos internos.")
		return
	}
	ps, e := payload.ReadFile("payload/Horizonte-AFK.ps1")
	if e != nil {
		messageBox("Horizonte AFK", "O painel interno não foi encontrado.")
		return
	}
	ps = applyFeaturePatch(ps)
	hero, _ := payload.ReadFile("payload/hero_bg.jpg")
	icon, _ := payload.ReadFile("payload/hz_icon.png")
	psPath := filepath.Join(root, "Horizonte-AFK.ps1")
	enginePath := filepath.Join(root, "RakSAMPClient.exe")
	engine, e := payload.ReadFile("payload/RakSAMPClient.exe")
	if e != nil {
		messageBox("Horizonte AFK", "O componente de conexão não foi encontrado.")
		return
	}
	sum := sha256.Sum256(engine)
	if hex.EncodeToString(sum[:]) != expectedEngineSHA256 {
		messageBox("Horizonte AFK", "A verificação do componente interno falhou.")
		return
	}
	if fileSHA256(enginePath) != expectedEngineSHA256 {
		if e := os.WriteFile(enginePath, engine, 0755); e != nil {
			messageBox("Horizonte AFK", "Não foi possível preparar a conexão.")
			return
		}
	}
	if len(hero) > 0 {
		_ = os.WriteFile(filepath.Join(root, "hero_bg.jpg"), hero, 0644)
	}
	if len(icon) > 0 {
		_ = os.WriteFile(filepath.Join(root, "hz_icon.png"), icon, 0644)
	}
	if len(ps) < 3 || !(ps[0] == 0xEF && ps[1] == 0xBB && ps[2] == 0xBF) {
		ps = append([]byte{0xEF, 0xBB, 0xBF}, ps...)
	}
	if e := os.WriteFile(psPath, ps, 0644); e != nil {
		messageBox("Horizonte AFK", "Não foi possível preparar o painel.")
		return
	}
	selfHash := fileSHA256(selfPath)
	c := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-STA", "-WindowStyle", "Hidden", "-File", psPath)
	c.Env = append(os.Environ(), "HZ_UPDATE_REQUEST_PATH="+marker, "HZ_LAUNCHER_HASH="+selfHash, "HZ_LAUNCHER_VERSION="+appVersion)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	if e := c.Start(); e != nil {
		messageBox("Horizonte AFK", "Falha ao abrir o painel.")
		return
	}
	_ = c.Wait()
}

func fileSHA256(path string) string {
	b, e := os.ReadFile(path)
	if e != nil {
		return ""
	}
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}
func messageBox(title, text string) {
	script := fmt.Sprintf("Add-Type -AssemblyName PresentationFramework; [System.Windows.MessageBox]::Show(%s,%s)|Out-Null", psQuote(text), psQuote(title))
	c := exec.Command("powershell.exe", "-NoProfile", "-WindowStyle", "Hidden", "-Command", script)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	_ = c.Run()
}
