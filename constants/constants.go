package constants

import (
	"os"
	"path/filepath"
)

// directories
var (
	HOMEDIR, _     = os.UserHomeDir()
	DEFAULTDATADIR = filepath.Join(HOMEDIR, ".bytedata") + string(os.PathSeparator)
	CURRENTDIR     = "./"
)

// filenames
const (
	AUTHFILENAME   = "auth.json"
	CONFIGFILENAME = "test_config.json"
	WALFILENAME    = "wal.log"
)

// filepaths
var (
	AUTHFILEPATH   = filepath.Join(DEFAULTDATADIR, AUTHFILENAME)
	CONFIGFILEPATH = filepath.Join(CURRENTDIR, CONFIGFILENAME)
	DEFAULTWALPATH = filepath.Join(DEFAULTDATADIR, WALFILENAME)
)

// configs
const (
	DEFAULTPORT = "8080"
	DEFAULTHOST = "localhost"
)
