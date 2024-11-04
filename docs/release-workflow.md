# Release Workflows 

We manage code changes with Git and [GitHub Action](https://github.com/features/actions) to automate release workflows.

## Git Flow

* `main` is the protected branch of our repository. It contains the latest code.
* All changes to `main` must be made via Pull Request.
* Merge commit is disabled, we use rebase and squash merge strategies.
* Write your commit messages with [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/).
* Changes to `main` branch are automatically deployed to `dev` environment.
