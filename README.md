# issuectl

> # WARNING
> this is still work in progress.

`issuectl`:

- gives you clean isolated space for each of your issues
- helps you manage multiple git repositories and branches involved

## Install

On macOS:

```bash
> brew install janekbaraniewski/janekbaraniewski/issuectl
```

## Usage

Basic idea of workflow and what systems this can interact with at each step
![](diagram.png)

### Repositories

In order to work on code you need to define your repositories:

```bash
> issuectl repo add user repoName git@github.com:user/repoName.git
```

### Backends

You'll also need to configure issue backend:

```bash
> issuectl backend add github github
```

### Profiles

Once you've set all of this up, you can create your default profile:

```bash
> issuectl profile add default /Path/To/Workspace repoName
> issuectl use default
```

This will create a profile which will clone `repoName` for each issue. You might want to clone multiple repositories, depending on your environment. For this, run

```bash
> issuectl profile addRepo repoName
```

This will add `repoName` to your profile and clone it when starting work on new issue.

### Issues

With all of this in place you can work on issues

#### Start

```bash
> issuectl start 15
```

#### Open PR

```bash
> issuectl openpr 15
```

#### Finish

```bash
> issuectl finish 15
```
