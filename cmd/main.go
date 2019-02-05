// Copyright [yyyy] [name of copyright owner]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"runtime"

	api "github.com/jmckind/gvent-api/pkg/gvent/api"
	"github.com/jmckind/gvent-api/version"
	log "github.com/sirupsen/logrus"
)

// main is the entrypoint for the application.
func main() {
	configureLogging()
	printVersion()
	api.Run()
}

// configureLogging sets up logging for the application.
func configureLogging() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//logrus.SetOutput(os.Stdout)

	// Allow the log level to be set using an environment variable
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

// printVersion outputs the architecture and application versions.
func printVersion() {
	log.Infof("Go Version: %s", runtime.Version())
	log.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Infof("operator-sdk Version: %v", version.Version)
}
