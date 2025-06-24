const { Arch } = require("electron-builder");
const pkg = require("./package.json");
const fs = require("fs");
const path = require("path");

const windowsShouldSign = !!process.env.SM_CODE_SIGNING_CERT_SHA1_HASH;

const config = {
    appId: pkg.build.appId,
    productName: pkg.productName,
    executableName: pkg.productName,
    artifactName: "${productName}-${platform}-${arch}-${version}.${ext}",
    generateUpdatesFilesForAllChannels: true,
    npmRebuild: false,
    nodeGypRebuild: false,
    electronCompile: false,
    files: [
        {
            from: "./dist",
            to: "./dist",
            filter: ["**/*", "!bin/*", "bin/wavesrv.${arch}*", "bin/wsh*"],
        },
        {
            from: ".",
            to: ".",
            filter: ["package.json"],
        },
        "!node_modules",
    ],
    directories: {
        output: "make",
    },
<<<<<<< HEAD
    extraResources: ["tools/run_aipyapp.sh", "tools/run_aipyapp.ps1"],
=======
    extraResources: ["tools/*.sh", "tools/*.ps1", "tools/*.yaml"],
>>>>>>> b38fa878093519204e2733ad558b24d47765f5c6
    asarUnpack: [
        "dist/bin/**/*",
        "dist/docsite/**/*",
        "tools/python-embed/**/*",
        "tools/python-embed-linux/**/*",
        "tools/python-embed-mac-x86_64/**/*",
        "tools/python-embed-mac-arm64/**/*",
<<<<<<< HEAD
        "tools/run_aipyapp.sh",
        "tools/run_aipyapp.ps1",
=======
        "tools/*.sh",
        "tools/*.ps1",
        "tools/*.yaml",
>>>>>>> b38fa878093519204e2733ad558b24d47765f5c6
    ],
    mac: {
        target: [
            {
                target: "zip",
                arch: ["arm64", "x64"],
            },
            {
                target: "dmg",
                arch: ["arm64", "x64"],
            },
        ],
        category: "public.app-category.developer-tools",
        minimumSystemVersion: "10.15.0",
        mergeASARs: true,
        singleArchFiles: "dist/bin/wavesrv.*",
        entitlements: "build/entitlements.mac.plist",
        entitlementsInherit: "build/entitlements.mac.plist",
        extendInfo: {
            NSContactsUsageDescription: "A CLI application running in Wave wants to use your contacts.",
            NSRemindersUsageDescription: "A CLI application running in Wave wants to use your reminders.",
            NSLocationWhenInUseUsageDescription:
                "A CLI application running in Wave wants to use your location information while active.",
            NSLocationAlwaysUsageDescription:
                "A CLI application running in Wave wants to use your location information, even in the background.",
            NSCameraUsageDescription: "A CLI application running in Wave wants to use the camera.",
            NSMicrophoneUsageDescription: "A CLI application running in Wave wants to use your microphone.",
            NSCalendarsUsageDescription: "A CLI application running in Wave wants to use Calendar data.",
            NSLocationUsageDescription: "A CLI application running in Wave wants to use your location information.",
            NSAppleEventsUsageDescription: "A CLI application running in Wave wants to use AppleScript.",
        },
<<<<<<< HEAD
        extraResources: ["tools/python-embed-mac-x86_64/**/*", "tools/python-embed-mac-arm64/**/*"],
=======
        extraResources: ["tools/prompt-optimizer/mac/**/*"],
>>>>>>> b38fa878093519204e2733ad558b24d47765f5c6
    },
    linux: {
        artifactName: "${name}-${platform}-${arch}-${version}.${ext}",
        category: "TerminalEmulator",
        executableName: pkg.name,
<<<<<<< HEAD
        target: ["zip", "deb", "rpm", "snap", "AppImage", "pacman"],
=======
        // target: ["zip", "deb", "rpm", "snap", "AppImage", "pacman"],
        target: ["AppImage"],
>>>>>>> b38fa878093519204e2733ad558b24d47765f5c6
        synopsis: pkg.description,
        description: null,
        desktop: {
            entry: {
                Name: pkg.productName,
                Comment: pkg.description,
                Keywords: "developer;terminal;emulator;",
                Categories: "Development;Utility;",
            },
        },
        executableArgs: ["--enable-features", "UseOzonePlatform", "--ozone-platform-hint", "auto"],
<<<<<<< HEAD
        extraResources: [
            "tools/python-embed-linux/**/*", // Linux 嵌入式 Python
        ],
=======
        extraResources: ["tools/python-embed-linux/**/*", "tools/prompt-optimizer/linux/**/*"],
>>>>>>> b38fa878093519204e2733ad558b24d47765f5c6
    },
    deb: {
        afterInstall: "build/deb-postinstall.tpl",
    },
    win: {
        target: ["nsis", "msi", "zip"],
        signtoolOptions: windowsShouldSign && {
            signingHashAlgorithms: ["sha256"],
            publisherName: "Command Line Inc",
            certificateSubjectName: "Command Line Inc",
            certificateSha1: process.env.SM_CODE_SIGNING_CERT_SHA1_HASH,
        },
        extraResources: [
            "tools/python-embed/**/*", // Windows 嵌入式 Python
<<<<<<< HEAD
=======
            "tools/prompt-optimizer/win/**/*",
>>>>>>> b38fa878093519204e2733ad558b24d47765f5c6
        ],
    },
    appImage: {
        license: "LICENSE",
    },
    snap: {
        base: "core22",
        confinement: "classic",
        allowNativeWayland: true,
        artifactName: "${name}_${version}_${arch}.${ext}",
    },
    publish: {
        provider: "generic",
        url: "https://dl.waveterm.dev/releases-w2",
    },
    afterPack: (context) => {
        if (context.electronPlatformName === "darwin" && context.arch === Arch.universal) {
            const packageBinDir = path.resolve(
                context.appOutDir,
                `${pkg.productName}.app/Contents/Resources/app.asar.unpacked/dist/bin`
            );
            fs.readdirSync(packageBinDir, {
                recursive: true,
                withFileTypes: true,
            })
                .filter((f) => f.isFile() && f.name.startsWith("wavesrv"))
                .forEach((f) => fs.chmodSync(path.resolve(f.parentPath ?? f.path, f.name), 0o755));
        }
<<<<<<< HEAD
=======

        if (context.electronPlatformName === "darwin") {
            const resourcesDir = path.join(context.appOutDir, `${pkg.productName}.app`, "Contents", "Resources");
            let srcDir;
            if (context.arch === "x64") {
                srcDir = "tools/python-embed-mac-x86_64";
            } else if (context.arch === "arm64") {
                srcDir = "tools/python-embed-mac-arm64";
            }
            if (srcDir) {
                const destDir = path.join(resourcesDir, path.basename(srcDir));
                fs.cpSync(srcDir, destDir, { recursive: true });
            }
        }
>>>>>>> b38fa878093519204e2733ad558b24d47765f5c6
    },
};

module.exports = config;
