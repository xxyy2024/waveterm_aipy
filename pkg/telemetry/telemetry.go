// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"context"
	"database/sql/driver"
	"log"
	"time"

	"github.com/wavetermdev/waveterm/pkg/panichandler"
	"github.com/wavetermdev/waveterm/pkg/util/daystr"
	"github.com/wavetermdev/waveterm/pkg/util/dbutil"
	"github.com/wavetermdev/waveterm/pkg/wavebase"
	"github.com/wavetermdev/waveterm/pkg/wconfig"
	"github.com/wavetermdev/waveterm/pkg/wstore"
)

const MaxTzNameLen = 50

type ActivityDisplayType struct {
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	DPR      float64 `json:"dpr"`
	Internal bool    `json:"internal,omitempty"`
}

type ActivityUpdate struct {
	FgMinutes     int                   `json:"fgminutes,omitempty"`
	ActiveMinutes int                   `json:"activeminutes,omitempty"`
	OpenMinutes   int                   `json:"openminutes,omitempty"`
	NumTabs       int                   `json:"numtabs,omitempty"`
	NewTab        int                   `json:"newtab,omitempty"`
	NumBlocks     int                   `json:"numblocks,omitempty"`
	NumWindows    int                   `json:"numwindows,omitempty"`
	NumSSHConn    int                   `json:"numsshconn,omitempty"`
	NumWSLConn    int                   `json:"numwslconn,omitempty"`
	NumMagnify    int                   `json:"nummagnify,omitempty"`
	NumPanics     int                   `json:"numpanics,omitempty"`
	Startup       int                   `json:"startup,omitempty"`
	Shutdown      int                   `json:"shutdown,omitempty"`
	SetTabTheme   int                   `json:"settabtheme,omitempty"`
	BuildTime     string                `json:"buildtime,omitempty"`
	Displays      []ActivityDisplayType `json:"displays,omitempty"`
	Renderers     map[string]int        `json:"renderers,omitempty"`
	Blocks        map[string]int        `json:"blocks,omitempty"`
	WshCmds       map[string]int        `json:"wshcmds,omitempty"`
	Conn          map[string]int        `json:"conn,omitempty"`
}

type ActivityType struct {
	Day           string        `json:"day"`
	Uploaded      bool          `json:"-"`
	TData         TelemetryData `json:"tdata"`
	TzName        string        `json:"tzname"`
	TzOffset      int           `json:"tzoffset"`
	ClientVersion string        `json:"clientversion"`
	ClientArch    string        `json:"clientarch"`
	BuildTime     string        `json:"buildtime"`
	OSRelease     string        `json:"osrelease"`
}

type TelemetryData struct {
	ActiveMinutes int                   `json:"activeminutes"`
	FgMinutes     int                   `json:"fgminutes"`
	OpenMinutes   int                   `json:"openminutes"`
	NumTabs       int                   `json:"numtabs"`
	NumBlocks     int                   `json:"numblocks,omitempty"`
	NumWindows    int                   `json:"numwindows,omitempty"`
	NumSSHConn    int                   `json:"numsshconn,omitempty"`
	NumWSLConn    int                   `json:"numwslconn,omitempty"`
	NumMagnify    int                   `json:"nummagnify,omitempty"`
	NewTab        int                   `json:"newtab"`
	NumStartup    int                   `json:"numstartup,omitempty"`
	NumShutdown   int                   `json:"numshutdown,omitempty"`
	NumPanics     int                   `json:"numpanics,omitempty"`
	SetTabTheme   int                   `json:"settabtheme,omitempty"`
	Displays      []ActivityDisplayType `json:"displays,omitempty"`
	Renderers     map[string]int        `json:"renderers,omitempty"`
	Blocks        map[string]int        `json:"blocks,omitempty"`
	WshCmds       map[string]int        `json:"wshcmds,omitempty"`
	Conn          map[string]int        `json:"conn,omitempty"`
}

func (tdata TelemetryData) Value() (driver.Value, error) {
	return dbutil.QuickValueJson(tdata)
}

func (tdata *TelemetryData) Scan(val interface{}) error {
	return dbutil.QuickScanJson(tdata, val)
}

func IsTelemetryEnabled() bool {
	settings := wconfig.GetWatcher().GetFullConfig()
	return settings.Settings.TelemetryEnabled
}

func IsAutoUpdateEnabled() bool {
	settings := wconfig.GetWatcher().GetFullConfig()
	return settings.Settings.AutoUpdateEnabled
}

func AutoUpdateChannel() string {
	settings := wconfig.GetWatcher().GetFullConfig()
	return settings.Settings.AutoUpdateChannel
}

