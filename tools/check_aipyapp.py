import importlib.util
import subprocess
import sys


def is_module_installed(module_name):
    return importlib.util.find_spec(module_name) is not None


def install_module(module_name):
    try:
        subprocess.check_call([sys.executable, "-m", "pip", "install", module_name,"--break-system-packages"])
        print(f"已成功安装 {module_name}。")
    except subprocess.CalledProcessError as e:
        print(f"安装 {module_name} 失败：{e}")
        sys.exit(1)


def main():
    # "promptify"
    for module_name in ["aipyapp"]:
        if is_module_installed(module_name):
            print(f"{module_name} 已安装。")
        else:
            print(f"未找到 {module_name}，开始安装...")
            install_module(module_name)


if __name__ == "__main__":
    main()