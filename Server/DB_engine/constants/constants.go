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
	BUCKETDIR      = "buckets"
	METABUCKETFILE = "buckets_meta.json"
)

// filepaths
var (
	AUTHFILEPATH   = filepath.Join(DEFAULTDATADIR, AUTHFILENAME)
	CONFIGFILEPATH = filepath.Join(CURRENTDIR, CONFIGFILENAME)
	DEFAULTWALPATH = filepath.Join(DEFAULTDATADIR, WALFILENAME)
	DBBUCKETSPATH = filepath.Join(DEFAULTDATADIR, BUCKETDIR)
	GLOBALMETAPATH = filepath.Join(DEFAULTDATADIR, METABUCKETFILE)
)

// configs
const (
	DEFAULTREEORDER = 4
)

const (
	SERVER_DEFAULT_PORT="9090"
)
// permissions
const (
	OWNERPERMISSION = 0755 
	REMOVEWRITEPERMISSION = 0555 
)