// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wavetermdev/waveterm/pkg/waveobj"
	"github.com/wavetermdev/waveterm/pkg/wshrpc"
	"github.com/wavetermdev/waveterm/pkg/wshrpc/wshclient"
	"github.com/wavetermdev/waveterm/pkg/wshutil"
)

var (
	rootCmd = &cobra.Command{
		Use:          "wsh",
		Short:        "CLI tool to control Wave Terminal",
		Long:         `wsh is a small utility that lets you do cool things with Wave Terminal, right from the command line`,
		SilenceUsage: true,
	}
)

var WrappedStdin io.Reader = os.Stdin
var RpcClient *wshutil.WshRpc
var RpcContext wshrpc.RpcContext
var UsingTermWshMode bool
var blockArg string
var WshExitCode int

func WriteStderr(fmtStr string, args ...interface{}) {
	output := fmt.Sprintf(fmtStr, args...)
	if UsingTermWshMode {
		output = strings.ReplaceAll(output, "\n", "\r\n")
	}
	fmt.Fprint(os.Stderr, output)
}

func WriteStdout(fmtStr string, args ...interface{}) {
	output := fmt.Sprintf(fmtStr, args...)
	if UsingTermWshMode {
		output = strings.ReplaceAll(output, "\n", "\r\n")
	}
	fmt.Print(output)
}

func preRunSetupRpcClient(cmd *cobra.Command, args []string) error {
	err := setupRpcClient(nil)
	if err != nil {
		return err
	}
	return nil
}

func getIsTty() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

type RunEFnType = func(*cobra.Command, []string) error

func activityWrap(activityStr string, origRunE RunEFnType) RunEFnType {
	return func(cmd *cobra.Command, args []string) (rtnErr error) {
		defer func() {
			sendActivity(activityStr, rtnErr == nil)
		}()
		return origRunE(cmd, args)
	}
}

func resolveBlockArg() (*waveobj.ORef, error) {
	oref := blockArg
	if oref == "" {
		oref = "this"
	}
	fullORef, err := resolveSimpleId(oref)
	if err != nil {
		return nil, fmt.Errorf("resolving blockid: %w", err)
	}
	return fullORef, nil
}

// returns the wrapped stdin and a new rpc client (that wraps the stdin input and stdout output)
func setupRpcClient(serverImpl wshutil.ServerImpl) error {
	jwtToken := os.Getenv(wshutil.WaveJwtTokenVarName)
	if jwtToken == "" {
		wshutil.SetTermRawModeAndInstallShutdownHandlers(true)
		UsingTermWshMode = true
		RpcClient, WrappedStdin = wshutil.SetupTerminalRpcClient(serverImpl)
		return nil
	}
	rpcCtx, err := wshutil.ExtractUnverifiedRpcContext(jwtToken)
	if err != nil {
		return fmt.Errorf("error extracting rpc context from %s: %v", wshutil.WaveJwtTokenVarName, err)
	}
	RpcContext = *rpcCtx
	sockName, err := wshutil.ExtractUnverifiedSocketName(jwtToken)
	if err != nil {
		return fmt.Errorf("error extracting socket name from %s: %v", wshutil.WaveJwtTokenVarName, err)
	}
	RpcClient, err = wshutil.SetupDomainSocketRpcClient(sockName, serverImpl)
	if err != nil {
		return fmt.Errorf("error setting up domain socket rpc client: %v", err)
	}
	wshclient.AuthenticateCommand(RpcClient, jwtToken, &wshrpc.RpcOpts{NoResponse: true})
	// note we don't modify WrappedStdin here (just use os.Stdin)
	return nil
}

func isFullORef(orefStr string) bool {
	_, err := waveobj.ParseORef(orefStr)
	return err == nil
}

func resolveSimpleId(id string) (*waveobj.ORef, error) {
	if isFullORef(id) {
		orefObj, err := waveobj.ParseORef(id)
		if err != nil {
			return nil, fmt.Errorf("error parsing full ORef: %v", err)
		}
		return &orefObj, nil
	}
	rtnData, err := wshclient.ResolveIdsCommand(RpcClient, wshrpc.CommandResolveIdsData{Ids: []string{id}}, &wshrpc.RpcOpts{Timeout: 2000})
	if err != nil {
		return nil, fmt.Errorf("error resolving ids: %v", err)
	}
	oref, ok := rtnData.ResolvedIds[id]
	if !ok {
		return nil, fmt.Errorf("id not found: %q", id)
	}
	return &oref, nil
}

// this will send wsh activity to the client running on *your* local machine (it does not contact any wave cloud infrastructure)
// if you've turned off telemetry in your local client, this data never gets sent to us
// no parameters or timestamps are sent, as you can see below, it just sends the name of the command (and if there was an error)
// (e.g. "wsh ai ..." would send "ai")
// this helps us understand which commands are actually being used so we know where to concentrate our effort
func sendActivity(wshCmdName string, success bool) {
	if RpcClient == nil || wshCmdName == "" {
		return
	}
	dataMap := make(map[string]int)
	dataMap[wshCmdName] = 1
	if !success {
		dataMap[wshCmdName+"#"+"error"] = 1
	}
	wshclient.WshActivityCommand(RpcClient, dataMap, nil)
}

// Execute executes the root command.
func Execute() {
	defer func() {
		r := recover()
		if r != nil {
			WriteStderr("[panic] %v\n", r)
			debug.PrintStack()
			wshutil.DoShutdown("", 1, true)
		} else {
			wshutil.DoShutdown("", WshExitCode, false)
		}
	}()
	rootCmd.PersistentFlags().StringVarP(&blockArg, "block", "b", "", "for commands which require a block id")
	err := rootCmd.Execute()
	if err != nil {
		wshutil.DoShutdown("", 1, true)
		return
	}
}
