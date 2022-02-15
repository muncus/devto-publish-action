package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/muncus/devto-publish-action/devto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var stateFile = flag.String("state", "", "File containing serialized version of post metadata.")
var postFiles = flag.String("post_files", "", "Comma-separated list of files to upload.")
var postDir = flag.String("post_dir", "", "A directory containing posts to upload.")
var apiKey = flag.String("apikey", "", "Dev.to api key")
var debugFlag = flag.Bool("debug", false, "Dump http request and response, for debugging")

func main() {
	flag.Parse()

	// Make pretty logs.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *postDir != "" && *postFiles != "" {
		log.Fatal().Msg("Cannot specify both --post_files and --post_dir\n")
		return
	}
	if *apiKey == "" {
		log.Fatal().Msg("--apikey must be specified.")
		return
	}

	devtoSyncer, err := devto.NewSyncer(*stateFile, *apiKey)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	devtoSyncer.SetDebug(*debugFlag)

	if len(*postFiles) > 0 {
		filesToSync := strings.Split(*postFiles, ",")
		for _, file := range filesToSync {
			_, err := devtoSyncer.SyncFile(file)
			if err != nil {
				log.Error().Str("file", file).Msg(fmt.Sprintf("Sync Failed: %s", err))
			} else {
				log.Info().Str("file", file).Msg("Success.")
			}
		}
	}
	if *postDir != "" {
		// articles and errors are logged to the Logger, so no need to check them here.
		_, _ = devtoSyncer.Sync(*postDir, log.Logger)
	}

	err = devtoSyncer.DumpState(*stateFile)
	if err != nil {
		fmt.Println("Failed to dump state:", err)
	}
}
