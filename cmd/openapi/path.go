package main

import (
	"log"
	"os"
	"path"
	"strings"
)

// GetExportPathFromArgs returns the export path from the command line arguments
func GetExportPathFromArgs() string {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	userExportPath := "."
	if len(os.Args) < 2 {
		log.Println("Default path is: " + rootPath + " (Please provide a --path to save in another location)")
	} else {
		pathArg := os.Args[1]
		splittedPath := strings.Split(pathArg, "=")
		if len(splittedPath) < 2 || splittedPath[1] == "" || splittedPath[0] != "--path" {
			log.Fatal("Please provide a correct --path=<path> to save in another location")
		}
		userExportPath = splittedPath[1]
	}
	exportPath := path.Join(rootPath, userExportPath)
	return exportPath
}

// GetInternalDirPathFromArgs returns the internal dir path from the command line arguments
func GetInternalDirPathFromArgs() string {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var userInternalDirPath string
	if len(os.Args) < 3 {
		log.Fatal("Please provide a --internal-dir-path=<path>")
	} else {
		internalDirArg := os.Args[2]
		splittedInternalDir := strings.Split(internalDirArg, "=")
		if len(splittedInternalDir) < 2 || splittedInternalDir[1] == "" || splittedInternalDir[0] != "--internal-dir-path" {
			log.Fatal("Please provide a correct --internal-dir-path=<path>")
		}
		userInternalDirPath = splittedInternalDir[1]
	}
	internalFolderPath := path.Join(rootPath, userInternalDirPath)
	return internalFolderPath
}
