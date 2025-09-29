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
	DEFAULTPORT = "8080"
	DEFAULTHOST = "localhost"
	DEFAULTREEORDER = 4
)


// permissions
const (
	OWNERPERMISSION = 0755 //// 755 is a common permission setting that allows the owner to read, write, and execute the directory,
	// while group members and others can read and execute it.
	BUCKETPERMISSION = 0700 //// 700 is a common permission setting that allows the owner to read, write, and execute the directory,
	// while group members and others have no permissions.
)