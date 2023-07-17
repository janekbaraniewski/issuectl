# issuectl

`issuectl`:

- gives you clean isolated space for each of your issues
- helps you manage multiple git repositories and branches involved

## workdir

You select your `workdir` using `ISSUECTL_WORKDIR` env var (TODO: move this to config).

```txt
workdir
├── issue-1
│   ├── project1@branch1
│   │   └── ...
│   └── project2@branch1
│       └── ...
└── issue-2
    ├── project1@branch2
    │   └── ...
    └── project2@branch2
        └── ...
```
