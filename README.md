# checkfileinrepo

`checkfileinrepo` is a tool to check if a specified file exists across all repos in a GitHub organization.
It will list all repos that do not have the file. Examples can be seen [here](examples.md).

A usecase of this tool can be to check if all repos in your GitHub org have the `CONTRIBUTING.md` file,
to make sure that your org is [friendly to new contributors](https://blog.github.com/2012-09-17-contributing-guidelines/).

Beware: This is very quick and hacky! :)

## Installation

1. Get the code

```
$ go get github.com/nikhita/checkfileinrepo
```

2. Build

```
$ cd $GOPATH/src/github.com/nikhita/checkfileinrepo
$ go install
```

## Usage

Since Github enforces a rate limit on requests, you will need a personal API token. You can find more details about generating an API token [here](https://github.com/blog/1509-personal-api-tokens).

The token, org and filename are mandatory.

```
checkfileinrepo : v0.1.0
USAGE:
checkfileinrepo --token=<your-token> --org=<org-name> <filename>
  -org string
    	Name of GitHub organization (mandatory)
  -token string
    	GitHub API token (mandatory)
  -v	print version and exit (shorthand)
  -version
    	print version and exit
```

Examples on how to use it along with sample outputs can be seen [here](examples.md).

## License

checkfileinrepo is licensed under the [MIT License](/LICENSE).
