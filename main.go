package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

func main() {
	gitlabToken, found := os.LookupEnv("GITLAB_TOKEN")
	if found != true {
		log.Fatal("GITLAB_TOKEN not found")
	}

	var query struct {
		Project struct {
			MergeRequests struct {
				Nodes []struct {
					Iid        string
					Title      string
					ApprovedBy struct {
						Nodes []struct {
							Username string
						}
					}
					Discussions struct {
						Nodes []struct {
							Resolved bool
							Notes    struct {
								Nodes []struct {
									System bool
									Body   string
									Author struct {
										Username string
									}
								}
							}
						}
					}
				}
			} `graphql:"mergeRequests(state: opened, draft: false)"`
		} `graphql:"project(fullPath: \"sourceability/pim\")"`
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitlabToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient("https://gitlab.com/api/graphql", httpClient)
	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		log.Fatalf("Failed to run query: %v", err)
	}

	var unapprovedMrs []string
	var unresolvedMrs = make(map[string]string)

	for _, mr := range query.Project.MergeRequests.Nodes {
		if len(mr.ApprovedBy.Nodes) == 0 {
			unapprovedMrs = append(unapprovedMrs, mr.Title)
		}

		for _, discussion := range mr.Discussions.Nodes {
			if discussion.Resolved == true {
				continue
			}

			for _, note := range discussion.Notes.Nodes {
				if note.System == true {
					continue
				}

				unresolvedMrs[mr.Iid] = mr.Title
			}
		}
	}

	fmt.Printf("Unapproved MRs\n===============\n")
	for _, mr := range unapprovedMrs {
		fmt.Printf("%v\n", mr)
	}
	fmt.Println()
	fmt.Printf("Unresolved MRs\n===============\n")
	for _, title := range unresolvedMrs {
		fmt.Printf("%v\n", title)
	}
}
