// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package waveobj

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type UpdatesRtnType = []WaveObjUpdate

type UIContext struct {
	WindowId    string `json:"windowid"`
	ActiveTabId string `json:"activetabid"`
}

const (
	UpdateType_Update = "update"
	UpdateType_Delete = "delete"
)

const (
	OType_Client      = "client"
	OType_Window      = "window"
	OType_Workspace   = "workspace"
	OType_Tab         = "tab"
	OType_LayoutState = "layout"
	OType_Block       = "block"
)

type WaveObjUpdate struct {
	UpdateType string  `json:"updatetype"`
	OType      string  `json:"otype"`
	OID        string  `json:"oid"`
	Obj        WaveObj `json:"obj,omitempty"`
}

func (update WaveObjUpdate) MarshalJSON() ([]byte, error) {
	rtn := make(map[string]any)
	rtn["updatetype"] = update.UpdateType
	rtn["otype"] = update.OType
	rtn["oid"] = update.OID
	if update.Obj != nil {
		var err error
		rtn["obj"], err = ToJsonMap(update.Obj)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(rtn)
}

func MakeUpdate(obj WaveObj) WaveObjUpdate {
	return WaveObjUpdate{
		UpdateType: UpdateType_Update,
		OType:      obj.GetOType(),
		OID:        GetOID(obj),
		Obj:        obj,
	}
}

func MakeUpdates(objs []WaveObj) []WaveObjUpdate {
	rtn := make([]WaveObjUpdate, 0, len(objs))
	for _, obj := range objs {
		rtn = append(rtn, MakeUpdate(obj))
	}
	return rtn
}

func (update *WaveObjUpdate) UnmarshalJSON(data []byte) error {
	var objMap map[string]any
	err := json.Unmarshal(data, &objMap)
	if err != nil {
		return err
	}
	var ok1, ok2, ok3 bool
	if _, found := objMap["updatetype"]; !found {
		return fmt.Errorf("missing updatetype (in WaveObjUpdate)")
	}
	update.UpdateType, ok1 = objMap["updatetype"].(string)
	if !ok1 {
		return fmt.Errorf("in WaveObjUpdate bad updatetype type %T", objMap["updatetype"])
	}
	if _, found := objMap["otype"]; !found {
		return fmt.Errorf("missing otype (in WaveObjUpdate)")
	}
	update.OType, ok2 = objMap["otype"].(string)
	if !ok2 {
		return fmt.Errorf("in WaveObjUpdate bad otype type %T", objMap["otype"])
	}
	if _, found := objMap["oid"]; !found {
		return fmt.Errorf("missing oid (in WaveObjUpdate)")
	}
	update.OID, ok3 = objMap["oid"].(string)
	if !ok3 {
		return fmt.Errorf("in WaveObjUpdate bad oid type %T", objMap["oid"])
	}
	if _, found := objMap["obj"]; found {
		objMap, ok := objMap["obj"].(map[string]any)
		if !ok {
			return fmt.Errorf("in WaveObjUpdate bad obj type %T", objMap["obj"])
		}
		waveObj, err := FromJsonMap(objMap)
		if err != nil {
			return fmt.Errorf("in WaveObjUpdate error decoding obj: %w", err)
		}
		update.Obj = waveObj
	}
	return nil
}

type Client struct {
	OID           string      `json:"oid"`
	Version       int         `json:"version"`
	WindowIds     []string    `json:"windowids"`
	Meta          MetaMapType `json:"meta"`
	TosAgreed     int64       `json:"tosagreed,omitempty"`
	HasOldHistory bool        `json:"hasoldhistory,omitempty"`
}

func (*Client) GetOType() string {
	return OType_Client
}

// stores the ui-context of the window
// workspaceid, active tab, active block within each tab, window size, etc.
type Window struct {
	OID         string      `json:"oid"`
	Version     int         `json:"version"`
	WorkspaceId string      `json:"workspaceid"`
	ActiveTabId string      `json:"activetabid"`
	Pos         Point       `json:"pos"`
	WinSize     WinSize     `json:"winsize"`
	LastFocusTs int64       `json:"lastfocusts"`
	Meta        MetaMapType `json:"meta"`
}

func (*Window) GetOType() string {
	return OType_Window
}

type Workspace struct {
	OID     string      `json:"oid"`
	Version int         `json:"version"`
	Name    string      `json:"name"`
	TabIds  []string    `json:"tabids"`
	Meta    MetaMapType `json:"meta"`
}

func (*Workspace) GetOType() string {
	return OType_Workspace
}

type Tab struct {
	OID         string      `json:"oid"`
	Version     int         `json:"version"`
	Name        string      `json:"name"`
	LayoutState string      `json:"layoutstate"`
	BlockIds    []string    `json:"blockids"`
	Meta        MetaMapType `json:"meta"`
}

func (*Tab) GetOType() string {
	return OType_Tab
}

func (t *Tab) GetBlockORefs() []ORef {
	rtn := make([]ORef, 0, len(t.BlockIds))
	for _, blockId := range t.BlockIds {
		rtn = append(rtn, ORef{OType: OType_Block, OID: blockId})
	}
	return rtn
}

type LayoutActionData struct {
	ActionType string `json:"actiontype"`
	BlockId    string `json:"blockid"`
	NodeSize   *uint  `json:"nodesize,omitempty"`
	IndexArr   *[]int `json:"indexarr,omitempty"`
	Focused    bool   `json:"focused"`
	Magnified  bool   `json:"magnified"`
}

type LeafOrderEntry struct {
	NodeId  string `json:"nodeid"`
	BlockId string `json:"blockid"`
}

type LayoutState struct {
	OID                   string              `json:"oid"`
	Version               int                 `json:"version"`
	RootNode              any                 `json:"rootnode,omitempty"`
	MagnifiedNodeId       string              `json:"magnifiednodeid,omitempty"`
	FocusedNodeId         string              `json:"focusednodeid,omitempty"`
	LeafOrder             *[]LeafOrderEntry   `json:"leaforder,omitempty"`
	PendingBackendActions *[]LayoutActionData `json:"pendingbackendactions,omitempty"`
	Meta                  MetaMapType         `json:"meta,omitempty"`
}

func (*LayoutState) GetOType() string {
	return OType_LayoutState
}

type FileDef struct {
	FileType string         `json:"filetype,omitempty"`
	Path     string         `json:"path,omitempty"`
	Url      string         `json:"url,omitempty"`
	Content  string         `json:"content,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
}

type BlockDef struct {
	Files map[string]*FileDef `json:"files,omitempty"`
	Meta  MetaMapType         `json:"meta,omitempty"`
}

type StickerClickOptsType struct {
	SendInput   string    `json:"sendinput,omitempty"`
	CreateBlock *BlockDef `json:"createblock,omitempty"`
}

type StickerDisplayOptsType struct {
	Icon    string `json:"icon"`
	ImgSrc  string `json:"imgsrc"`
	SvgBlob string `json:"svgblob,omitempty"`
}

type StickerType struct {
	StickerType string                  `json:"stickertype"`
	Style       map[string]any          `json:"style"`
	ClickOpts   *StickerClickOptsType   `json:"clickopts,omitempty"`
	Display     *StickerDisplayOptsType `json:"display"`
}

type RuntimeOpts struct {
	TermSize TermSize `json:"termsize,omitempty"`
	WinSize  WinSize  `json:"winsize,omitempty"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type WinSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Block struct {
	OID         string         `json:"oid"`
	Version     int            `json:"version"`
	BlockDef    *BlockDef      `json:"blockdef"`
	RuntimeOpts *RuntimeOpts   `json:"runtimeopts,omitempty"`
	Stickers    []*StickerType `json:"stickers,omitempty"`
	Meta        MetaMapType    `json:"meta"`
}

func (*Block) GetOType() string {
	return OType_Block
}

func AllWaveObjTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&Client{}),
		reflect.TypeOf(&Window{}),
		reflect.TypeOf(&Workspace{}),
		reflect.TypeOf(&Tab{}),
		reflect.TypeOf(&Block{}),
		reflect.TypeOf(&LayoutState{}),
	}
}

type TermSize struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}
