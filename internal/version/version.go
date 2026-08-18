package version

import (
	"fmt"
	"runtime/debug"
	"time"
)

var IS_MAIN_BUILD bool = false // MUST BE SET TO FALSE FOR EVERY VERSION BUMP
var SITE string = "https://pokeapi.co/api/v2/"
var VERSION string = "0.1.8"
var DEV_TAG string = "dev"
var PROJECT_NAME string = "Poké7DB"

// {PROJECT_NAME} v{VERSION}-{DEV_TAG}-{timestamp}_{commitID}

func GenerateVersionNumber() {
	// If not development build, don't print commit details
	// %TODO pre-release / beta/ alpha / release tag support
	if IS_MAIN_BUILD {
		VERSION = fmt.Sprintf("%s v%s", PROJECT_NAME, VERSION)
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		VERSION = fmt.Sprintf("%s UNKNOWN VERSION", PROJECT_NAME)
		return
	}

	commit := "unknown"
	modified := ""
	revTime := ""

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.time":
			t, err := time.Parse(time.RFC3339, setting.Value)
			if err != nil {
				fmt.Println(err)
				return
			}
			revTime = fmt.Sprintf("%x", t.Unix())
		case "vcs.revision":
			commit = setting.Value[:7]
		case "vcs.modified":
			modified = " (modified)"
		}
	}

	VERSION = fmt.Sprintf("%s v%s-%s [%s.%s]%s", PROJECT_NAME, VERSION, DEV_TAG, revTime, commit, modified)
}
