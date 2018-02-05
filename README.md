# LGTM

Command line tool to list Github pull requests waiting for review.
It is an early stage project used to automated my workflow (and learn golang).
Please use with caution :)

## Installation

Download it from https://github.com/brunograsselli/lgtm/releases .

Alternatively, you can build it from the source with:

```shell
git clone git@github.com:brunograsselli/lgtm.git
cd lgtm
make install
make build
./bin/lgtm
```

## Configuration

Configure your GitHub user name and projects you would like to watch by adding the file `~/.lgtm.yml`.

Eg:
```yaml
repos:
  - brunograsselli/lgtm
  - brunograsslli/dotvim
user: brunograsselli
```

## Usage

```shell
lgtm
```
```
Watch pull requests waiting for your review.

Usage:
  lgtm [command]

Available Commands:
  help        Help about any command
  list        List pull requests waiting for your review
  login       Login to GitHub

Flags:
      --config string   config file (default is $HOME/.lgtm.yml)
  -h, --help            help for lgtm

Use "lgtm [command] --help" for more information about a command.
```

## TODO
* Encrypt saved token
* Implement command to logout
* Save user from login
* Implement command to open a PR
