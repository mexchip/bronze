package main

import (
	// "fmt"
	"net/url"
	"os"
	// "regexp"
	"strings"

	. "github.com/reujab/bronze/types"
	"gopkg.in/src-d/go-git.v4"
)

// // reformat scp-like url (e. g. ssh)
// // https://github.com/motemen/ghq/blob/master/url.go
// func fixUrl(url string) string {
// 	hasSchemePattern := regexp.MustCompile("^[^:]+://")
// 	scpLikeUrlPattern := regexp.MustCompile("^([^@]+@)?([^:]+):/?(.+)$")

// 	if !hasSchemePattern.MatchString(url) && scpLikeUrlPattern.MatchString(url) {
// 		matched := scpLikeUrlPattern.FindStringSubmatch(url)
// 		user := matched[1]
// 		host := matched[2]
// 		path := matched[3]
// 		return fmt.Sprintf("ssh://%s%s/%s", user, host, path)
// 	}

// 	return url
// }

// the git segment provides useful information about a git repository such as the domain of the "origin" remote (with an icon), the current branch, and whether the HEAD is dirty
func gitnewSegment(segment *Segment) {
	dir, err := os.Getwd()
	die(err)
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return
	}
	// defer repo.Free()

	var domainName string
	list, err := repo.Remotes()
	if err == nil {
		for _, r := range list {
			if strings.Contains(r.Config().Name, "origin") {
				remoteUrl := fixUrl(r.Config().URLs[0])
				uri, err := url.Parse(remoteUrl)
				if err == nil && len(uri.Hostname()) > 4 {
					// strip the tld off the hostname
					domainName = uri.Hostname()[:len(uri.Hostname())-4]
				}
				break
			}
		}

		// fmt.Println("domainName:", domainName)
	}

	// stash is not supported by go-git
	// (https://github.com/src-d/go-git/blob/master/COMPATIBILITY.md)
	// var stashes int
	// repo.Stashes.Foreach(func(int, string, *git.Oid) error {
	// 	stashes++
	// 	return nil
	// })

	// var ahead, behind int
	var branch string
	// head, err := repo.Head()
	// if err == nil {
	// 	upstream, err := head.Branch().Upstream()
	// 	if err == nil {
	// 		ahead, behind, err = repo.AheadBehind(head.Branch().Target(), upstream.Target())
	// 		die(err)
	// 	}

	head, err := repo.Head()
	if err == nil {
		ref := head.Name()
		if ref.IsBranch() {
			refSlice := strings.Split(ref.String(), "/")
			branch = refSlice[len(refSlice) - 1]
		} else if "HEAD" == ref.String() {
			branch = "HEAD"
		}
	}

	var dirty, modified, staged bool
	workTree, err := repo.Worktree()
	if err == nil {
		status, err := workTree.Status()

		dirty = !status.IsClean()

		if err == nil {
			for path := range status {
				// fmt.Println(path)
				fileStatus := status.File(path)
				// fmt.Println("staging: ", string(fileStatus.Staging), "worktree: ", string(fileStatus.Worktree))
				switch fileStatus.Staging {
				case 'M':
					staged = true
				case 'A':
					staged = true
				case 'D':
					staged = true
				case 'R':
					staged = true
				case 'C':
					staged = true
				case 'U':
					staged = true
				}
				switch fileStatus.Worktree {
				case '?':
					modified = true
				case 'M':
					modified = true
				case 'A':
					modified = true
				case 'D':
					modified = true
				case 'R':
					modified = true
				case 'C':
					modified = true
				case 'U':
					modified = true
				}
			}
		}
	}

	var segments []string
	domainIcon := icons[domainName]
	if domainIcon == "" {
		domainIcon = icons["git"]
	}
	if domainIcon != "" {
		segments = append(segments, domainIcon)
	}

	// if stashes != 0 || ahead != 0 || behind != 0 {
	// 	section := strings.Repeat(icons["stash"], stashes) + strings.Repeat(icons["ahead"], ahead) + strings.Repeat(icons["behind"], behind)
	// 	if section != "" {
	// 		segments = append(segments, section)
	// 	}
	// }

	if branch != "" {
		segments = append(segments, branch)
	}
	if dirty {
		segment.Background = "yellow"

		var section string
		if modified {
			section += icons["modified"]
		}
		if staged {
			section += icons["staged"]
		}
		if section != "" {
			segments = append(segments, section)
		}
	}
	segment.Value = strings.Join(segments, " ")
}
