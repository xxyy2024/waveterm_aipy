# 获取脚本所在目录的上级目录
$BASE_DIR = Resolve-Path (Join-Path $PSScriptRoot ".." -Resolve)


$APP_DIR = Join-Path "tools" "prompt-optimizer" "win"   

# 调试：打印路径
Write-Output "BASE_DIR: $BASE_DIR"
Write-Output "APP_DIR: $APP_DIR"


$APP_EXEC = Join-Path $baseDir $APP_DIR "app"

# 明确指定配置文件路径，例如，对于大多数应用，配置可能在用户目录下
$configFile = Join-Path $HOME ".prompt-optimizer" "config.yaml"
$sourceFile = Join-Path "tools" "prompt_config.yaml"

if (!(Test-Path $configFile)) {
    New-Item -ItemType Directory -Force -Path (Split-Path $configFile) | Out-Null
    Copy-Item $sourceFile $configFile
}


# 启动 Python 应用
& $$APP_EXEC  "interactive" "-c" $configFile