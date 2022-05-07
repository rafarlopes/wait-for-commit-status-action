# wait-for-commit-status-action
GitHub Action to wait for the commit status with a given context

# Usage

This action uses an environment variable name `GITHUB_TOKEN` to authenticate and checkout the repository with the default token provided by GitHub Actions.

We also use the `GITHUB_REPOSITORY` with the default value of the current repository where this actions runs.

<!-- start usage -->
```yaml
- uses: rafarlopes/wait-for-commit-status-action@v1
  with:
    # Context for which we should look for the matching status
    context: 'cd/my-web-api/development'

    # The commit sha we should look for the status
    sha: 'ead549b4ab21b7d6653556b2772c2338f11a3082'
```

Example with overriden environment variables in case of different repository or private repository with PAT:

```yaml
- uses: rafarlopes/wait-for-commit-status-action@v1
  env:
    GITHUB_REPOSITORY: 'myorg/myprivaterepo'
    GITHUB_TOKEN: ${{ secrets.MY_PAT }}
  with:
    # Context for which we should look for the matching status
    context: 'cd/my-web-api/development'

    # The commit sha we should look for the status
    sha: 'ead549b4ab21b7d6653556b2772c2338f11a3082'
```
<!-- end usage -->
