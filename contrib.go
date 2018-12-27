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
	BANNER = "contribcheck : %s\n"
	// USAGE is an example of how the command should be used.
	USAGE = "USAGE:\ncontribcheck -token=<your-token> <org>"
	// VERSION is the binary version.
	VERSION = "v0.1.0"
)

var (
	token   string
	version bool
)

func init() {
	flag.StringVar(&token, "token", "", "Mandatory GitHub API token")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")

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
		usageAndExit("GitHub token cannot be empty", 1)
	}
}

func main() {
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Wrong number of arguments!")
		os.Exit(1)
	}

	org := args[0]
	ctx := context.Background()

	// Create an authenticated client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	getReposWithoutContribFile(ctx, client, org)
}

// getAllRepos checks for the CONTRIBUTING.md file across all public repos in an org.
func getReposWithoutContribFile(ctx context.Context, client *github.Client, org string) {
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
	var reposWithoutContribFile []string

	for _, repository := range repos {
		repo := repository.GetName()
		fmt.Printf("Checking %s...\n", repo)
		_, _, _, err := client.Repositories.GetContents(ctx, org, repo, "code-of-conduct.md", &github.RepositoryContentGetOptions{Ref: "master"})

		if err != nil && !strings.Contains(err.Error(), "404") {
			fmt.Printf("getting code-of-conduct.md file failed: %v", err)
			os.Exit(1)
		}

		if err != nil && strings.Contains(err.Error(), "404") {
			reposWithoutContribFile = append(reposWithoutContribFile, fmt.Sprintf("%s/%s", org, repo))
		}
	}

	if len(reposWithoutContribFile) != 0 {
		fmt.Printf("\nThe following repos in the %s org do not have the code-of-conduct.md file:\n", org)

		for _, line := range reposWithoutContribFile {
			fmt.Println(line)
		}
		return
	}

	fmt.Printf("\nYay! All repos in the %s org have the code-of-conduct.md file.\n", org)
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
