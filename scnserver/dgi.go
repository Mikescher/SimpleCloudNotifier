package server

import (
	_ "embed"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"strings"
)

//go:embed DOCKER_GIT_INFO
var FileDockerGitInfo string

var CommitHash *string
var VCSType *string
var CommitTime *string
var BranchName *string
var RemoteURL *string

func init() {
	for _, v := range strings.Split(FileDockerGitInfo, "\n") {
		if v == "" {
			continue
		} else if strings.HasPrefix(v, "VCSTYPE=") {
			VCSType = langext.Ptr(v[len("VCSTYPE="):])
			fmt.Printf("Found DGI Config: '%s' := '%s'\n", "VCSType", *VCSType)
		} else if strings.HasPrefix(v, "BRANCH=") {
			BranchName = langext.Ptr(v[len("BRANCH="):])
			fmt.Printf("Found DGI Config: '%s' := '%s'\n", "BranchName", *BranchName)
		} else if strings.HasPrefix(v, "HASH=") {
			CommitHash = langext.Ptr(v[len("HASH="):])
			fmt.Printf("Found DGI Config: '%s' := '%s'\n", "CommitHash", *CommitHash)
		} else if strings.HasPrefix(v, "COMMITTIME=") {
			CommitTime = langext.Ptr(v[len("COMMITTIME="):])
			fmt.Printf("Found DGI Config: '%s' := '%s'\n", "CommitTime", *CommitTime)
		} else if strings.HasPrefix(v, "REMOTE=") {
			RemoteURL = langext.Ptr(v[len("REMOTE="):])
			fmt.Printf("Found DGI Config: '%s' := '%s'\n", "RemoteURL", *RemoteURL)
		} else {
			fmt.Printf("[ERROR] Failed to parse DGI Config '%s'\n", v)
		}
	}
}
