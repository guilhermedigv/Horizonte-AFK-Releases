$ErrorActionPreference='Stop'
$root=$PSScriptRoot;$baseline=Join-Path $root 'baseline';New-Item -ItemType Directory -Force $baseline|Out-Null;$baselineExe=Join-Path $baseline 'Horizonte_AFK.exe'
if(!(Test-Path $baselineExe)){Invoke-WebRequest -UseBasicParsing -Uri 'https://github.com/guilhermedigv/Horizonte-AFK-Releases/releases/download/v1.6.55-beta.1/Horizonte_AFK.exe' -OutFile $baselineExe}
if((Get-FileHash $baselineExe).Hash-ne'0079D132CA1963E1073C8914F4E321FE65F0894AE6A4875AB88DB8CBF7DC941C'){throw'Unexpected baseline SHA-256'}
$embedded=[Text.Encoding]::UTF8.GetString([IO.File]::ReadAllBytes($baselineExe));$start=$embedded.IndexOf('# Horizonte AFK Launcher');$end=$embedded.IndexOf('# PendingDialogText:',$start);if($start-lt0-or$end-lt0){throw'Baseline script boundaries missing'};$end=$embedded.IndexOf("`r`n",$end)
$utf8=New-Object Text.UTF8Encoding($false);$original=$embedded.Substring($start,$end-$start);$source=$original.Replace("`r`n","`n")
function Get-Fn([string]$name){$t=$null;$e=$null;$ast=[Management.Automation.Language.Parser]::ParseInput($script:source,[ref]$t,[ref]$e);$fn=$ast.Find({param($a)$a-is[Management.Automation.Language.FunctionDefinitionAst]-and$a.Name-eq$name},$false);if(!$fn){throw"Missing function $name"};return $fn}
function Replace-InFunction([string]$fnName,[string]$old,[string]$new,[string]$label){$fn=Get-Fn $fnName;$segment=$script:source.Substring($fn.Extent.StartOffset,$fn.Extent.EndOffset-$fn.Extent.StartOffset);$i=$segment.IndexOf($old,[StringComparison]::Ordinal);if($i-lt0){throw"Missing patch anchor $label in $fnName"};if($segment.IndexOf($old,$i+$old.Length,[StringComparison]::Ordinal)-ge0){throw"Non-unique patch anchor $label in $fnName"};$segment=$segment.Substring(0,$i)+$new+$segment.Substring($i+$old.Length);$script:source=$script:source.Substring(0,$fn.Extent.StartOffset)+$segment+$script:source.Substring($fn.Extent.EndOffset)}
function Insert-BeforeFunction([string]$name,[string]$text){$fn=Get-Fn $name;$script:source=$script:source.Substring(0,$fn.Extent.StartOffset)+$text+"`n"+$script:source.Substring($fn.Extent.StartOffset)}

$helpers=@'
$script:hzCloudDialogCoordinatorAt=[DateTime]::MinValue
$script:lastCloudDialogKey=''
function Get-HZAccountByCloudSessionId([string]$sid){if([string]::IsNullOrWhiteSpace($sid)){return $null};foreach($a in @($script:accounts)){try{if(([string](Get-HZAccountCloudSessionId $a)).Trim().Equals($sid,[StringComparison]::OrdinalIgnoreCase)){return $a}}catch{}};return $null}
function Set-HZCloudDialogOwner($d,[string]$sid){if($null-eq$d){return $null};$a=Get-HZAccountByCloudSessionId $sid;$aid=$(if($null-eq$a){''}else{[string]$a.Id});$sig='';try{$sig=[string]$d.Signature}catch{};if([string]::IsNullOrWhiteSpace($sig)){$sig=([string]$d.Generation+'|'+[string]$d.Seq+'|'+[string]$d.Id)};foreach($p in @(@('Source','cloud'),@('CloudSessionId',$sid),@('CloudAccountId',$aid),@('CloudDialogKey',($sid+'|'+$sig)))){try{$d|Add-Member -NotePropertyName $p[0] -NotePropertyValue $p[1] -Force}catch{}};return $d}
function Update-HZPendingCloudDialogs{
 if($env:HZ_MULTI_CHILD-eq'1'-or[string]::IsNullOrWhiteSpace([string]$script:cloudToken)){return};try{if($script:process-and!$script:process.HasExited){return}}catch{};if($DialogOverlay.Visibility-eq'Visible'-and$null-ne$script:activeDialogData){return}
 $now=Get-Date;if($script:hzCloudDialogCoordinatorAt-ne[DateTime]::MinValue-and($now-$script:hzCloudDialogCoordinatorAt).TotalMilliseconds-lt1200){return};$script:hzCloudDialogCoordinatorAt=$now;$sessions=@(Get-HZCloudSessionsSafe);if($sessions.Count-lt1){return};$current=([string]$script:cloudSessionId).Trim();$sessions=@($sessions|Sort-Object @{Expression={if(([string]$_.session_id).Equals($current,[StringComparison]::OrdinalIgnoreCase)){0}else{1}}})
 foreach($s in $sessions){if(!$s.enabled-or([string]$s.phase)-ne'awaiting_response'){continue};$sid=([string]$s.session_id).Trim();if([string]::IsNullOrWhiteSpace($sid)){continue};try{$x=Invoke-HZCloudApi 'GET' ('/v1/client/sessions/'+$sid+'/dialog') $null 5;if(!$x.pending-or$null-eq$x.dialog){continue};$d=Set-HZCloudDialogOwner $x.dialog $sid;if($null-eq$d){continue};Show-IntegratedDialog $d;$a=Get-HZAccountByCloudSessionId $sid;if($null-ne$a){Set-Status 'Dialog aguardando resposta' ('Conta '+[string]$a.Nickname+' • responda para continuar o login desta sessão.')}return}catch{}}
}
'@
Insert-BeforeFunction 'Initialize-HZAccountCenter' $helpers

