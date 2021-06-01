// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libablyd

import (
	"fmt"
	"strings"
)

const (
	gatewayInterface = "websocketd-CGI/0.1"
)

var headerNewlineToSpace = strings.NewReplacer("\n", " ", "\r", " ")
var headerDashToUnderscore = strings.NewReplacer("-", "_")

func createEnv(handler *AblyDHandler, log *LogScope) []string {
	standardEnvCount := 20

	parentLen := len(handler.config.ParentEnv)
	env := make([]string, 0, standardEnvCount+parentLen+len(handler.config.Env))

	// This variable could be rewritten from outside
	env = appendEnv(env, "SERVER_SOFTWARE", handler.config.ServerSoftware)

	parentStarts := len(env)
	env = append(env, handler.config.ParentEnv...)

	// IMPORTANT ---> Adding a header? Make sure standardEnvCount (above) is up to date.

	// Standard CGI specification headers.
	// As defined in http://tools.ietf.org/html/rfc3875
	env = appendEnv(env, "SCRIPT_NAME", handler.URLInfo.ScriptPath)
	env = appendEnv(env, "PATH_INFO", handler.URLInfo.PathInfo)

	// Not supported, but we explicitly clear them so we don't get leaks from parent environment.
	env = appendEnv(env, "AUTH_TYPE", "")
	env = appendEnv(env, "CONTENT_LENGTH", "")
	env = appendEnv(env, "CONTENT_TYPE", "")
	env = appendEnv(env, "REMOTE_IDENT", "")
	env = appendEnv(env, "REMOTE_USER", "")

	// Non standard, but commonly used headers.
	// env = appendEnv(env, "UNIQUE_ID", handler.Id) // Based on Apache mod_unique_id.

	// The following variables are part of the CGI specification, but are optional
	// and not set by websocketd:
	//
	//   AUTH_TYPE, REMOTE_USER, REMOTE_IDENT
	//     -- Authentication left to the underlying programs.
	//
	//   CONTENT_LENGTH, CONTENT_TYPE
	//     -- makes no sense for WebSocket connections.
	if log.MinLevel == LogDebug {
		for i, v := range env {
			if i >= parentStarts && i < parentLen+parentStarts {
				log.Debug("env", "Parent envvar: %v", v)
			} else {
				log.Debug("env", "Std. variable: %v", v)
			}
		}
	}

	for _, v := range handler.config.Env {
		env = append(env, v)
		log.Debug("env", "External variable: %s", v)
	}

	return env
}

// Adapted from net/http/header.go
func appendEnv(env []string, k string, v ...string) []string {
	if len(v) == 0 {
		return env
	}

	vCleaned := make([]string, 0, len(v))
	for _, val := range v {
		vCleaned = append(vCleaned, strings.TrimSpace(headerNewlineToSpace.Replace(val)))
	}
	return append(env, fmt.Sprintf("%s=%s",
		strings.ToUpper(k),
		strings.Join(vCleaned, ", ")))
}
