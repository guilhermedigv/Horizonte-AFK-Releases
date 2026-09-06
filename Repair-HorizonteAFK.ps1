param([Parameter(Mandatory=$true)][string]$TargetPath,[switch]$Start)
$ErrorActionPreference='Stop'
Import-Module (Join-Path $PSHOME 'Modules\Microsoft.PowerShell.Utility\Microsoft.PowerShell.Utility.psd1') -Force
$target=[IO.Path]::GetFullPath($TargetPath)
if([IO.Path]::GetExtension($target)-ne'.exe'){throw 'Informe o caminho do seu executável Horizonte AFK.'}
$parent=Split-Path -Parent $target
if(!(Test-Path -LiteralPath $parent -PathType Container)){throw 'A pasta informada não existe.'}
if(Get-Process | Where-Object { try{$_.Path-eq$target}catch{$false} }){throw 'Feche esse Horizonte AFK e execute a reparação novamente.'}
$manifest='https://raw.githubusercontent.com/guilhermedigv/Horizonte-AFK-Releases/main/latest.json'
$m=Invoke-RestMethod -Uri ($manifest+'?t='+[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()) -UseBasicParsing -TimeoutSec 15
if(!$m.available-or([string]$m.sha256)-notmatch'^[a-fA-F0-9]{64}$'){throw 'A atualização oficial ainda não está disponível.'}
$uri=[Uri]([string]$m.download_url)
if($uri.Scheme-ne'https'-or$uri.Host-notin@('github.com','raw.githubusercontent.com')){throw 'Endereço de atualização inválido.'}
$stage=Join-Path $parent ('.HorizonteAFK-'+[Guid]::NewGuid().ToString('N')+'.tmp')
$backup=$target+'.hz-backup'
try{
    Invoke-WebRequest -Uri $uri.AbsoluteUri -UseBasicParsing -TimeoutSec 120 -OutFile $stage
    if((Get-FileHash -LiteralPath $stage -Algorithm SHA256).Hash.ToLowerInvariant()-ne([string]$m.sha256).ToLowerInvariant()){throw 'O arquivo recebido falhou na verificação SHA-256.'}
    if((Get-Item -LiteralPath $stage).Length-ne[int64]$m.size_bytes){throw 'O tamanho recebido está incorreto.'}
    $s=[IO.File]::OpenRead($stage);try{if($s.ReadByte()-ne77-or$s.ReadByte()-ne90){throw 'O arquivo recebido não é um executável.'}}finally{$s.Dispose()}
    $existed=Test-Path -LiteralPath $target
    if($existed){[IO.File]::Replace($stage,$target,$backup)}else{Move-Item -LiteralPath $stage -Destination $target}
    if($Start){try{Start-Process -FilePath $target -ErrorAction Stop}catch{if($existed){[IO.File]::Replace($backup,$target,[NullString]::Value)};throw}}
    Remove-Item -LiteralPath $backup -Force -ErrorAction SilentlyContinue
    Write-Output ('Horizonte AFK '+$m.version+' instalado em '+$target)
}finally{Remove-Item -LiteralPath $stage -Force -ErrorAction SilentlyContinue}
