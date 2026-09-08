param([Parameter(Mandatory=$true)][string]$ExePath)
$bytes=[IO.File]::ReadAllBytes((Resolve-Path $ExePath))
$text=[Text.Encoding]::UTF8.GetString($bytes)
foreach($marker in @('1.6.54-beta.1','function Get-HZPageCloudSession','/v1/client/sessions/','CloudSessionId','Get-HZCloudSessionsSafe $force')){
    if(-not $text.Contains($marker)){throw "Missing marker: $marker"}
}
foreach($name in @('Get-HZPageCloudSession','Get-HZActivePageCloudSessionId','Sync-HZActivePageCloudAttachment','Get-HZConnectionStatusLabel','Get-HZConnectionChatText','Open-HZUcpConnection')){
    $start=$text.IndexOf('function '+$name)
    $end=$text.IndexOf("`nfunction ",$start+10)
    if($start-lt0-or$end-lt0){throw "Missing function: $name"}
    [void][scriptblock]::Create($text.Substring($start,$end-$start))
    Write-Host "Parse OK: $name"
}
$start=$text.IndexOf('function Get-HZPageCloudSession')
$end=$text.IndexOf("`nfunction ",$start+10)
Invoke-Expression $text.Substring($start,$end-$start)
$script:testAccount=[pscustomobject]@{Id='acc1';Nickname='bluey';Host='ip3.horizonte-rp.com';Port=7777;CloudSessionId='sid1'}
$page=[pscustomobject]@{AccountId='acc1'}
$script:cloudSessionId='sid1'
$script:activeAccountId='acc1'
$script:directCalls=0
function Get-HZPageAccount($value){return $script:testAccount}
function Get-HZAccountCloudSessionId($value){return [string]$value.CloudSessionId}
function Get-HZScalarInt($value,$default){return [int]$value}
function Get-HZCloudSessionsSafe([bool]$force){return @()}
function Register-HZDeviceCloudSession($account,$remote){}
function Set-HZAccountCloudProperty($account,[string]$name,$value){$account.$name=$value}
function Save-HZAccounts{}
function Save-HZDeviceRegistry{}
function Invoke-HZCloudApi($method,$path,$body,$timeout){
    $script:directCalls++
    return [pscustomobject]@{session_id='sid1';configured_nickname='bluey';host='ip3.horizonte-rp.com';port=7777;state='online';phase='spawned'}
}
1..5 | ForEach-Object {
    $remote=Get-HZPageCloudSession $page $false
    if($null-eq$remote-or[string]$remote.session_id-ne'sid1'){throw 'Live session lost during repeated refresh'}
}
if($script:directCalls-lt5){throw 'Direct live session resolver was not exercised'}
Write-Host 'Repeated live-session refresh simulation OK'
