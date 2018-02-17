# LGTM
[![Build Status](https://travis-ci.org/brunograsselli/lgtm.svg?branch=master)](https://travis-ci.org/brunograsselli/lgtm)

Command line tool to list Github pull requests waiting for review.
Early stage project used to automate my workflow (and learn golang).
Please use with caution :)

## Installation

Download it from https://github.com/brunograsselli/lgtm/releases .

Alternatively, you can build it from the source with:

```shell
go get github.com/brunograsselli/lgtm
```

## Configuration

Create the file `~/.lgtm.yml` with your GitHub user and the projects you would like to watch.

```yaml
repos:
  - brunograsselli/lgtm
  - brunograsslli/dotvim
username: brunograsselli
```

## Usage
```
Watch pull requests waiting for your review

Usage:
  lgtm [command]

Available Commands:
  config      Show configuration
  help        Help about any command
  list        List pull requests waiting for your review
  login       Login
  logout      Logout
  show        Show pull request on the browser

Flags:
      --config string   config file (default is $HOME/.lgtm.yml)
  -h, --help            help for lgtm

Use "lgtm [command] --help" for more information about a command.
```

## TODO
* Encrypt saved token
* Save user from login
