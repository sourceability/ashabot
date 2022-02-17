package main

import (
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

func main() {
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	git, err := gitlab.NewClient(gitlabToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	state := "opened"
	wip := "no"
	notLabels := gitlab.Labels{"dependencies"}

	opts := &gitlab.ListProjectMergeRequestsOptions{}
	opts.State = &state
	opts.WIP = &wip
	opts.NotLabels = &notLabels

	mrs, resp, err := git.MergeRequests.ListProjectMergeRequests(4835773, opts)
	if err != nil {
		log.Fatalf("Failed to list MRs: %v", err)
	}
	log.Printf("Resp: %v", resp.TotalItems)
	for _, mr := range mrs {
		log.Printf("MR: %v", mr.Title)
	}
}
