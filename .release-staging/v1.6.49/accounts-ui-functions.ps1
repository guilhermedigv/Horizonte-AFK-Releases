function Refresh-HZAccountsCenter([bool]$forceCloud=$false){
 if($null-eq$script:ucpList){return}
 $selectedId='';try{$selectedId=[string]$script:ucpList.SelectedItem.AccountId}catch{}
 Load-HZAccounts;$cloud=@(Get-HZCloudSessionsSafe $forceCloud);Sync-HZAccountsFromCloud $cloud
 $script:ucpProfiles=@($script:accounts|Sort-Object LastUsed -Descending);$views=@();$onlineCount=0
 foreach($account in @($script:ucpProfiles)){
  $status=Get-HZUcpStatus $account $cloud;if($status-ne'OFFLINE'){$onlineCount++}
  $ruleCount=@(Get-HZAccountRules $account).Count;$idx=Get-HZScalarInt $account.ServerIndex -1
  $serverName=$(if($idx-ge0){'Horizonte RP #'+($idx+1)}else{[string]$account.Host})
  $bg='#182536';$fg='#A8C1D5'
  if($status-match'BLOQUEADA|RECUSADO|ERRO|FALHA'){$bg='#40202B';$fg='#FFB7C3'}
  elseif($status-match'AGUARDANDO|CONECTANDO|ENTRANDO'){$bg='#3A321C';$fg='#FFD778'}
  elseif($status-ne'OFFLINE' -and $status-notmatch'OFFLINE'){$bg='#123B30';$fg='#69F0B0'}
  $nick=([string]$account.Nickname).Trim();$initial=$(if($nick.Length-gt0){$nick.Substring(0,1).ToUpperInvariant()}else{'?'})
  $views += [pscustomobject]@{AccountId=[string]$account.Id;Status=$status;StatusBg=$bg;StatusFg=$fg;Nickname=$nick;Initial=$initial;Server=$serverName;Rules=('login rápido: '+$ruleCount);RuleHint=$(if($ruleCount-eq1){'1 resposta aprendida'}else{$ruleCount.ToString()+' respostas aprendidas'})}
 }
 $script:ucpList.ItemsSource=$null;$script:ucpList.ItemsSource=@($views)
 if($script:ucpSummary){$script:ucpSummary.Text=$script:ucpProfiles.Count.ToString()+' contas salvas • '+$onlineCount.ToString()+' ativas'}
 $found=$false;if(![string]::IsNullOrWhiteSpace($selectedId)){for($i=0;$i-lt$views.Count;$i++){if([string]$views[$i].AccountId-eq$selectedId){$script:ucpList.SelectedIndex=$i;$found=$true;break}}}
 if(!$found-and$views.Count-gt0){$script:ucpList.SelectedIndex=0}
}

function Get-HZUcpSelected{
 if($null-eq$script:ucpList-or$null-eq$script:ucpList.SelectedItem){return $null}
 $id='';try{$id=([string]$script:ucpList.SelectedItem.AccountId).Trim()}catch{}
 if([string]::IsNullOrWhiteSpace($id)){return $null};return Get-HZAccountById -id $id
}

function Open-HZUcpConnection([string]$id){
 $account=Get-HZAccountById -id $id;if($null-eq$account){throw'Esta conta não está mais disponível.'}
 $idx=Get-HZScalarInt $account.ServerIndex -1;$nick=([string]$account.Nickname).Trim();$page=$null
 foreach($p in @($script:hzAccountPages)){if([string]$p.AccountId-eq$id){$page=$p;break}}
 if($null-eq$page){foreach($p in @($script:hzAccountPages)){if(([string]$p.Nick).Equals($nick,[StringComparison]::OrdinalIgnoreCase)-and[int]$p.ServerIndex-eq$idx){$page=$p;break}}}
 if($null-eq$page){foreach($p in @($script:hzAccountPages)){if([string]::IsNullOrWhiteSpace([string]$p.AccountId)-and[string]::IsNullOrWhiteSpace([string]$p.Nick)){$page=$p;break}}}
 if($null-eq$page-and@($script:hzAccountPages).Count-lt5){New-HZAccountPage $false $nick $idx $id $false|Out-Null;foreach($p in @($script:hzAccountPages)){if([string]$p.AccountId-eq$id){$page=$p;break}}}
 if($null-eq$page){throw'As 5 páginas já estão ocupadas. Libere uma página para abrir esta conta.'}
 $page.AccountId=$id;$page.Nick=$nick;$page.ServerIndex=$idx;Save-HZAccountPages;Select-HZAccountPage ([string]$page.Id);Show-HZConnectionCenter
}

