// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

import { app, ipcMain } from "electron";
import os from "os";
import path from "path";
import { WaveDevVarName, WaveDevViteVarName } from "../frontend/util/isdev";
import * as keyutil from "../frontend/util/keyutil";

const WaveHomeVarName = "WAVETERM_HOME";

const isDev = !app.isPackaged;
const isDevVite = isDev && process.env.ELECTRON_RENDERER_URL;
if (isDev) {
    process.env[WaveDevVarName] = "1";
}
if (isDevVite) {
    process.env[WaveDevViteVarName] = "1";
}

app.setName(isDev ? "Wave (Dev)" : "Wave");
const unamePlatform = process.platform;
const unameArch: string = process.arch;
keyutil.setKeyUtilPlatform(unamePlatform);

ipcMain.on("get-is-dev", (event) => {
    event.returnValue = isDev;
});
ipcMain.on("get-platform", (event, url) => {
    event.returnValue = unamePlatform;
});
ipcMain.on("get-user-name", (event) => {
    const userInfo = os.userInfo();
    event.returnValue = userInfo.username;
});
ipcMain.on("get-host-name", (event) => {
    event.returnValue = os.hostname();
});

// must match golang
function getWaveHomeDir() {
    const override = process.env[WaveHomeVarName];
    if (override) {
        return override;
    }
    return path.join(os.homedir(), isDev ? ".waveterm-dev" : ".waveterm");
}

function getElectronAppBasePath(): string {
    return path.dirname(import.meta.dirname);
}

function getGoAppBasePath(): string {
    return getElectronAppBasePath().replace("app.asar", "app.asar.unpacked");
}

const wavesrvBinName = `wavesrv.${unameArch}`;

function getWaveSrvPath(): string {
    if (process.platform === "win32") {
        const winBinName = `${wavesrvBinName}.exe`;
        const appPath = path.join(getGoAppBasePath(), "bin", winBinName);
        return `${appPath}`;
    }
    return path.join(getGoAppBasePath(), "bin", wavesrvBinName);
}

function getWaveSrvCwd(): string {
    return getWaveHomeDir();
}

export {
    getElectronAppBasePath,
    getGoAppBasePath,
    getWaveHomeDir,
    getWaveSrvCwd,
    getWaveSrvPath,
    isDev,
    isDevVite,
    unameArch,
    unamePlatform,
};
