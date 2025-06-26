# 获取脚本所在目录的上级目录
$BASE_DIR = Resolve-Path (Join-Path $PSScriptRoot ".." -Resolve)

# 设置 Python 嵌入目录
$PYTHON_DIR = [IO.Path]::Combine("tools", "python-embed")  # 使用 / 确保 PowerShell 兼容性




# 构建 Python 解释器路径
$PYTHON_EXEC = [IO.Path]::Combine( $BASE_DIR, $PYTHON_DIR, "python.exe")

# 验证 Python 解释器
if (-not (Test-Path $PYTHON_EXEC)) {
	Write-Error "Error: python.exe not found at $PYTHON_EXEC"
	exit 1
}

# 设置 PYTHONPATH
$Env:PYTHONPATH = [IO.Path]::Combine($BASE_DIR, $PYTHON_DIR, "Lib", "site-packages")
Write-Output "PYTHONPATH: $Env:PYTHONPATH"

$Env:PYTHONHOME = [IO.Path]::Combine($BASE_DIR, $PYTHON_DIR)
Write-Output "PYTHONHOME: $Env:PYTHONHOME"

echo "Start AiPy..."
# 启动 Python 应用
& $PYTHON_EXEC -c "from aipyapp.__main__ import main; main()"

