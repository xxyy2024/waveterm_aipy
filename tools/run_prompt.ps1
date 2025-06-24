# 获取脚本所在目录的上级目录
$BASE_DIR = Resolve-Path (Join-Path $PSScriptRoot ".." -Resolve)

$APP_DIR = [IO.Path]::Combine("tools", "prompt-optimizer", "win")
$APP_EXEC = [IO.Path]::Combine($BASE_DIR, $APP_DIR, "app.exe")
$configFile = [IO.Path]::Combine($HOME, ".prompt-optimizer", "config.yaml")

# 调试：打印路径
Write-Output "BASE_DIR: $BASE_DIR"
Write-Output "APP_DIR: $APP_DIR"





$sourceFile = [IO.Path]::Combine($BASE_DIR, "tools", "prompt_config.yaml")

if (!(Test-Path $configFile)) {
	New-Item -ItemType Directory -Force -Path (Split-Path $configFile) | Out-Null
	Copy-Item $sourceFile $configFile
}


# 启动 Python 应用
& $APP_EXEC  "interactive" "-c" $configFile