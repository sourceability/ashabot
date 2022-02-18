package ashabot

import (
	"context"
)

func (app *appEnv) fetchMRsForReview() (*mergeRequestsForReview, error) {
	query := openNonDraftMergeRequests{}
	err := app.qc.Query(context.Background(), &query, nil)
	if err != nil {
		return nil, err
	}

	unapprovedMrs := make(map[string]mergeRequest)
	unresolvedMrs := make(map[string]mergeRequest)

	for _, mr := range query.Project.MergeRequests.Nodes {
		if len(mr.ApprovedBy.Nodes) < 2 {
			mr := mergeRequest{iid: mr.Iid, title: mr.Title, url: mr.WebUrl}
			unapprovedMrs[mr.iid] = mr
		}

		for _, discussion := range mr.Discussions.Nodes {
			if discussion.Resolved == true {
				continue
			}

			for _, note := range discussion.Notes.Nodes {
				if note.System == true {
					continue
				}
				mr := mergeRequest{iid: mr.Iid, title: mr.Title, url: mr.WebUrl}
				unresolvedMrs[mr.iid] = mr
			}
		}
	}

	return &mergeRequestsForReview{unapprovedMRs: values(unapprovedMrs), unresolvedMRs: values(unresolvedMrs)}, nil
}

func values(m map[string]mergeRequest) []mergeRequest {
	values := make([]mergeRequest, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

type mergeRequestsForReview struct {
	unapprovedMRs []mergeRequest
	unresolvedMRs []mergeRequest
}

type mergeRequest struct {
	iid        string
	title      string
	approvedBy string
	url        string
}

type openNonDraftMergeRequests struct {
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
				WebUrl      string
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
