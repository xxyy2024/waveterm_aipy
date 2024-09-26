// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

// wave core application coordinator
package wcore

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wavetermdev/waveterm/pkg/blockcontroller"
	"github.com/wavetermdev/waveterm/pkg/telemetry"
	"github.com/wavetermdev/waveterm/pkg/waveobj"
	"github.com/wavetermdev/waveterm/pkg/wps"
	"github.com/wavetermdev/waveterm/pkg/wstore"
)

// the wcore package coordinates actions across the storage layer
// orchestrating the wave object store, the wave pubsub system, and the wave rpc system

// TODO bring Tx infra into wcore

const DefaultTimeout = 2 * time.Second
const DefaultActivateBlockTimeout = 60 * time.Second

func DeleteBlock(ctx context.Context, tabId string, blockId string) error {
	err := wstore.DeleteBlock(ctx, tabId, blockId)
	if err != nil {
		return fmt.Errorf("error deleting block: %w", err)
	}
	go blockcontroller.StopBlockController(blockId)
	sendBlockCloseEvent(tabId, blockId)
	return nil
}

func sendBlockCloseEvent(tabId string, blockId string) {
	waveEvent := wps.WaveEvent{
		Event: wps.Event_BlockClose,
		Scopes: []string{
			waveobj.MakeORef(waveobj.OType_Tab, tabId).String(),
			waveobj.MakeORef(waveobj.OType_Block, blockId).String(),
		},
		Data: blockId,
	}
	wps.Broker.Publish(waveEvent)
}

func DeleteTab(ctx context.Context, workspaceId string, tabId string) error {
	tabData, err := wstore.DBGet[*waveobj.Tab](ctx, tabId)
	if err != nil {
		return fmt.Errorf("error getting tab: %w", err)
	}
	if tabData == nil {
		return nil
	}
	// close blocks (sends events + stops block controllers)
	for _, blockId := range tabData.BlockIds {
		err := DeleteBlock(ctx, tabId, blockId)
		if err != nil {
			return fmt.Errorf("error deleting block %s: %w", blockId, err)
		}
	}
	// now delete tab (also deletes layout)
	err = wstore.DeleteTab(ctx, workspaceId, tabId)
	if err != nil {
		return fmt.Errorf("error deleting tab: %w", err)
	}

	return nil
}

// returns tabid
func CreateTab(ctx context.Context, windowId string, tabName string, activateTab bool) (string, error) {
	windowData, err := wstore.DBMustGet[*waveobj.Window](ctx, windowId)
	if err != nil {
		return "", fmt.Errorf("error getting window: %w", err)
	}
	tab, err := wstore.CreateTab(ctx, windowData.WorkspaceId, tabName)
	if err != nil {
		return "", fmt.Errorf("error creating tab: %w", err)
	}
	if activateTab {
		err = wstore.SetActiveTab(ctx, windowId, tab.OID)
		if err != nil {
			return "", fmt.Errorf("error setting active tab: %w", err)
		}
	}
	return tab.OID, nil
}

func CreateWindow(ctx context.Context, winSize *waveobj.WinSize) (*waveobj.Window, error) {
	windowId := uuid.NewString()
	workspaceId := uuid.NewString()
	if winSize == nil {
		winSize = &waveobj.WinSize{
			Width:  1200,
			Height: 850,
		}
	}
	window := &waveobj.Window{
		OID:         windowId,
		WorkspaceId: workspaceId,
		Pos: waveobj.Point{
			X: 100,
			Y: 100,
		},
		WinSize: *winSize,
	}
	err := wstore.DBInsert(ctx, window)
	if err != nil {
		return nil, fmt.Errorf("error inserting window: %w", err)
	}
	ws := &waveobj.Workspace{
		OID:  workspaceId,
		Name: "w" + workspaceId[0:8],
	}
	err = wstore.DBInsert(ctx, ws)
	if err != nil {
		return nil, fmt.Errorf("error inserting workspace: %w", err)
	}
	_, err = CreateTab(ctx, windowId, "T1", true)
	if err != nil {
		return nil, fmt.Errorf("error inserting tab: %w", err)
	}
	client, err := wstore.DBGetSingleton[*waveobj.Client](ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	client.WindowIds = append(client.WindowIds, windowId)
	err = wstore.DBUpdate(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("error updating client: %w", err)
	}
	return wstore.DBMustGet[*waveobj.Window](ctx, windowId)
}

func checkAndFixWindow(ctx context.Context, windowId string) {
	window, err := wstore.DBMustGet[*waveobj.Window](ctx, windowId)
	if err != nil {
		log.Printf("error getting window %q (in checkAndFixWindow): %v\n", windowId, err)
		return
	}
	workspace, err := wstore.DBMustGet[*waveobj.Workspace](ctx, window.WorkspaceId)
	if err != nil {
		log.Printf("error getting workspace %q (in checkAndFixWindow): %v\n", window.WorkspaceId, err)
		return
	}
	if len(workspace.TabIds) == 0 {
		log.Printf("fixing workspace with no tabs %q (in checkAndFixWindow)\n", workspace.OID)
		_, err = CreateTab(ctx, windowId, "T1", true)
		if err != nil {
			log.Printf("error creating tab (in checkAndFixWindow): %v\n", err)
		}
	}
}

// returns (new-window, first-time, error)
func EnsureInitialData() (*waveobj.Window, bool, error) {
	// does not need to run in a transaction since it is called on startup
	ctx, cancelFn := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFn()
	firstRun := false
	client, err := wstore.DBGetSingleton[*waveobj.Client](ctx)
	if err == wstore.ErrNotFound {
		client, err = CreateClient(ctx)
		if err != nil {
			return nil, false, fmt.Errorf("error creating client: %w", err)
		}
		firstRun = true
	}
	log.Printf("clientid: %s\n", client.OID)
	if len(client.WindowIds) == 1 {
		checkAndFixWindow(ctx, client.WindowIds[0])
	}
	if len(client.WindowIds) > 0 {
		return nil, false, nil
	}
	window, err := CreateWindow(ctx, &waveobj.WinSize{Height: 0, Width: 0})
	if err != nil {
		return nil, false, fmt.Errorf("error creating window: %w", err)
	}
	return window, firstRun, nil
}

func CreateClient(ctx context.Context) (*waveobj.Client, error) {
	client := &waveobj.Client{
		OID:       uuid.NewString(),
		WindowIds: []string{},
	}
	err := wstore.DBInsert(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("error inserting client: %w", err)
	}
	return client, nil
}

func CreateBlock(ctx context.Context, tabId string, blockDef *waveobj.BlockDef, rtOpts *waveobj.RuntimeOpts) (*waveobj.Block, error) {
	if blockDef == nil {
		return nil, fmt.Errorf("blockDef is nil")
	}
	if blockDef.Meta == nil || blockDef.Meta.GetString(waveobj.MetaKey_View, "") == "" {
		return nil, fmt.Errorf("no view provided for new block")
	}
	blockData, err := wstore.CreateBlock(ctx, tabId, blockDef, rtOpts)
	if err != nil {
		return nil, fmt.Errorf("error creating block: %w", err)
	}
	go func() {
		blockView := blockDef.Meta.GetString(waveobj.MetaKey_View, "")
		if blockView == "" {
			return
		}
		tctx, cancelFn := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelFn()
		telemetry.UpdateActivity(tctx, telemetry.ActivityUpdate{
			Renderers: map[string]int{blockView: 1},
		})
	}()
	return blockData, nil
}
