package version

import (
	"fmt"
	"runtime/debug"
	"time"
)

var G_IsMainBuild bool = false // MUST BE SET TO FALSE FOR EVERY VERSION BUMP
var G_Site string = "https://pokeapi.co/api/v2/"
var G_Version string = "0.1.9"
var G_DevTag string = "dev"
var G_ProjectName string = "Poké7DB"

// {PROJECT_NAME} v{VERSION}-{DEV_TAG}-{timestamp}_{commitID}

func GenerateVersionNumber() {
	// If not development build, don't print commit details
	// %TODO pre-release / beta/ alpha / release tag support
	if G_IsMainBuild {
		G_Version = fmt.Sprintf("%s v%s", G_ProjectName, G_Version)
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		G_Version = fmt.Sprintf("%s UNKNOWN VERSION", G_ProjectName)
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

	G_Version = fmt.Sprintf("%s v%s-%s [%s.%s]%s", G_ProjectName, G_Version, G_DevTag, revTime, commit, modified)
}
