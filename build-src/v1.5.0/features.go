package main

import "strings"

// applyFeaturePatch keeps the validated network/dialog payload intact and only
// layers desktop UX features over the embedded PowerShell panel at runtime.
func applyFeaturePatch(ps []byte) []byte {
	s := strings.ReplaceAll(string(ps), "\r\n", "\n")
	s = replaceRequired(s, "$appVersion = '1.4.1'", "$appVersion = '1.5.0'")

	// Add maximize/full-screen and explicit exit controls while keeping the
	// existing minimize button untouched.
	s = replaceRequired(s, `                        <Button x:Name="CloseButton" Content="×" Width="42" Height="34"
                                Foreground="#C9D5E4" Background="Transparent" BorderThickness="0"
                                FontSize="22" Cursor="Hand"/>`, `                        <Button x:Name="MaximizeButton" Content="□" Width="42" Height="34"
                                Foreground="#C9D5E4" Background="Transparent" BorderThickness="0"
                                FontSize="17" Cursor="Hand" ToolTip="Tela cheia"/>
                        <Button x:Name="CloseButton" Content="×" Width="42" Height="34"
                                Foreground="#C9D5E4" Background="Transparent" BorderThickness="0"
                                FontSize="22" Cursor="Hand" ToolTip="Minimizar ao fechar"/>
                        <Button x:Name="ExitButton" Content="⏻" Width="42" Height="34"
                                Foreground="#FF9EAE" Background="Transparent" BorderThickness="0"
                                FontSize="16" Cursor="Hand" ToolTip="Sair completamente"/>`)

	s = replaceRequired(s, `                            <Button x:Name="ConnectButton" Content="CONECTAR E FICAR AFK  ▶"
                                    Style="{StaticResource PrimaryButton}" Margin="0,0,0,9"/>`, `                            <Button x:Name="ConnectButton" Content="CONECTAR E FICAR AFK  ▶"
                                    Style="{StaticResource PrimaryButton}" Margin="0,0,0,9"/>
                            <Button x:Name="MultiServerButton" Content="CONECTAR EM VÁRIOS SERVIDORES"
                                    Style="{StaticResource FlatButton}" Margin="0,0,0,9"
                                    BorderBrush="#1A6B93" Foreground="#8DDFFF"/>`)

	s = replaceRequired(s, `$MinimizeButton = $window.FindName('MinimizeButton')
$CloseButton = $window.FindName('CloseButton')`, `$MinimizeButton = $window.FindName('MinimizeButton')
$MaximizeButton = $window.FindName('MaximizeButton')
$CloseButton = $window.FindName('CloseButton')
$ExitButton = $window.FindName('ExitButton')`)

	s = replaceRequired(s, `$ConnectButton = $window.FindName('ConnectButton')
$StopButton = $window.FindName('StopButton')`, `$ConnectButton = $window.FindName('ConnectButton')
$MultiServerButton = $window.FindName('MultiServerButton')
$StopButton = $window.FindName('StopButton')`)

	multiFunctions := `
function Start-MultiServerSessions {
    if ($script:process -and !$script:process.HasExited) {
        Set-Status 'Sessão já ativa' 'Desconecte a sessão atual antes de abrir o modo multi-servidor.'
        return
    }
    if ([string]::IsNullOrWhiteSpace($NickBox.Text)) {
        Set-Status 'Informe o nick' 'O mesmo nick será usado nas sessões selecionadas.'
        return
    }

    $mx = @'
<Window xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
        xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml"
        Title="Multi-servidor" Width="470" Height="470"
        WindowStartupLocation="CenterOwner" ResizeMode="NoResize"
        Background="#081321" Foreground="White" FontFamily="Segoe UI">
    <Grid Margin="24">
        <Grid.RowDefinitions>
            <RowDefinition Height="Auto"/><RowDefinition Height="Auto"/>
            <RowDefinition Height="*"/><RowDefinition Height="Auto"/><RowDefinition Height="Auto"/>
        </Grid.RowDefinitions>
        <TextBlock Text="CONECTAR EM VÁRIOS SERVIDORES" FontSize="18" FontWeight="Bold" Foreground="#F4F8FD"/>
        <TextBlock Grid.Row="1" Text="Cada servidor abre uma sessão isolada. Login, 2FA e dialogs continuam separados para não misturar respostas."
                   Foreground="#8398B4" FontSize="10" TextWrapping="Wrap" Margin="0,7,0,18"/>
        <StackPanel Grid.Row="2">
            <CheckBox x:Name="S1" Content="Horizonte RP #1" Margin="0,0,0,13" FontSize="13"/>
            <CheckBox x:Name="S2" Content="Horizonte RP #2" Margin="0,0,0,13" FontSize="13"/>
            <CheckBox x:Name="S3" Content="Horizonte RP #3" Margin="0,0,0,13" FontSize="13"/>
            <CheckBox x:Name="S4" Content="Horizonte RP #4" Margin="0,0,0,13" FontSize="13"/>
            <CheckBox x:Name="S5" Content="Horizonte RP #5" Margin="0,0,0,13" FontSize="13"/>
        </StackPanel>
        <TextBlock x:Name="Info" Grid.Row="3" Foreground="#FFB36B" FontSize="10" Margin="0,0,0,12" TextWrapping="Wrap"/>
        <Grid Grid.Row="4">
            <Grid.ColumnDefinitions><ColumnDefinition Width="*"/><ColumnDefinition Width="10"/><ColumnDefinition Width="*"/></Grid.ColumnDefinitions>
            <Button x:Name="Start" Content="ABRIR SESSÕES" Padding="12,10" Background="#087CFF" Foreground="White" FontWeight="Bold"/>
            <Button x:Name="Cancel" Grid.Column="2" Content="CANCELAR" Padding="12,10" Background="#16243A" Foreground="#DCE8F5"/>
        </Grid>
    </Grid>
</Window>
'@
    $mr = New-Object System.Xml.XmlNodeReader([xml]$mx)
    $mw = [Windows.Markup.XamlReader]::Load($mr)
    $mw.Owner = $window
    $checks = @($mw.FindName('S1'),$mw.FindName('S2'),$mw.FindName('S3'),$mw.FindName('S4'),$mw.FindName('S5'))
    if ($ProfileBox.SelectedIndex -ge 0 -and $ProfileBox.SelectedIndex -lt 5) { $checks[$ProfileBox.SelectedIndex].IsChecked = $true }
    $info = $mw.FindName('Info')
    $start = $mw.FindName('Start')
    $cancel = $mw.FindName('Cancel')
    $script:multiPicked = @()

    $start.Add_Click({
        $picked = @()
        for ($i=0; $i -lt $checks.Count; $i++) {
            if ($checks[$i].IsChecked) { $picked += $i }
        }
        if ($picked.Count -lt 2) {
            $info.Text = 'Selecione pelo menos 2 servidores.'
            return
        }
        $script:multiPicked = @($picked)
        $mw.DialogResult = $true
    })
    $cancel.Add_Click({ $mw.DialogResult = $false })
    if ($mw.ShowDialog() -ne $true) { return }

    $started = 0
    foreach ($idx in @($script:multiPicked)) {
        try {
            $sessionParent = Join-Path $dataRoot ('multi-' + [Guid]::NewGuid().ToString('N'))
            $sessionBase = Join-Path $sessionParent 'runtime'
            New-Item -ItemType Directory -Path $sessionBase -Force | Out-Null
            Copy-Item -LiteralPath $exePath -Destination (Join-Path $sessionBase 'RakSAMPClient.exe') -Force
            $childScript = Join-Path $sessionBase 'Horizonte-AFK.ps1'
            Copy-Item -LiteralPath (Join-Path $base 'Horizonte-AFK.ps1') -Destination $childScript -Force
            if (Test-Path $heroPath) { Copy-Item -LiteralPath $heroPath -Destination (Join-Path $sessionBase 'hero_bg.jpg') -Force }
            if (Test-Path $iconPath) { Copy-Item -LiteralPath $iconPath -Destination (Join-Path $sessionBase 'hz_icon.png') -Force }

            $oldChild = $env:HZ_MULTI_CHILD
            $oldServer = $env:HZ_AUTO_SERVER
            $oldNick = $env:HZ_AUTO_NICK
            $oldUpdate = $env:HZ_UPDATE_REQUEST_PATH
            try {
                $env:HZ_MULTI_CHILD = '1'
                $env:HZ_AUTO_SERVER = [string]$idx
                $env:HZ_AUTO_NICK = $NickBox.Text.Trim()
                $env:HZ_UPDATE_REQUEST_PATH = ''
                Start-Process -FilePath 'powershell.exe' -WindowStyle Hidden -ArgumentList @(
                    '-NoProfile','-ExecutionPolicy','Bypass','-STA','-WindowStyle','Hidden','-File',('"' + $childScript + '"')
                ) | Out-Null
                $started++
            } finally {
                $env:HZ_MULTI_CHILD = $oldChild
                $env:HZ_AUTO_SERVER = $oldServer
                $env:HZ_AUTO_NICK = $oldNick
                $env:HZ_UPDATE_REQUEST_PATH = $oldUpdate
            }
        } catch {
            Add-SafeActivity ('Falha ao abrir uma sessão multi-servidor: ' + $_.Exception.Message)
        }
    }

    if ($started -gt 0) {
        Set-Status 'Multi-servidor iniciado' (($started.ToString()) + ' sessões independentes foram abertas.')
        Add-SafeActivity (($started.ToString()) + ' sessões multi-servidor iniciadas.')
    } else {
        Set-Status 'Falha no multi-servidor' 'Nenhuma sessão pôde ser iniciada.'
    }
}
`
	s = replaceRequired(s, "# Window chrome", multiFunctions+"\n# Window chrome")

	s = replaceRequired(s, `# Window chrome
$TopBar.Add_MouseLeftButtonDown({
    if ($_.ButtonState -eq [Windows.Input.MouseButtonState]::Pressed) {
        $window.DragMove()
    }
})
$MinimizeButton.Add_Click({ $window.WindowState = 'Minimized' })
$CloseButton.Add_Click({ Stop-Client; $window.Close() })
$window.Add_Closing({
    try { if ($script:updateTimer) { $script:updateTimer.Stop() } } catch {}
    if ($script:process -and !$script:process.HasExited) {
        try { $script:process.Kill() } catch {}
    }
    try { if (Test-Path $logPath) { Remove-Item $logPath -Force -ErrorAction SilentlyContinue } } catch {}
})`, `# Window chrome
$script:allowExit = $false
function Toggle-WindowSize {
    if ($window.WindowState -eq [Windows.WindowState]::Maximized) {
        $window.WindowState = [Windows.WindowState]::Normal
        $MaximizeButton.Content = '□'
        $MaximizeButton.ToolTip = 'Tela cheia'
    } else {
        $window.WindowState = [Windows.WindowState]::Maximized
        $MaximizeButton.Content = '❐'
        $MaximizeButton.ToolTip = 'Restaurar janela'
    }
}
$TopBar.Add_MouseLeftButtonDown({
    if ($_.ClickCount -ge 2) { Toggle-WindowSize; return }
    if ($_.ButtonState -eq [Windows.Input.MouseButtonState]::Pressed) {
        try { $window.DragMove() } catch {}
    }
})
$MinimizeButton.Add_Click({ $window.WindowState = [Windows.WindowState]::Minimized })
$MaximizeButton.Add_Click({ Toggle-WindowSize })
$CloseButton.Add_Click({ $window.WindowState = [Windows.WindowState]::Minimized })
$ExitButton.Add_Click({
    $script:allowExit = $true
    try { Stop-Client } catch {}
    $window.Close()
})
$window.Add_KeyDown({
    if ($_.Key -eq [Windows.Input.Key]::F11) {
        Toggle-WindowSize
        $_.Handled = $true
    }
})
$window.Add_Closing({
    param($sender,$e)
    if (!$script:allowExit) {
        $e.Cancel = $true
        $window.WindowState = [Windows.WindowState]::Minimized
        return
    }
    try { if ($script:updateTimer) { $script:updateTimer.Stop() } } catch {}
    if ($script:process -and !$script:process.HasExited) {
        try { $script:process.Kill() } catch {}
    }
    try { if (Test-Path $logPath) { Remove-Item $logPath -Force -ErrorAction SilentlyContinue } } catch {}
})`)

	s = replaceRequired(s, `$ConnectButton.Add_Click({
    try { Start-Client }`, `$MultiServerButton.Add_Click({ Start-MultiServerSessions })
$ConnectButton.Add_Click({
    try { Start-Client }`)

	s = replaceRequired(s, `[IO.File]::WriteAllText($updateRequestPath, 'update', (New-Object Text.UTF8Encoding($false)))
        $window.Close()`, `[IO.File]::WriteAllText($updateRequestPath, 'update', (New-Object Text.UTF8Encoding($false)))
        $script:allowExit = $true
        $window.Close()`)

	s = replaceRequired(s, `$script:updateTimer = New-Object Windows.Threading.DispatcherTimer
$script:updateTimer.Interval = [TimeSpan]::FromSeconds(60)
$script:updateTimer.Add_Tick({ Test-InAppUpdate })
$script:updateTimer.Start()

$window.ShowDialog() | Out-Null`, `$script:multiAutoStarted = $false
if ($env:HZ_MULTI_CHILD -eq '1') {
    $parsedServer = -1
    if ([int]::TryParse([string]$env:HZ_AUTO_SERVER, [ref]$parsedServer) -and $parsedServer -ge 0 -and $parsedServer -lt 5) {
        Select-ServerIndex $parsedServer
        $window.Title = ('Horizonte AFK - Servidor #' + ($parsedServer + 1))
    }
    if (![string]::IsNullOrWhiteSpace([string]$env:HZ_AUTO_NICK)) {
        $NickBox.Text = [string]$env:HZ_AUTO_NICK
    }
    $MultiServerButton.IsEnabled = $false
    $UpdateStatusText.Text = 'Atualizações gerenciadas pelo painel principal'
    $window.Add_ContentRendered({
        if (!$script:multiAutoStarted) {
            $script:multiAutoStarted = $true
            try { Start-Client }
            catch {
                Set-Status 'Não foi possível conectar' 'Verifique a sessão e tente novamente.'
                Add-SafeActivity 'Falha ao iniciar sessão multi-servidor.'
            }
        }
    })
} else {
    $script:updateTimer = New-Object Windows.Threading.DispatcherTimer
    $script:updateTimer.Interval = [TimeSpan]::FromSeconds(60)
    $script:updateTimer.Add_Tick({ Test-InAppUpdate })
    $script:updateTimer.Start()
}

$window.ShowDialog() | Out-Null
if ($env:HZ_MULTI_CHILD -eq '1') {
    try { Remove-Item -LiteralPath $dataRoot -Recurse -Force -ErrorAction SilentlyContinue } catch {}
}`)

	return []byte(s)
}

func replaceRequired(s, old, new string) string {
	if !strings.Contains(s, old) {
		panic("Horizonte AFK feature patch anchor not found")
	}
	return strings.Replace(s, old, new, 1)
}
