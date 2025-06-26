<p align="center">
  <a href="https://www.waveterm.dev">
	<picture>
		<source media="(prefers-color-scheme: dark)" srcset="./assets/wave-dark.png">
		<source media="(prefers-color-scheme: light)" srcset="./assets/wave-light.png">
		<img alt="Wave Terminal Logo" src="./assets/wave-light.png" width="240">
	</picture>
  </a>
  <br/>
</p>

# Wave Terminal

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwavetermdev%2Fwaveterm.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwavetermdev%2Fwaveterm?ref=badge_shield)

Wave is an open-source terminal that combines traditional terminal features with graphical capabilities like file previews, web browsing, and AI assistance. It runs on MacOS, Linux, and Windows.

Modern development involves constantly switching between terminals and browsers - checking documentation, previewing files, monitoring systems, and using AI tools. Wave brings these graphical tools directly into the terminal, letting you control them from the command line. This means you can stay in your terminal workflow while still having access to the visual interfaces you need.

![WaveTerm Screenshot](./assets/wave-screenshot.webp)

## Key Features

- Flexible drag & drop interface to organize terminal blocks, editors, web browsers, and AI assistants
- Built-in editor for seamlessly editing remote files with syntax highlighting and modern editor features
- Rich file preview system for remote files (markdown, images, video, PDFs, CSVs, directories)
- Integrated AI chat with support for multiple models (OpenAI, Claude, Azure, Perplexity, Ollama)
- Command Blocks for isolating and monitoring individual commands with auto-close options
- One-click remote connections with full terminal and file system access
- Rich customization including tab themes, terminal styles, and background images
- Powerful `wsh` command system for managing your workspace from the CLI and sharing data between terminal sessions

## Installation

Wave Terminal works on macOS, Linux, and Windows.

Platform-specific installation instructions can be found [here](https://docs.waveterm.dev/gettingstarted).

You can also install Wave Terminal directly from: [www.waveterm.dev/download](https://www.waveterm.dev/download).

### Minimum requirements

Wave Terminal runs on the following platforms:

- macOS 11 or later (arm64, x64)
- Windows 10 1809 or later (x64)
- Linux based on glibc-2.28 or later (Debian 10, RHEL 8, Ubuntu 20.04, etc.) (arm64, x64)

The WSH helper runs on the following platforms:

- macOS 11 or later (arm64, x64)
- Windows 10 or later (arm64, x64)
- Linux Kernel 2.6.32 or later (x64), Linux Kernel 3.1 or later (arm64)

## Roadmap

Wave is constantly improving! Our roadmap will be continuously updated with our goals for each release. You can find it [here](./ROADMAP.md).

Want to provide input to our future releases? Connect with us on [Discord](https://discord.gg/XfvZ334gwU) or open a [Feature Request](https://github.com/wavetermdev/waveterm/issues/new/choose)!

## Links

- Homepage &mdash; https://www.waveterm.dev
- Download Page &mdash; https://www.waveterm.dev/download
- Documentation &mdash; https://docs.waveterm.dev
- Legacy Documentation &mdash; https://legacydocs.waveterm.dev
- Blog &mdash; https://blog.waveterm.dev
- X &mdash; https://x.com/wavetermdev
- Discord Community &mdash; https://discord.gg/XfvZ334gwU

## Building from Source

See [Building Wave Terminal](BUILD.md).

## Contributing

Wave uses GitHub Issues for issue tracking.

Find more information in our [Contributions Guide](CONTRIBUTING.md), which includes:

- [Ways to contribute](CONTRIBUTING.md#contributing-to-wave-terminal)
- [Contribution guidelines](CONTRIBUTING.md#before-you-start)
- [Storybook](https://docs.waveterm.dev/storybook)

## License

Wave Terminal is licensed under the Apache-2.0 License. For more information on our dependencies, see [here](./ACKNOWLEDGEMENTS.md).

## Python 解释器

```bash

brew install python@3.11
python3.11 -m venv python-embed-mac-arm64
./python-embed-mac-arm64/bin/python3 -m ensurepip --upgrade
./python-embed-mac-arm64/bin/pip install aipyapp
cd python-embed-mac-arm64
rm -rf ./lib/python3.11/test
rm -rf ./lib/python3.11/pycache
# copy python-embed-mac-arm64 to target

b##  uild

task package  -- --mac --x64
task package  -- --mac --arm64
```

## 窗口的Tab源码

pkg/wcore/layout.go

## 实现本地图标支持

为了让 widgets.json 的 icon 字段支持本地图片路径（如 icons/terminal.png），需要修改 Widget 组件的渲染逻辑，并可能调整 makeIconClass 函数。以下是具体步骤：
a. 准备本地图标
创建图标文件：
准备 PNG 或 SVG 图标，建议尺寸为 20x20 或 24x24 像素，适配 widget 栏。

存储到 ./public/icons/：
bash

mkdir -p ./public/icons
cp /path/to/terminal.png ./public/icons/

配置 Vite 复制静态资源：
编辑 ./electron.vite.config.ts：
typescript

import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import { viteStaticCopy } from 'vite-plugin-static-copy';

export default defineConfig({
plugins: [
react(),
viteStaticCopy({
targets: [
{ src: 'public/fontawesome/webfonts/*', dest: 'fontawesome/webfonts' },
{ src: 'public/icons/*', dest: 'icons' },
],
}),
],
});

说明：确保 ./public/icons/terminal.png 复制到 ./dist/frontend/icons/terminal.png，运行时通过 /icons/terminal.png 访问。

更新 widgets.json：
修改 ~/.config/waveterm/config/widgets.json：
json

{
"defwidget@terminal": {
"display:order": -5,
"icon": "icons/terminal.png",
"label": "terminal",
"blockdef": {
"meta": {
"view": "term",
"controller": "shell"
}
}
}
}

b. 修改 Widget 组件
修改 ./frontend/app/workspace/workspace.tsx 中的 Widget 组件，支持图片路径渲染：
更新 Widget 组件：
tsx

const Widget = memo(({ widget }: { widget: WidgetConfigType }) => {
const isImage = widget.icon && (widget.icon.endsWith('.png') || widget.icon.endsWith('.svg'));

    return (
        <div
            className={clsx(
                "flex flex-col justify-center items-center w-full py-1.5 pr-0.5 text-secondary text-lg overflow-hidden rounded-sm hover:bg-hoverbg hover:text-white cursor-pointer",
                widget["display:hidden"] && "hidden"
            )}
            onClick={() => handleWidgetSelect(widget)}
            title={widget.description || widget.label}
        >
            <div style={{ color: widget.color }}>
                {isImage ? (
                    <img
                        src={`/${widget.icon}`}
                        alt={widget.label}
                        style={{ width: '20px', height: '20px', objectFit: 'contain' }}
                    />
                ) : (
                    <i className={makeIconClass(widget.icon, true, { defaultIcon: "browser" })}></i>
                )}
            </div>
            {!isBlank(widget.label) ? (
                <div className="text-xxs mt-0.5 w-full px-0.5 text-center whitespace-nowrap overflow-hidden">
                    {widget.label}
                </div>
            ) : null}
        </div>
    );

});

说明：
添加 isImage 检查，判断 widget.icon 是否以 .png 或 .svg 结尾。

如果是图片路径，使用 <img> 标签渲染；否则，使用 <i> 渲染 Font Awesome 图标。

src={/${widget.icon}} 假设图片在 ./dist/frontend/icons/（例如 /icons/terminal.png）。

保留 widget.color 应用到 <div>，确保图片颜色与 Font Awesome 图标一致（如果需要，可直接应用到 <img> 的 filter 属性）。
