// +build ignore

// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func genLDFlags(version string) string {
	ldflagsStr := "-s -w"
	ldflagsStr += " -X github.com/minio/minio/cmd.Version=" + version
	ldflagsStr += " -X github.com/minio/minio/cmd.ReleaseTag=" + releaseTag(version)
	ldflagsStr += " -X github.com/minio/minio/cmd.CommitID=" + commitID()
	ldflagsStr += " -X github.com/minio/minio/cmd.ShortCommitID=" + commitID()[:12]
	ldflagsStr += " -X github.com/minio/minio/cmd.GOPATH=" + os.Getenv("GOPATH")
	ldflagsStr += " -X github.com/minio/minio/cmd.GOROOT=" + os.Getenv("GOROOT")
	return ldflagsStr
}

// genReleaseTag prints release tag to the console for easy git tagging.
func releaseTag(version string) string {
	relPrefix := "DEVELOPMENT"
	if prefix := os.Getenv("MINIO_RELEASE"); prefix != "" {
		relPrefix = prefix
	}

	relSuffix := ""
	if hotfix := os.Getenv("MINIO_HOTFIX"); hotfix != "" {
		relSuffix = hotfix
	}

	relTag := strings.Replace(version, " ", "-", -1)
	relTag = strings.Replace(relTag, ":", "-", -1)
	relTag = strings.Replace(relTag, ",", "", -1)
	relTag = relPrefix + "." + relTag

	if relSuffix != "" {
		relTag += "." + relSuffix
	}

	return relTag
}

// commitID returns the abbreviated commit-id hash of the last commit.
func commitID() string {
	// git log --format="%H" -n1
	var (
		commit []byte
		err    error
	)
	cmdName := "git"
	cmdArgs := []string{"log", "--format=%H", "-n1"}
	if commit, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Error generating git commit-id: ", err)
		os.Exit(1)
	}

	return strings.TrimSpace(string(commit))
}

func main() {
	fmt.Println(genLDFlags(time.Now().UTC().Format(time.RFC3339)))
}
