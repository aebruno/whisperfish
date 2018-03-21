// Copyright 2016 Andrew E. Bruno
//
// This file is part of Whisperfish.
//
// Whisperfish is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Whisperfish is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Whisperfish.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/aebruno/whisperfish/ui"
)

var (
	Version        = "dev-build"
	versionFlag    bool
	debugFlag      bool
	convertFlag    bool
	attachmentFlag bool
)

func init() {
	flag.BoolVar(&versionFlag, "version", false, "show version")
	flag.BoolVar(&versionFlag, "v", false, "show version (shorthand)")
	flag.BoolVar(&debugFlag, "debug", false, "debug to file")
	flag.BoolVar(&debugFlag, "d", false, "debug to file (shorthand)")
	flag.BoolVar(&convertFlag, "convert", false, "convert datastore")
	flag.BoolVar(&attachmentFlag, "fix-attachments", false, "Add extensions to attachments")
}

func main() {
	flag.Parse()
	if versionFlag {
		fmt.Printf("Whisperfish v%s\n", Version)
		os.Exit(0)
	}

	if convertFlag {
		ui.ConvertDataStore()
		os.Exit(0)
	} else if attachmentFlag {
		ui.AddAttachmentExtensions()
		os.Exit(0)
	}

	if debugFlag {
		logFile, err := os.OpenFile("/tmp/whisperfish-app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
			os.Exit(-1)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	ui.Run(Version)
}
