// main.go
//
// Copyright (c) 2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package main

import (
	"context"
	"fmt"
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var CmdName = "xfwder"

// openLogFile retrieves an appropriate path for a log file and opens it for writing.
func openLogFile() (io.WriteCloser, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	logDir := filepath.Join(home, "Library", "Logs", CmdName)
	if err = os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logFile := filepath.Join(logDir, fmt.Sprint(CmdName, ".log"))
	return os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
}

// newClient creates a new HTTP client that connects to the given socket path.
func newClient(socketPath string) *http.Client {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, "unix", socketPath)
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

func main() {
	// Prepare logger.
	logFile, err := openLogFile()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to open log file for writing:", err)
		return
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "failed to properly close log file:", err)
		}
	}()
	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	// Check if the user launches the application.
	if len(os.Args) >= 2 {
		logger.Info("initializing application", "args", os.Args)
		return
	}

	// Prepare the forwarder and run it.
	runtime.LockOSThread()
	app := appkit.Application_SharedApplication()
	delegate := &appkit.ApplicationDelegate{}
	delegate.SetApplicationOpenURLs(func(app appkit.Application, urls []foundation.URL) {
		defer app.Terminate(app)

		if len(urls) == 0 {
			logger.Error("no urls are given")
			return
		}
		rawURL := urls[0].AbsoluteString()

		logger.Info("received a URL", "url", rawURL)
		u, err := url.Parse(rawURL)
		if err != nil {
			logger.Error("failed to parse url", "error", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprint("http://localhost", u.RequestURI()), nil)
		if err != nil {
			logger.Error("failed to create request", "error", err)
			return
		}

		_, err = newClient(fmt.Sprintf("/tmp/%v.sock", u.Host)).Do(req)
		if err != nil {
			logger.Error("failed to send request", "error", err)
		}
	})

	app.SetDelegate(delegate)
	app.Run()
}
