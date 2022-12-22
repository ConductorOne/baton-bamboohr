![Baton Logo](./docs/images/baton-logo.png)

# `baton-bamboohr` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-bamboohr.svg)](https://pkg.go.dev/github.com/conductorone/baton-bamboohr) ![main ci](https://github.com/conductorone/baton-bamboohr/actions/workflows/main.yaml/badge.svg)

`baton-bamboohr` is a connector for BambooHR built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the BambooHR API to sync data about users.

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-bamboohr
baton-bamboohr
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_COMPANY_DOMAIN=customerID -e BATON_API_KEY=apiKey ghcr.io/conductorone/baton-bamboohr:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-bamboohr/cmd/baton-bamboohr@main

BATON_COMPANY_DOMAIN=customerID BATON_API_KEY=apiKey
baton resources
```

# Data Model

`baton-bamboohr` will pull down information about the following BambooHR resources:
- Users
  - Users supervisors

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-bamboohr` Command Line Usage

```
baton-bamboohr

Usage:
  baton-bamboohr [flags]
  baton-bamboohr [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
  -f, --file string                         The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
      --company-domain string               The company domain used for the BambooHR account. ($BATON_COMPANY_DOMAIN)
      --api-key string                      The api-key for the BambooHR account. ($BATON_API_KEY)
  -h, --help                                help for baton-bamboohr
      --log-format string                   The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                    The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -v, --version                             version for baton-bamboohr

Use "baton-bamboohr [command] --help" for more information about a command.

```
