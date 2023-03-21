# cnay

A command-line tool that resolves a list of hostnames to their corresponding IP addresses.

It filters the results to include only unique IPv4 addresses from:

- Hostnames with an A record.
- A records from CNAMEs on the same domain (`foo.example.com -> bar.example.com -> 1.1.1.1`).

Is useful for quickly eliminating subdomains that are likely pointing to third-party services.

## Installation

```
go install github.com/melvinsh/cnay@latest
```

## Usage

```
cnay [-l file] [-d] [-r] [-pb]
```

### Flags

| Flag | Description |
|------|-------------|
| `-l` | Path to the file containing the list of hostnames |
| `-d` | Enable debug output |
| `-r` | Show original hostname in brackets |
| `-pb` | Enable progress bar |

If no `-l` flag is provided, `cnay` reads from `stdin`.

### Examples

Resolve hostnames from a file:

```
$ cnay -l hostnames.txt
8.8.8.8
8.8.4.4
```

Resolve hostnames from `stdin`:

```
$ echo "www.google.com" | cnay
172.217.6.4
```

Show original hostname in brackets:

```
$ cnay -r -l hostnames.txt
8.8.8.8 [google-public-dns-a.google.com]
8.8.4.4 [google-public-dns-b.google.com]
```

## Dependencies

This tool relies on the `golang.org/x/net/publicsuffix` package to determine if two hostnames are on the same domain. This package is a Go implementation of the Public Suffix List, which is a cross-vendor initiative to provide an accurate list of domain name suffixes.

To install this dependency, run:

```
go get -u golang.org/x/net/publicsuffix
```
