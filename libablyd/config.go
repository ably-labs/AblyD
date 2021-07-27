// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libablyd

type Config struct {
	// base initiaization fields
	CommandName    string    // Command to execute.
	CommandArgs    []string  // Additional args to pass to command.
	
	// settings
	MaxForks       int 		// Max forks
	LogLevel       LogLevel

	// created environment
	Env       []string // Additional environment variables to pass to process ("key=value").
}