// Wraps UpdateCurrentActivity, spawns goroutine, and logs errors
func GoUpdateActivityWrap(update ActivityUpdate, debugStr string) {
	go func() {
		defer panichandler.PanicHandlerNoTelemetry("GoUpdateActivityWrap")
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		err := UpdateActivity(ctx, update)
		if err != nil {
			// ignore error, just log, since this is not critical
			log.Printf("error updating current activity (%s): %v\n", debugStr, err)
		}
	}()
}

func UpdateActivity(ctx context.Context, update ActivityUpdate) error {
	now := time.Now()
	dayStr := daystr.GetCurDayStr()
	txErr := wstore.WithTx(ctx, func(tx *wstore.TxWrap) error {
		var tdata TelemetryData
		query := `SELECT tdata FROM db_activity WHERE day = ?`
		found := tx.Get(&tdata, query, dayStr)
		if !found {
			query = `INSERT INTO db_activity (day, uploaded, tdata, tzname, tzoffset, clientversion, clientarch, buildtime, osrelease)
                                      VALUES (  ?,        0,     ?,      ?,        ?,             ?,          ?,         ?,         ?)`
			tzName, tzOffset := now.Zone()
			if len(tzName) > MaxTzNameLen {
				tzName = tzName[0:MaxTzNameLen]
			}
			tx.Exec(query, dayStr, tdata, tzName, tzOffset, wavebase.WaveVersion, wavebase.ClientArch(), wavebase.BuildTime, wavebase.UnameKernelRelease())
		}
		tdata.FgMinutes += update.FgMinutes
		tdata.ActiveMinutes += update.ActiveMinutes
		tdata.OpenMinutes += update.OpenMinutes
		tdata.NewTab += update.NewTab
		tdata.NumStartup += update.Startup
		tdata.NumShutdown += update.Shutdown
		tdata.SetTabTheme += update.SetTabTheme
		tdata.NumMagnify += update.NumMagnify
		tdata.NumPanics += update.NumPanics
		if update.NumTabs > 0 {
			tdata.NumTabs = update.NumTabs
		}
		if update.NumBlocks > 0 {
			tdata.NumBlocks = update.NumBlocks
		}
		if update.NumWindows > 0 {
			tdata.NumWindows = update.NumWindows
		}
		if update.NumSSHConn > 0 && update.NumSSHConn > tdata.NumSSHConn {
			tdata.NumSSHConn = update.NumSSHConn
		}
		if update.NumWSLConn > 0 && update.NumWSLConn > tdata.NumWSLConn {
			tdata.NumWSLConn = update.NumWSLConn
		}
		if len(update.Renderers) > 0 {
			if tdata.Renderers == nil {
				tdata.Renderers = make(map[string]int)
			}
			for key, val := range update.Renderers {
				tdata.Renderers[key] += val
			}
		}
		if len(update.WshCmds) > 0 {
			if tdata.WshCmds == nil {
				tdata.WshCmds = make(map[string]int)
			}
			for key, val := range update.WshCmds {
				tdata.WshCmds[key] += val
			}
		}
		if len(update.Conn) > 0 {
			if tdata.Conn == nil {
				tdata.Conn = make(map[string]int)
			}
			for key, val := range update.Conn {
				tdata.Conn[key] += val
			}
		}
		if len(update.Displays) > 0 {
			tdata.Displays = update.Displays
		}
		if len(update.Blocks) > 0 {
			tdata.Blocks = update.Blocks
		}
		query = `UPDATE db_activity
                 SET tdata = ?,
                     clientversion = ?,
                     buildtime = ?
                 WHERE day = ?`
		tx.Exec(query, tdata, wavebase.WaveVersion, wavebase.BuildTime, dayStr)
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return nil
}

func GetNonUploadedActivity(ctx context.Context) ([]*ActivityType, error) {
	var rtn []*ActivityType
	txErr := wstore.WithTx(ctx, func(tx *wstore.TxWrap) error {
		query := `SELECT * FROM db_activity WHERE uploaded = 0 ORDER BY day DESC LIMIT 30`
		tx.Select(&rtn, query)
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return rtn, nil
}

func MarkActivityAsUploaded(ctx context.Context, activityArr []*ActivityType) error {
	dayStr := daystr.GetCurDayStr()
	txErr := wstore.WithTx(ctx, func(tx *wstore.TxWrap) error {
		query := `UPDATE db_activity SET uploaded = 1 WHERE day = ?`
		for _, activity := range activityArr {
			if activity.Day == dayStr {
				continue
			}
			tx.Exec(query, activity.Day)
		}
		return nil
	})
	return txErr
}
