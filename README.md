# contribcheck

`contribcheck` is a tool to check if the `CONTRIBUTING.md` file is present across all repos in a Github organization,
to make sure that your GitHub org is [friendly to new contributors](https://blog.github.com/2012-09-17-contributing-guidelines/).
It will list the repos that do not have this file. An example output can be seen [here](example-output.md).

Beware: This is very quick and hacky! :)

## Installation

**Prerequisites**: Go version 1.7 or greater.

1. Get the code

```
$ go get github.com/nikhita/contribcheck
```

2. Build

```
$ cd $GOPATH/src/github.com/nikhita/contribcheck
$ go install
```

## Usage

Since Github enforces a rate limit on requests, you will need a personal API token. You can find more details about generating an API token [here](https://github.com/blog/1509-personal-api-tokens).

The org and token are mandatory.

```
contribcheck : v0.1.0
USAGE:
contribcheck -token=<your-token> <org>
  -token string
    	Mandatory GitHub API token
  -v	print version and exit (shorthand)
  -version
    	print version and exit
```

An example on how to use it and a sample output can be see [here](example-output.md).

## License

github-contrib is licensed under the [MIT License](/LICENSE).