function Show-HZAccountsCenter{
 if($env:HZ_MULTI_CHILD-eq'1'){return};$script:ucpOpenConnectionId=''
 $ux=@'
<Window xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation" xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml" Title="Centro de Contas" Width="1040" Height="670" MinWidth="900" MinHeight="570" WindowStartupLocation="CenterOwner" Background="#07111D" Foreground="White" FontFamily="Segoe UI">
 <Window.Resources>
  <Style TargetType="ListBoxItem"><Setter Property="HorizontalContentAlignment" Value="Stretch"/><Setter Property="Background" Value="Transparent"/><Setter Property="BorderThickness" Value="0"/><Setter Property="Margin" Value="0,0,0,8"/><Setter Property="Padding" Value="0"/><Setter Property="Template"><Setter.Value><ControlTemplate TargetType="ListBoxItem"><Border x:Name="Card" Background="#0A1826" BorderBrush="#183A52" BorderThickness="1" CornerRadius="11"><ContentPresenter/></Border><ControlTemplate.Triggers><Trigger Property="IsSelected" Value="True"><Setter TargetName="Card" Property="Background" Value="#0C2B45"/><Setter TargetName="Card" Property="BorderBrush" Value="#0A9EFF"/></Trigger><Trigger Property="IsMouseOver" Value="True"><Setter TargetName="Card" Property="BorderBrush" Value="#24709A"/></Trigger></ControlTemplate.Triggers></ControlTemplate></Setter.Value></Setter></Style>
  <DataTemplate x:Key="AccountRow"><Grid Margin="12,10"><Grid.ColumnDefinitions><ColumnDefinition Width="170"/><ColumnDefinition Width="54"/><ColumnDefinition Width="150"/><ColumnDefinition Width="190"/><ColumnDefinition Width="165"/><ColumnDefinition Width="*"/></Grid.ColumnDefinitions><Border Background="{Binding StatusBg}" CornerRadius="13" Padding="10,5" HorizontalAlignment="Left" VerticalAlignment="Center"><StackPanel Orientation="Horizontal"><Ellipse Width="7" Height="7" Fill="{Binding StatusFg}" Margin="0,0,7,0"/><TextBlock Text="{Binding Status}" Foreground="{Binding StatusFg}" FontSize="10" FontWeight="Bold"/></StackPanel></Border><Border Grid.Column="1" Width="38" Height="38" CornerRadius="19" Background="#132C43" BorderBrush="#2C5C82" BorderThickness="1"><TextBlock Text="{Binding Initial}" Foreground="#D8EEFF" FontSize="16" FontWeight="Bold" HorizontalAlignment="Center" VerticalAlignment="Center"/></Border><StackPanel Grid.Column="2" VerticalAlignment="Center"><TextBlock Text="{Binding Nickname}" Foreground="#F3F8FC" FontSize="13" FontWeight="SemiBold"/><TextBlock Text="Conta salva" Foreground="#6689A2" FontSize="9" Margin="0,2,0,0"/></StackPanel><StackPanel Grid.Column="3" VerticalAlignment="Center"><TextBlock Text="{Binding Server}" Foreground="#DCEAF6" FontSize="12"/><TextBlock Text="Servidor" Foreground="#6689A2" FontSize="9" Margin="0,2,0,0"/></StackPanel><StackPanel Grid.Column="4" VerticalAlignment="Center"><TextBlock Text="{Binding Rules}" Foreground="#DCEAF6" FontSize="11"/><TextBlock Text="{Binding RuleHint}" Foreground="#6689A2" FontSize="9" Margin="0,2,0,0"/></StackPanel><Button x:Name="OpenChat" Grid.Column="5" Tag="{Binding AccountId}" Content="CONEXÕES E CHAT" Background="#087FE0" Foreground="White" BorderBrush="#18A9FF" FontWeight="Bold" FontSize="10" Padding="15,10" HorizontalAlignment="Right" VerticalAlignment="Center"/></Grid></DataTemplate>
 </Window.Resources>
 <Grid Margin="24"><Grid.RowDefinitions><RowDefinition Height="Auto"/><RowDefinition Height="Auto"/><RowDefinition Height="*"/><RowDefinition Height="Auto"/></Grid.RowDefinitions>
  <Grid><Grid.ColumnDefinitions><ColumnDefinition Width="*"/><ColumnDefinition Width="Auto"/></Grid.ColumnDefinitions><StackPanel><TextBlock Text="CONTAS E SESSÕES" FontSize="25" FontWeight="Bold"/><TextBlock Text="Sessões 24/7 na Discloud • selecione uma conta para gerenciar ou abrir o chat" Foreground="#7895AA" FontSize="11" Margin="0,5,0,0"/></StackPanel><Button x:Name="Refresh" Grid.Column="1" Content="↻  ATUALIZAR" Background="#10263A" Foreground="#BFEAFF" BorderBrush="#285A76" Padding="15,9" VerticalAlignment="Center"/></Grid>
  <Grid Grid.Row="1" Margin="0,18,0,11"><Grid.ColumnDefinitions><ColumnDefinition Width="*"/><ColumnDefinition Width="Auto"/></Grid.ColumnDefinitions><TextBlock x:Name="Summary" Text="Carregando..." Foreground="#64DDA2" FontWeight="SemiBold" VerticalAlignment="Center"/><TextBlock Grid.Column="1" Text="Acesso rápido por conta • status atualizado automaticamente" Foreground="#617F95" FontSize="9" VerticalAlignment="Center"/></Grid>
  <Border Grid.Row="2" Background="#081521" BorderBrush="#18364C" BorderThickness="1" CornerRadius="12" Padding="8"><ListBox x:Name="Accounts" Background="Transparent" BorderThickness="0" ItemTemplate="{StaticResource AccountRow}" ScrollViewer.HorizontalScrollBarVisibility="Disabled"/></Border>
  <Grid Grid.Row="3" Margin="0,15,0,0"><Grid.ColumnDefinitions><ColumnDefinition Width="Auto"/><ColumnDefinition Width="Auto"/><ColumnDefinition Width="Auto"/><ColumnDefinition Width="*"/><ColumnDefinition Width="Auto"/></Grid.ColumnDefinitions><Button x:Name="SaveCurrent" Content="SALVAR CONTA ATUAL" Background="#17273A" Foreground="#DCEBFA" BorderBrush="#31506D" Padding="13,9" Margin="0,0,8,0"/><Button x:Name="Cloud" Grid.Column="1" Content="☁  CONECTAR 24/7" Background="#103C2D" Foreground="#79F0B0" BorderBrush="#23895D" Padding="13,9" Margin="0,0,8,0"/><Button x:Name="Stop" Grid.Column="2" Content="DESCONECTAR" Background="#422028" Foreground="#FFC3CA" BorderBrush="#71313C" Padding="13,9"/><Button x:Name="Delete" Grid.Column="4" Content="EXCLUIR PERFIL" Background="#251A20" Foreground="#E1AEBB" BorderBrush="#633541" Padding="13,9"/></Grid>
 </Grid>
</Window>
'@
 $reader=New-Object System.Xml.XmlNodeReader([xml]$ux);$dialog=[Windows.Markup.XamlReader]::Load($reader);$dialog.Owner=$window;$script:ucpWindow=$dialog;$script:ucpList=$dialog.FindName('Accounts');$script:ucpSummary=$dialog.FindName('Summary')
 $dialog.AddHandler([Windows.Controls.Button]::ClickEvent,[Windows.RoutedEventHandler]{param($sender,$e);$b=$e.OriginalSource;if($b-is[Windows.Controls.Button]-and$b.Name-eq'OpenChat'){$script:ucpOpenConnectionId=[string]$b.Tag;try{$script:ucpWindow.Close()}catch{};$e.Handled=$true}})
 $dialog.FindName('Refresh').Add_Click({Refresh-HZAccountsCenter $true})
 $dialog.FindName('SaveCurrent').Add_Click({try{$null=Ensure-HZAccountCurrent;Refresh-HZAccountsCenter}catch{}})
 $dialog.FindName('Cloud').Add_Click({try{$account=Get-HZUcpSelected;if($null-eq$account){return};$ruleCount=@(Get-HZAccountRules $account).Count;$mainBusy=$false;try{$mainBusy=($script:process-and!$script:process.HasExited)}catch{};if($mainBusy-and$ruleCount-lt1){throw'Faça o primeiro login desta conta sem outra sessão local ativa para o launcher aprender os dialogs.'};Start-HZAccountCloud $account (!$mainBusy);Refresh-HZAccountsCenter $true}catch{[Windows.MessageBox]::Show($_.Exception.Message,'Horizonte AFK')|Out-Null}})
 $dialog.FindName('Stop').Add_Click({try{$account=Get-HZUcpSelected;if($null-eq$account){return};$sid=Get-HZAccountCloudSessionId $account;foreach($cloud in @(Get-HZCloudSessionsSafe)){if([string]$cloud.session_id-eq$sid){$null=Invoke-HZCloudApi 'POST' ('/v1/client/sessions/'+$sid+'/stop') @{confirm=$true} 8}};Stop-HZAccountLocalSession -id ([string]$account.Id);Stop-HZAdditiveMultiSessionForAccount -id ([string]$account.Id);Stop-HZPageSessionForAccount -accountId ([string]$account.Id);try{if([string]$script:activeAccountId-eq[string]$account.Id-and$script:process-and!$script:process.HasExited){Stop-Client}}catch{};if([string]$script:cloudSessionId-eq$sid){Detach-HZCloudUi};Refresh-HZAccountsCenter $true}catch{}})
 $dialog.FindName('Delete').Add_Click({$account=Get-HZUcpSelected;if($null-eq$account){return};$answer=[Windows.MessageBox]::Show(('Excluir o perfil salvo de '+[string]$account.Nickname+'?'),'Horizonte AFK',[Windows.MessageBoxButton]::YesNo);if($answer-ne[Windows.MessageBoxResult]::Yes){return};$id=[string]$account.Id;$script:accounts=@($script:accounts|Where-Object{[string]$_.Id-ne$id});Save-HZAccounts;Save-HZDeviceRegistry;Refresh-HZAccountsCenter})
 Refresh-HZAccountsCenter $true;if($script:ucpRefreshTimer){try{$script:ucpRefreshTimer.Stop()}catch{}};$script:ucpRefreshTimer=New-Object Windows.Threading.DispatcherTimer;$script:ucpRefreshTimer.Interval=[TimeSpan]::FromSeconds(2);$script:ucpRefreshTimer.Add_Tick({try{Refresh-HZAccountsCenter $false}catch{}});$script:ucpRefreshTimer.Start();$null=$dialog.ShowDialog()
 if($script:ucpRefreshTimer){try{$script:ucpRefreshTimer.Stop()}catch{};$script:ucpRefreshTimer=$null};$openId=[string]$script:ucpOpenConnectionId;$script:ucpWindow=$null;$script:ucpList=$null;$script:ucpSummary=$null;$script:ucpProfiles=@();$script:ucpOpenConnectionId='';if(![string]::IsNullOrWhiteSpace($openId)){try{Open-HZUcpConnection $openId}catch{[Windows.MessageBox]::Show($_.Exception.Message,'Horizonte AFK')|Out-Null}}
}

