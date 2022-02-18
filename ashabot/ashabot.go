package ashabot

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

func CLI(args []string) int {
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}
	mrs := app.fetchMRsForReview()
	app.output.write(&mrs)

	return 0
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("ashabot", flag.ContinueOnError)
	fl.BoolVar(&app.debug, "debug", false, "Set to true to enable debug output")
	outputDesination := fl.String("output", "cli", "Output destination: cli or slack")

	if err := fl.Parse(args); err != nil {
		return err
	}

	if *outputDesination != "cli" && *outputDesination != "slack" {
		fmt.Fprintf(os.Stderr, "Invalid output destination: %s\n", *outputDesination)
		fl.Usage()
		return flag.ErrHelp
	}
	if *outputDesination == "slack" {
		app.output = &slackOutput{}
	} else {
		app.output = &cliOutput{}
	}

	gitlabToken, found := os.LookupEnv("GITLAB_TOKEN")
	if found != true {
		log.Fatal("GITLAB_TOKEN not found")
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitlabToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	app.qc = graphql.NewClient("https://gitlab.com/api/graphql", httpClient)

	return nil
}

type outputClient interface {
	write(mrs *mergeRequestsForReview)
}

type appEnv struct {
	qc     *graphql.Client
	output outputClient
	debug  bool
}

type cliOutput struct{}
type slackOutput struct{}

func (out *cliOutput) write(mrs *mergeRequestsForReview) {
	fmt.Printf("MRs with fewer than 2 approvals (%v)\n========================================\n", len(mrs.unapprovedMRs))

	for i, mr := range mrs.unapprovedMRs {
		fmt.Printf("%v. %v (%v)\n", i+1, mr.title, mr.url)
	}
	fmt.Println()
	fmt.Printf("MRs with unresolved comments (%v)\n==================================\n", len(mrs.unresolvedMRs))
	for i, mr := range mrs.unresolvedMRs {
		fmt.Printf("%v. %v (%v)\n", i+1, mr.title, mr.url)
	}
}

func (out *slackOutput) write(mrs *mergeRequestsForReview) {
	fmt.Printf("Send to slack: %v", mrs)
}
