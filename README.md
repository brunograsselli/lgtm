# LGTM
[![Build Status](https://travis-ci.org/brunograsselli/lgtm.svg?branch=master)](https://travis-ci.org/brunograsselli/lgtm)

Command line tool to list Github pull requests waiting for review.
Early stage project used to automate my workflow (and learn golang).
Please use with caution :)

## Installation

Download it from https://github.com/brunograsselli/lgtm/releases .

Alternatively, you can build it from the source with:

```shell
$ go get github.com/brunograsselli/lgtm
```

## Configuration

Login to GitHub:

```shell
$ lgtm login
```

Edit the file `~/.lgtm.yml` and add the repositories you would like to watch:

```yaml
username: brunograsselli
repos:
  - brunograsselli/lgtm
  - brunograsslli/dotvim
```

## Usage
```shell
$ lgtm
```

## Help
```
Watch pull requests waiting for your review

Usage:
  lgtm [flags]
  lgtm [command]

Available Commands:
  config      Show configuration
  help        Help about any command
  list        List pull requests waiting for your review
  login       Login
  logout      Logout

Flags:
  -a, --all             List all open pull requests
      --config string   config file (default is $HOME/.lgtm.yml)
  -h, --help            help for lgtm

Use "lgtm [command] --help" for more information about a command.
```

## TODO
* Encrypt saved token
