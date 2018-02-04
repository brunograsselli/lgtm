# LGTM

Command line tool to list Github pull requests waiting for review.

## Configuration Example (~/.lgtm.yml)

```yaml
repos:
  - brunograsselli/lgtm
  - brunograsslli/dotvim
user: github_user
```

## Usage

```shell
$ lgtm login
$ lgtm list
```

## TODO
* Encrypt saved token
* Implement command to logout
* Save user from login
* Implement command to open a PR
