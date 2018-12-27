package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = "checkfileinrepo : %s\n"
	// USAGE is an example of how the command should be used.
	USAGE = "USAGE:\ncheckfileinrepo --token=<your-token> --org=<org-name> <filename>"
	// VERSION is the binary version.
	VERSION = "v0.1.0"
)

var (
	token, org, filename string
	version              bool
)

func init() {
	flag.StringVar(&token, "token", "", "GitHub API token (mandatory)")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.StringVar(&org, "org", "", "Name of GitHub organization (mandatory)")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		fmt.Println(USAGE)
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}

	if token == "" {
		usageAndExit("--token cannot be empty.", 1)
	}

	if org == "" {
		usageAndExit("--orgs cannot be empty. Please specify at least one GitHub organization.", 1)
	}
}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		usageAndExit("No arguments specified. Please specify the filename.", 1)
	}
	if len(args) > 1 {
		usageAndExit(fmt.Sprintf("Only one argument is allowed. %d arguments found.", len(args)), 1)
	}

	filename := args[0]

	ctx := context.Background()

	// Create an authenticated client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	getReposWithoutFile(ctx, client, org, filename)
}

// getReposWithoutFile checks for the specified file across all public repos in an org.
func getReposWithoutFile(ctx context.Context, client *github.Client, org, filename string) {
	fmt.Printf("Listing repos for the %s org...\n", org)
	opt := &github.RepositoryListByOrgOptions{
		Type: "public",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}
	repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sleepIfRateLimitExceeded(ctx, client)
	var reposWithoutFile []string

	for _, repository := range repos {
		repo := repository.GetName()
		fmt.Printf("Checking %s...\n", repo)
		_, _, _, err := client.Repositories.GetContents(ctx, org, repo, filename, &github.RepositoryContentGetOptions{Ref: "master"})

		if err != nil && !strings.Contains(err.Error(), "404") {
			fmt.Printf("getting %s file failed: %v", filename, err)
			os.Exit(1)
		}

		if err != nil && strings.Contains(err.Error(), "404") {
			reposWithoutFile = append(reposWithoutFile, fmt.Sprintf("%s/%s", org, repo))
		}
	}

	if len(reposWithoutFile) != 0 {
		fmt.Printf("\nThe following repos in the %s org do not have the %s file:\n", org, filename)

		for _, line := range reposWithoutFile {
			fmt.Println(line)
		}
		return
	}

	fmt.Printf("\nYay! All repos in the %s org have the %s file.\n", org, filename)
}

func sleepIfRateLimitExceeded(ctx context.Context, client *github.Client) {
	rateLimit, _, err := client.RateLimits(ctx)
	if err != nil {
		fmt.Printf("Problem in getting rate limit information %v\n", err)
		return
	}

	if rateLimit.Search.Remaining == 1 {
		timeToSleep := rateLimit.Search.Reset.Sub(time.Now()) + time.Second
		time.Sleep(timeToSleep)
	}
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
