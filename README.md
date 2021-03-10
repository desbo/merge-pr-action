# Merge PR action

This action merges PRs from automatic dependency upgrade services.

If the PR title includes two SemVer version numbers, and the type of update (patch, minor or major) is allowed by the action configuration, it'll be merged. 

Intended to be included in a workflow that builds and tests the project, to be run as a separate job after these steps have passed successfully.

You should also include a condition in the merge job to only run against PRs created by your depenency bot (e.g. `if: github.actor == 'some-bot'` in usage example).

## Inputs
### `GITHUB_TOKEN`
The token of a GitHub user with `repo` access (required to merge PRs). This should be provided by a secret, of course.

### `ALLOWED_UPDATE`
Set to either `patch`, `minor` or `major` to control the type of upgrade allowed. Defaults to `patch`. 

### `MERGE_METHOD`
The [merge method](https://docs.github.com/en/github/administering-a-repository/about-merge-methods-on-github) to use: `merge`, `squash` or `rebase`. Defaults to `merge`.

### Example usage

```yaml
jobs:
  build:
    name: Build and test
    runs-on: ubuntu-18.04
    steps:
    - name: Check out project
      uses: actions/checkout@v2

    - name: Test
      run: sbt test

  merge:
    name: Merge dependency update
    if: github.actor == 'some-bot'
    needs:
      - build
    runs-on: ubuntu-latest
    steps:
    - name: merge PR
      uses: desbo/merge-pr-action@v0
      with:
        GITHUB_TOKEN: ${{ secrets.CI_GITHUB_TOKEN }}
        ALLOWED_UPDATE: minor
        MERGE_METHOD: rebase
```
