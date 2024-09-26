// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

import { getApi } from "./global";

class ContextMenuModelType {
    handlers: Map<string, () => void> = new Map(); // id -> handler

    constructor() {
        getApi().onContextMenuClick(this.handleContextMenuClick.bind(this));
    }

    handleContextMenuClick(id: string): void {
        const handler = this.handlers.get(id);
        if (handler) {
            handler();
        }
    }

    _convertAndRegisterMenu(menu: ContextMenuItem[]): ElectronContextMenuItem[] {
        const electronMenuItems: ElectronContextMenuItem[] = [];
        for (const item of menu) {
            const electronItem: ElectronContextMenuItem = {
                role: item.role,
                type: item.type,
                label: item.label,
                id: crypto.randomUUID(),
                checked: item.checked,
            };
            if (item.click) {
                this.handlers.set(electronItem.id, item.click);
            }
            if (item.submenu) {
                electronItem.submenu = this._convertAndRegisterMenu(item.submenu);
            }
            electronMenuItems.push(electronItem);
        }
        return electronMenuItems;
    }

    showContextMenu(menu: ContextMenuItem[], ev: React.MouseEvent<any>): void {
        this.handlers.clear();
        const electronMenuItems = this._convertAndRegisterMenu(menu);
        getApi().showContextMenu(electronMenuItems);
    }
}

const ContextMenuModel = new ContextMenuModelType();

export { ContextMenuModel, ContextMenuModelType };
