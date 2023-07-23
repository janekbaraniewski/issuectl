# issuectl

![Release](https://badgen.net/github/release/janekbaraniewski/issuectl)
[![License](https://img.shields.io/github/license/janekbaraniewski/issuectl.svg)](LICENSE)
![GH workflow](https://github.com/janekbaraniewski/issuectl/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/janekbaraniewski/issuectl)](https://goreportcard.com/report/github.com/janekbaraniewski/issuectl)
[![Go Reference](https://pkg.go.dev/badge/github.com/janekbaraniewski/issuectl.svg)](https://pkg.go.dev/github.com/janekbaraniewski/issuectl)

> # WARNING
> this is still work in progress.

<img src="logo.png" alt="logo" width="200" />

- gives you clean isolated space for each of your issues
- helps you manage multiple git repositories and branches involved

## Install

On macOS:

```bash
➜ brew install janekbaraniewski/janekbaraniewski/issuectl
```

## Quick start

`issuectl init` will create config file with minimal setup to get you going

```bash
➜ issuectl init
? Enter Git user name: John Doe
? Enter Git user email: john@doe.com
? Enter SSH key path: /Users/johndoe/.ssh/id_rsa
? Do you want to configure a backend? Yes
? Select backend type: github
? Enter GitHub Host (Skip for https://api.github.com/):
? Enter GitHub Token: ****************
? Enter GitHub Username: johndoe
? Enter working directory for profile: /Users/johndoe/Workdir
? Enter repository name: my-repo
? Enter repository owner: johndoe
? Enter repository URL: git@github.com:johndoe/my-repo.git
```

To see generated config run

```bash
➜ issuectl config get
currentProfile: default
repositories:
  my-repo:
    name: my-repo
    owner: johndoe
    url: git@github.com:johndoe/my-repo.git
profiles:
  default:
    name: default
    workDir: /Users/johndoe/Workdir
    issueBackend: default
    repoBackend: default
    gituser: John Doe
    repositories:
    - my-repo
    defaultRepository: my-repo
backends:
  default:
    name: default
    backendType: github
    github:
      token: supersecrettoken
      username: johndoe
gitUsers:
  John Doe:
    name: John Doe
    email: john@doe.com
    sshkey: /Users/johndoe/.ssh/id_rsa
```

## Usage

```bash
➜ issuectl --help

issuectl helps managing separate environments for work on multiple issues.

Start work on issue:
	issuectl start [issue_number]

Open PR and link it to issue
	issuectl openpr [issue_number]

Finish work, close the issue
	issuectl finish [issue_number]

Usage:
  issuectl [flags]
  issuectl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manage config
  finish      Cleanup resources and close issue
  help        Help about any command
  init        Initialize configuration
  list        List all issues
  openpr      Opens a pull request for the specified issue
  start       Start work on issue
  workon      Open specified issue in the preferred code editor

Flags:
  -h, --help      help for issuectl
  -v, --version   version for issuectl

Use "issuectl [command] --help" for more information about a command.
```

> This is a basic idea of workflow and what systems this can interact with at each step.
![diagram](diagram.png)

### Repositories

In order to work on code you need to define your repositories:

```bash
➜ issuectl repo add user repoName git@github.com:user/repoName.git
```

### Backends

You'll also need to configure issue backend:

```bash
➜ issuectl backend add github github
```

### Profiles

Once you've set all of this up, you can create your default profile:

```bash
➜ issuectl profile add default /Path/To/Workspace repoName
➜ issuectl use default
```

This will create a profile which will clone `repoName` for each issue. You might want to clone multiple repositories, depending on your environment. For this, run

```bash
➜ issuectl profile addRepo repoName
```

This will add `repoName` to your profile and clone it when starting work on new issue.

### Issues

With all of this in place you can work on issues

#### Start

```bash
➜ issuectl start [issueNumber]
```

#### Open PR

```bash
➜ issuectl openpr [issueNumber]
```

#### Finish

```bash
➜ issuectl finish [issueNumber]
```
