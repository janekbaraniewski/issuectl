# issuectl

> # WARNING
> this is still work in progress.

`issuectl`:

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
? Select backend type: github
? Enter backend access token: **************************
? Enter working directory for profile: /Users/johndoe/Workspace/work_issues
? Enter repository name: myRepo
? Enter repository owner: myOrg
? Enter repository URL: git@github.com:myOrg/myRepo.git
```

To see generated config run

```bash
➜ issuectl config get
currentprofile: ""
repositories:
  myRepo:
    name: myRepo
    owner: myOrg
    repourl: git@github.com:myOrg/myRepo.git
issues: {}
profiles:
  default:
    name: default
    workdir: /Users/johndoe/Workspace/work_issues
    repository: ""
    backend: ""
    gitusername: ""
    repositories: []
backends:
  default:
    name: default
    type: github
    token: Zm9tZXRva2VubWFnaWNoZWhlaGVoYWhhaGE
gitusers:
  John Doe:
    name: John Doe
    email: john@doe.com
    sshkey: /Users/johndoe/.ssh/id_rsa
```

## Usage

```bash
➜ issuectl --help
issuectl
        issuectl
                issuectl

Usage:
  issuectl [flags]
  issuectl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manage config
  finish      Cleanup resources and close issue
  help        Help about any command
  openpr      Opens a pull request for the specified issue
  start       Start work on issue

Flags:
  -h, --help   help for issuectl

Use "issuectl [command] --help" for more information about a command.
```

Basic idea of workflow and what systems this can interact with at each step
![](diagram.png)

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
➜ issuectl start 15
```

#### Open PR

```bash
➜ issuectl openpr 15
```

#### Finish

```bash
➜ issuectl finish 15
```
