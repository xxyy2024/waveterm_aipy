// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package waveobj

import (
	"strings"
)

const Entity_Any = "any"

// for typescript typing
type MetaTSType struct {
	// shared
	View           string   `json:"view,omitempty"`
	Controller     string   `json:"controller,omitempty"`
	Title          string   `json:"title,omitempty"`
	File           string   `json:"file,omitempty"`
	Url            string   `json:"url,omitempty"`
	Connection     string   `json:"connection,omitempty"`
	Edit           bool     `json:"edit,omitempty"`
	History        []string `json:"history,omitempty"`
	HistoryForward []string `json:"history:forward,omitempty"`

	DisplayName  string  `json:"display:name,omitempty"`
	DisplayOrder float64 `json:"display:order,omitempty"`

	Icon      string `json:"icon,omitempty"`
	IconColor string `json:"icon:color,omitempty"`

	Frame                    bool   `json:"frame,omitempty"`
	FrameClear               bool   `json:"frame:*,omitempty"`
	FrameBorderColor         string `json:"frame:bordercolor,omitempty"`
	FrameBorderColor_Focused string `json:"frame:bordercolor:focused,omitempty"`

	Cmd               string            `json:"cmd,omitempty"`
	CmdClear          bool              `json:"cmd:*,omitempty"`
	CmdInteractive    bool              `json:"cmd:interactive,omitempty"`
	CmdLogin          bool              `json:"cmd:login,omitempty"`
	CmdRunOnStart     bool              `json:"cmd:runonstart,omitempty"`
	CmdClearOnStart   bool              `json:"cmd:clearonstart,omitempty"`
	CmdClearOnRestart bool              `json:"cmd:clearonrestart,omitempty"`
	CmdEnv            map[string]string `json:"cmd:env,omitempty"`
	CmdCwd            string            `json:"cmd:cwd,omitempty"`
	CmdNoWsh          bool              `json:"cmd:nowsh,omitempty"`

	GraphClear     bool     `json:"graph:*,omitempty"`
	GraphNumPoints int      `json:"graph:numpoints,omitempty"`
	GraphMetrics   []string `json:"graph:metrics,omitempty"`

	// for tabs
	Bg          string  `json:"bg,omitempty"`
	BgClear     bool    `json:"bg:*,omitempty"`
	BgOpacity   float64 `json:"bg:opacity,omitempty"`
	BgBlendMode string  `json:"bg:blendmode,omitempty"`

	TermClear          bool   `json:"term:*,omitempty"`
	TermFontSize       int    `json:"term:fontsize,omitempty"`
	TermFontFamily     string `json:"term:fontfamily,omitempty"`
	TermMode           string `json:"term:mode,omitempty"`
	TermTheme          string `json:"term:theme,omitempty"`
	TermLocalShellPath string `json:"term:localshellpath,omitempty"` // matches settings

	Count int `json:"count,omitempty"` // temp for cpu plot. will remove later
}

type MetaDataDecl struct {
	Key        string   `json:"key"`
	Desc       string   `json:"desc,omitempty"`
	Type       string   `json:"type"` // string, int, float, bool, array, object
	Default    any      `json:"default,omitempty"`
	StrOptions []string `json:"stroptions,omitempty"`
	NumRange   []*int   `json:"numrange,omitempty"` // inclusive, null means no limit
	Entity     []string `json:"entity"`             // what entities this applies to, e.g. "block", "tab", "any", etc.
	Special    []string `json:"special,omitempty"`  // special handling.  things that need to happen if this gets updated
}

type MetaPresetDecl struct {
	Preset string   `json:"preset"`
	Desc   string   `json:"desc,omitempty"`
	Keys   []string `json:"keys"`
	Entity []string `json:"entity"` // what entities this applies to, e.g. "block", "tab", etc.
}

// returns a clean copy of meta with mergeMeta merged in
// if mergeSpecial is false, then special keys will not be merged (like display:*)
func MergeMeta(meta MetaMapType, metaUpdate MetaMapType, mergeSpecial bool) MetaMapType {
	rtn := make(MetaMapType)
	for k, v := range meta {
		rtn[k] = v
	}
	// deal with "section:*" keys
	for k := range metaUpdate {
		if !strings.HasSuffix(k, ":*") {
			continue
		}
		if !metaUpdate.GetBool(k, false) {
			continue
		}
		prefix := strings.TrimSuffix(k, ":*")
		if prefix == "" {
			continue
		}
		// delete "[prefix]" and all keys that start with "[prefix]:"
		prefixColon := prefix + ":"
		for k2 := range rtn {
			if k2 == prefix || strings.HasPrefix(k2, prefixColon) {
				delete(rtn, k2)
			}
		}
	}
	// now deal with regular keys
	for k, v := range metaUpdate {
		if !mergeSpecial && strings.HasPrefix(k, "display:") {
			continue
		}
		if strings.HasSuffix(k, ":*") {
			continue
		}
		if v == nil {
			delete(rtn, k)
			continue
		}
		rtn[k] = v
	}
	return rtn
}