Replace-InFunction 'Remember-HZAccountDialogResponse' '$account = Get-HZAccountById -id $script:activeAccountId' @'
$dialogAccountId=[string]$script:activeAccountId
    try{if($data.PSObject.Properties['CloudAccountId']-and![string]::IsNullOrWhiteSpace([string]$data.CloudAccountId)){$dialogAccountId=[string]$data.CloudAccountId}}catch{}
    $account = Get-HZAccountById -id $dialogAccountId
'@ 'remember-owner'
Replace-InFunction 'Try-HZQuickAnswerMainDialog' 'if ($signature -eq $script:lastQuickDialogSignature) { return }' @'
try{if($script:activeDialogData.PSObject.Properties['CloudSessionId']-and![string]::IsNullOrWhiteSpace([string]$script:activeDialogData.CloudSessionId)){$signature=([string]$script:activeDialogData.CloudSessionId)+'|'+$signature}}catch{}
    if ($signature -eq $script:lastQuickDialogSignature) { return }
'@ 'quick-signature'
Replace-InFunction 'Try-HZQuickAnswerMainDialog' 'foreach ($rule in @(Get-HZRulesForActiveAccount)) {' @'
$dialogRulesAccount=$null
    try{if($script:activeDialogData.PSObject.Properties['CloudAccountId']){$dialogRulesAccount=Get-HZAccountById -id ([string]$script:activeDialogData.CloudAccountId)}}catch{}
    if($null-eq$dialogRulesAccount){$dialogRulesAccount=Get-HZAccountById -id $script:activeAccountId}
    $dialogRules=$(if($null-eq$dialogRulesAccount){@(Get-HZRulesForActiveAccount)}else{@(Get-HZAccountRules $dialogRulesAccount)})
    foreach ($rule in @($dialogRules)) {
'@ 'quick-owner-rules'

Replace-InFunction 'Update-HZCloudSession' '$seq = [uint32]$d.Seq' @'
$seq = [uint32]$d.Seq
            $d=Set-HZCloudDialogOwner $d ([string]$script:cloudSessionId)
            $dialogKey=[string]$d.CloudDialogKey
'@ 'dialog-owner'
Replace-InFunction 'Update-HZCloudSession' 'if ($seq -ne $script:lastCloudDialogSeq) {' @'
$mayShow=$true
            if($DialogOverlay.Visibility-eq'Visible'-and$null-ne$script:activeDialogData){$shown='';try{$shown=[string]$script:activeDialogData.CloudDialogKey}catch{};if([string]::IsNullOrWhiteSpace($shown)-or!$shown.Equals($dialogKey,[StringComparison]::OrdinalIgnoreCase)){$mayShow=$false}}
            if ($mayShow -and $dialogKey -ne [string]$script:lastCloudDialogKey) {
'@ 'dialog-key-check'
Replace-InFunction 'Update-HZCloudSession' 'Show-IntegratedDialog $d' '$script:lastCloudDialogKey=$dialogKey; Show-IntegratedDialog $d' 'dialog-key-store'

Replace-InFunction 'Invoke-HZCloudDialogButton' 'if ([string]::IsNullOrWhiteSpace($script:cloudSessionId)) { return }' '' 'response-global-guard'
Replace-InFunction 'Invoke-HZCloudDialogButton' '$data = $script:activeDialogData' @'
$data = $script:activeDialogData
        $targetSid='';try{$targetSid=([string]$data.CloudSessionId).Trim()}catch{};if([string]::IsNullOrWhiteSpace($targetSid)){$targetSid=([string]$script:cloudSessionId).Trim()};if([string]::IsNullOrWhiteSpace($targetSid)){return}
'@ 'response-owner'
Replace-InFunction 'Invoke-HZCloudDialogButton' "('/v1/client/sessions/' + `$script:cloudSessionId + '/dialog-response')" "('/v1/client/sessions/' + `$targetSid + '/dialog-response')" 'response-url'

Replace-InFunction 'Start-HZAccountCloud' '$script:lastCloudDialogSeq = [uint32]0' @'
$script:lastCloudDialogSeq = [uint32]0
        $script:lastCloudDialogKey = ''
        try{Hide-IntegratedDialog}catch{}
'@ 'attach-reset'
Replace-InFunction 'Initialize-HZAccountCenter' '$script:accountAutoTimer.Add_Tick({try{Try-HZQuickAnswerMainDialog}catch{};try{Update-HZAccountLocalSessions}catch{}})' '$script:accountAutoTimer.Add_Tick({try{Update-HZPendingCloudDialogs}catch{};try{Try-HZQuickAnswerMainDialog}catch{};try{Update-HZAccountLocalSessions}catch{}})' 'coordinator-timer'

$source=$source.Replace('1.6.55-beta.1','1.6.56-beta.1');$source=[regex]::Replace($source,'(?m)^[ \t]+$','')
function To-Crlf([string]$s){return $s.Replace("`r`n","`n").Replace("`n","`r`n")}
$oldBytes=$utf8.GetBytes($original);$candidateText=To-Crlf $source;$newBytes=$utf8.GetBytes($candidateText)
if($newBytes.Length-gt$oldBytes.Length){$t=$null;$e=$null;[void][Management.Automation.Language.Parser]::ParseInput($source,[ref]$t,[ref]$e);foreach($c in @($t|Where-Object{$_.Kind-eq'Comment'-and$_.Text.Length-gt30}|Sort-Object{$_.Extent.StartOffset}-Descending)){if($utf8.GetByteCount((To-Crlf $source))-le$oldBytes.Length){break};$source=$source.Remove($c.Extent.StartOffset,$c.Extent.EndOffset-$c.Extent.StartOffset)};$candidateText=To-Crlf $source;$newBytes=$utf8.GetBytes($candidateText)}
if($newBytes.Length-gt$oldBytes.Length){throw"Script exceeds baseline capacity: $($newBytes.Length) > $($oldBytes.Length)"}
$t=$null;$e=$null;[void][Management.Automation.Language.Parser]::ParseInput($candidateText,[ref]$t,[ref]$e);if($e.Count){throw($e|Out-String)}
foreach($m in @('CloudDialogKey','CloudSessionId','CloudAccountId','Update-HZPendingCloudDialogs','Get-HZAccountByCloudSessionId','1.6.56-beta.1')){if(!$candidateText.Contains($m)){throw"Missing marker $m"}}
$candidate=Join-Path $root 'candidate';New-Item -ItemType Directory -Force $candidate|Out-Null;[IO.File]::WriteAllText((Join-Path $candidate 'runtime.ps1'),$candidateText,$utf8)
$bytes=[IO.File]::ReadAllBytes($baselineExe);$latin=[Text.Encoding]::GetEncoding(28591);$binary=$latin.GetString($bytes);$offset=$binary.IndexOf($latin.GetString($oldBytes),[StringComparison]::Ordinal);if($offset-lt0){throw'Exact embedded baseline not found'};for($i=0;$i-lt$oldBytes.Length;$i++){$bytes[$offset+$i]=32};[Array]::Copy($newBytes,0,$bytes,$offset,$newBytes.Length);$binary=$latin.GetString($bytes);$binary=$binary.Replace('1.6.55-beta.1','1.6.56-beta.1');[IO.File]::WriteAllBytes((Join-Path $candidate 'Horizonte_AFK.exe'),$latin.GetBytes($binary))
'PATCH offset='+$offset+' oldBytes='+$oldBytes.Length+' newBytes='+$newBytes.Length;Get-FileHash (Join-Path $candidate 'Horizonte_AFK.exe')
