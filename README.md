# gh-environments

A GitHub `gh` [CLI](https://cli.github.com/) extension to list environments and their associated metadata for an organization and/or specific repositories. 

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation) instructions.

2. Install the extension:
   ```sh
   gh extension install katiem0/gh-environments
   ```

For more information: [`gh extension install`](https://cli.github.com/manual/gh_extension_install).

## Usage

The `gh-environments` extension supports `GitHub.com` and GitHub Enterprise Server, through the use of `--hostname` and the following commands:

```sh
$ gh environments -h

List and create repo environments and metadata, including listing and creating environment secrets and variables.

Usage:
  environments [command]

Available Commands:
  create      Create environments and metadata.
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

#### Report Output

The output `csv` file contains the following information:

| Field Name | Description |
|:-----------|:------------|
|`RepositoryName` | The name of the repository where the data is extracted from. |
|`RepositoryID`| The `ID` associated with the Repository, for API usage. |
|`EnvironmentName`| The name of the repository specific environment. |
|`AdminBypass`| `True`/`False` flag to indicate if administrators are allowed to bypass configured protection rules. |
|`WaitTimer`| The an amount of time to wait before allowing deployments to proceed. |
|`Reviewers`| Specified people or teams that have the ability to approve workflow runs when they access the environment. In the format `<UserOrTeam>;Name;ID` and reviewers delimited by `|` |
|`PreventSelfReview` | Indicates if a Reviewer is able to approve/deny the workflow run on a specific environment |
|`BranchPolicyType`| Indicates if the environment can only be deployed to specific branches. (Values: `protected`, `custom`, or `null`, where `null` indicates **any branch from the repo can deploy**.)|
|`Branches`| If `BranchPolicyType = custom`, list of specific branch name patterns the environment deployment is limited to. In the format `Name;<BranchOrTag>` and policies delimited by `|`|
|`CustomDeploymentProtectionPolicy`| Lists the custom deployment protection rules that are enabled for an environment. In the format: `PolicyID;Enabled;AppID;AppSlug` and policies delimited by `|`|
|`SecretsTotalCount`| The number of Actions secrets that are associated with the environment. |
|`VariablesTotalCount`| The number of Actions variables that are associated with the environment. |

### Create Environments

The `gh environments create` command will create environments from a `csv` file using `--from-file` following the format outlined in [`gh environments create`](#environment-create).

```sh
$ gh environments create -h

Create environments and metadata for specified environments per repository in an organization from a file.

Usage:
  environments create  <target organization> [flags]

Flags:
  -d, --debug              To debug logging
  -f, --from-file string   Path and Name of CSV file to create environments from
      --hostname string    GitHub Enterprise Server hostname (default "github.com")
  -t, --token string       GitHub personal access token for organization to write to (default "gh auth token")

Global Flags:
      --help   Show help for command
```

The `create` command utilizes the following fields in their given format but expects all headers listed [Report Output](#report-output): 

| Field Name | Description |
|:-----------|:------------|
|`RepositoryName` | The name of the repository where the data is extracted from. |
|`EnvironmentName`| The name of the repository specific environment. |
|`AdminBypass`| `True`/`False` flag to indicate if administrators are allowed to bypass configured protection rules. |
|`WaitTimer`| The an amount of time to wait before allowing deployments to proceed. |
|`Reviewers`| Specified people or teams that have the ability to approve workflow runs when they access the environment. In the format `<UserOrTeam>;Name;ID` and reviewers delimited by `|` |
|`PreventSelfReview` | Indicates if a Reviewer is able to approve/deny the workflow run on a specific environment |
|`BranchPolicyType`| Indicates if the environment can only be deployed to specific branches. (Values: `protected`, `custom`, or `null`, where `null` indicates **any branch from the repo can deploy**.)|
|`Branches`| If `BranchPolicyType = custom`, list of specific branch name patterns the environment deployment is limited to. In the format `Name;<BranchOrTag>` and policies delimited by `|`|

### Environment Secrets

The `gh environment secrets` command comprises of two subcommands, `list` and `create`, to access and create Environment specific Secrets.

```sh
$ gh environments secrets -h

List and Create Environment specific secrets in repositories.

Usage:
  environments secrets [command]

Available Commands:
  create      Create Environment secrets.
  list        Generate a report of Environment secrets.

Flags:
      --help   Show help for command

Use "environments secrets [command] --help" for more information about a command.
```

Both the `create` and `list` commands utilize the following fields: 

| Field Name | Description |
|:-----------|:------------|
|`RepositoryID`| The `ID` associated with the Repository, for API usage. |
|`RepositoryName` | The name of the repository where the data is extracted from. |
|`EnvironmentName`| The name of the repository specific environment. |
|`SecretName`| The name of the secret|
|`SecretValue`| Will be blank for `list`, and is required for `create` |
|`SecretCreatedAt`| The timestamp associated with when the secret was initially created. |
|`SecretUpdatedAt`| The timestamp associated with the last time the secret was modified. |

#### Create Secrets

The `gh environments secrets create` command will create secrets from a `csv` file using `--from-file` following the format outlined in [`gh environments secrets`](#environment-secrets).

>**Note**
> The `SecretValue` specified in the `csv` file will be [encrypted using the associated `public key`](https://docs.github.com/en/actions/security-guides/encrypted-secrets) before the environment secret is created.

```sh
$ gh environments secrets create -h

Create Environment secrets for specified environments per repository in an organization from a file.

Usage:
  environments secrets create <organization> [flags]

Flags:
  -d, --debug              To debug logging
  -f, --from-file string   Path and Name of CSV file to create secrets from
      --hostname string    GitHub Enterprise Server hostname (default "github.com")
  -t, --token string       GitHub personal access token for organization to write to (default "gh auth token")

Global Flags:
      --help   Show help for command
```

#### List Secrets

The `gh environments secrets list` command generates a `csv` report of environment specific secrets for the specified `<organization>` or `[repo ..]` list. If `[repo ...]` is specified, **secrets associated to environments across all repositories will be captured**. The report will contain secrets produces a `csv` report containing the fields outlined in [`gh environments secrets`](#environment-secrets).

>**Note**
> The `SecretValue` specified in the `csv` file will be left blank. **Secret values will NOT be extracted.**


```sh
$ gh environments secrets list -h

Generate a report of secrets for each environment per repository in an organization.

Usage:
  environments secrets list [flags] <organization> [repo ...] 

Flags:
  -d, --debug                To debug logging
      --hostname string      GitHub Enterprise Server hostname (default "github.com")
  -o, --output-file string   Name of file to write CSV report (default "report-20230512134718.csv")
  -t, --token string         GitHub Personal Access Token (default "gh auth token")

Global Flags:
      --help   Show help for command
```


### Environment Variables

The `gh environment variables` command comprises of two subcommands, `list` and `create`, to access and create Environment specific variables.

```sh
$  gh environments variables -h

List and Create Environment specific variables in repositories under an organization.

Usage:
  environments variables [command]

Available Commands:
  create      Create Environment variables.
  list        Generate a report of Environment variable.

Flags:
      --help   Show help for command

Use "environments variables [command] --help" for more information about a command.
```

Both the `create` and `list` commands utilize the following fields: 

| Field Name | Description |
|:-----------|:------------|
|`RepositoryID`| The `ID` associated with the Repository, for API usage. |
|`RepositoryName` | The name of the repository where the data is extracted from. |
|`EnvironmentName`| The name of the repository specific environment. |
|`VariableName`| The name of the variable|
|`VariableValue`| The value of the variable |
|`VariableCreatedAt`| The timestamp associated with when the variable was initially created. |
|`VariableUpdatedAt`| The timestamp associated with the last time the variable was modified. |

#### Create Variables

The `gh environments variables create` command will create variables from a `csv` file using `--from-file` following the format outlined in [`gh environments variables`](#environment-variables).



```sh
$ gh environments variables create -h

Create Environment variables for specified environments per repository in an organization from a file.

Usage:
  environments variables create <organization> [flags]

Flags:
  -d, --debug              To debug logging
  -f, --from-file string   Path and Name of CSV file to create variables from
      --hostname string    GitHub Enterprise Server hostname (default "github.com")
  -t, --token string       GitHub personal access token for organization to write to (default "gh auth token")

Global Flags:
      --help   Show help for command
```

#### List Variables

The `gh environments variables list` command generates a `csv` report of environment specific secrets for the specified `<organization>` or `[repo ..]` list. If `[repo ...]` is specified, **variables associated to environments across all repositories will be captured**. The report will contain variables produces a `csv` report containing the fields outlined in [`gh environments variables`](#environment-variables).


```sh
$ gh environments variables list -h

Generate a report of variables for each environment per repository in an organization.

Usage:
  environments variables list [flags] <organization> [repo ...] 

Flags:
  -d, --debug                To debug logging
      --hostname string      GitHub Enterprise Server hostname (default "github.com")
  -o, --output-file string   Name of file to write CSV report (default "report-20230512135332.csv")
  -t, --token string         GitHub Personal Access Token (default "gh auth token")

Global Flags:
      --help   Show help for command
