package main

import (
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

func main() {
	gitlabToken, found := os.LookupEnv("GITLAB_TOKEN")
	if found != true {
		log.Fatal("GITLAB_TOKEN not found")
	}

	git, err := gitlab.NewClient(gitlabToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	state := "opened"
	wip := "no"
	notLabels := gitlab.Labels{"dependencies"}
	orderBy := "updated_at"
	sort := "asc"

	opts := &gitlab.ListProjectMergeRequestsOptions{}
	opts.State = &state
	opts.WIP = &wip
	opts.NotLabels = &notLabels
	opts.Sort = &sort
	opts.OrderBy = &orderBy

	mrs, resp, err := git.MergeRequests.ListProjectMergeRequests(4835773, opts)
	if err != nil {
		log.Fatalf("Failed to list MRs: %v", err)
	}
	log.Printf("Resp: %v", resp.TotalItems)
	for _, mr := range mrs {
		log.Printf("MR: %v", mr.Title)
	}
}
