import sys
import importlib.util

def check_library(lib_name):
    return importlib.util.find_spec(lib_name) is not None

def install_libraries(libs):
    import subprocess
    pip_cmd = [sys.executable, "-m", "pip", "install", "--no-warn-script-location"] + libs
    try:
        result = subprocess.run(pip_cmd, capture_output=True, text=True, check=True)
        print(result.stdout)
        return {"success": True, "message": result.stdout}
    except subprocess.CalledProcessError as e:
        print(e.stderr, file=sys.stderr)
        return {"success": False, "message": e.stderr}

if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "--check":
        if check_library("aipyapp"):
            print("aipyapp 已安装")
            sys.exit(0)
        else:
            print("aipyapp 未安装")
            sys.exit(1)
    else:
        libraries = ["aipyapp"]  # 动态安装其他库
        result = install_libraries(libraries)
        sys.exit(0 if result["success"] else 1)