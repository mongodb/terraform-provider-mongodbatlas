[![GoDoc](https://godoc.org/github.com/mongodb-forks/digest?status.svg)](https://godoc.org/github.com/mongodb-forks/digest)

![.github/workflows/pr.yml](https://github.com/mongodb-forks/digest/workflows/.github/workflows/pr.yml/badge.svg?branch=master&event=push)

# Golang HTTP Digest Authentication

## Overview

This is a fork of the (unmaintained) code.google.com/p/mlab-ns2/gae/ns/digest package.
There's a descriptor leak in the original package, so this fork was created to patch
the leak.

### Update 2020

This is a fork of the now unmaintained fork of [digest](https://github.com/bobziuchkovski/digest).
This implementation now supports the SHA-256 algorithm which was added as part of [rfc 7616](https://tools.ietf.org/html/rfc7616).

## Usage

See the [godocs](https://godoc.org/github.com/bobziuchkovski/digest) for details.

## Contributing

**Contributions are welcome!**

The code is linted with [golangci-lint](https://golangci-lint.run/).  This library also defines *git hooks* that format and lint the code.

Before submitting a PR, please run `make setup link-git-hooks` to set up your local development environment.

## Fork Maintainer

Bob Ziuchkovski (@bobziuchkovski)

## Original Authors

Bipasa Chattopadhyay <bipasa@cs.unc.edu>
Eric Gavaletz <gavaletz@gmail.com>
Seon-Wook Park <seon.wook@swook.net>

## License

Apache 2.0
