# LGTM

CLI tool to list Github pull requests waiting for review.

# Configuration Example (~/.lgtm.yml)

```shell
repos:
  - brunograsselli/lgtm
  - brunograsslli/dotvim
user: github_user
token: github_oauth_token
```

# Usage

```shell
$ lgtm list
```

# TODO
* Generate OAUTH token automatically
* Implement command to open a PR
