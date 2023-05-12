# gh-environments

A GitHub `gh` [CLI](https://cli.github.com/) extension to list environments and their associated metadata for an organization and/or specific repositories. 

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation) instructions.

2. Install the extension:
   ```sh
   gh extension install katiem0/gh-organization-webhooks
   ```

For more information: [`gh extension install`](https://cli.github.com/manual/gh_extension_install).

## Usage

The `gh-environments` extension supports `GitHub.com` and GitHub Enterprise Server, through the use of `--hostname` and `--source-hostname`, and the following commands:

```sh
$ gh environments -h

List repo environments and metadata, including listing and creating environment secrets and variables.

Usage:
  environments [command]

Available Commands:
  list        Generate a report of environments and metadata.
  secrets     List and Create Environment secrets.
  variables   List and Create Environment variables.

Flags:
      --help   Show help for command

Use "environments [command] --help" for more information about a command.
```

### List Environments

Environment metadata can be listed and written to a `csv` file for an organization or specific repository.


```sh
$ gh environments list -h

Generate a report of environments and metadata for a single repository or all repositories in an organization.

Usage:
  environments list [flags] <organization> [repo ...] 

Flags:
  -d, --debug                To debug logging
      --hostname string      GitHub Enterprise Server hostname (default "github.com")
  -o, --output-file string   Name of file to write CSV report (default "report-20230512095310.csv")
  -t, --token string         GitHub Personal Access Token (default "gh auth token")

Global Flags:
      --help   Show help for command
```

The output `csv` file contains the following information:

| Field Name | Description |
|:-----------|:------------|
|`RepositoryName` | The name of the repository where the data is extracted from. |
|`RepositoryID`| The `ID` associated with the Repository, for API usage. |
|`EnvironmentName`| The name of the repository specific environment. |
|`AdminBypass`| `True`/`False` flag to indicate if administrators are allowed to bypass configured protection rules. |
|`WaitTimer| The an amount of time to wait before allowing deployments to proceed. |
|`Reviewers`| Specified people or teams that have the ability to approve workflow runs when tey access the environment. |
|`BranchPolicyType`| Indicates if the environment can only be deployed to specific branches. (Values: `protected`, `custom`, or `null`, where `null` indicates **any branch from the repo can deploy**.)|
|`Branches`| If `BranchPolicyType = custom`, list of specific branch name patterns the environment deployment is limited to|
|`SecretsTotalCount`| The number of Actions secrets that are associated with the environment. |
|`VariablesTotalCount`| The number of Actions variables that are associated with the environment. |

### Environment Secrets

### Environment Variables