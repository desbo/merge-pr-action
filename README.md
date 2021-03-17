# Merge PR action

This action merges PRs from automatic dependency upgrade services.

If the PR title includes two SemVer version numbers, and the type of update (patch, minor or major) is allowed by the action configuration, it'll be merged. 

Intended to be included in a workflow that builds and tests the project, to be run as a separate job after these steps have passed successfully.

You should also include a condition in the merge job to only run against PRs created by your depenency bot (e.g. `if: github.actor == 'some-bot'` in usage example).

## Use with scala-steward
Update your `.scala-steward.conf` so that PR titles include both the old and new version number, for example:

```
commits.message = "Upgrade ${artifactName} from ${currentVersion} to ${nextVersion}"
```

## Caveats and future work
At the moment the action only runs in response to [`pull_request` events](https://docs.github.com/en/actions/reference/events-that-trigger-workflows#pull_request). This means that, for the merge to take place, the check suite triggered by the PR being opened must pass and if this isn't the case (for example if a test fails), the merge action won't run, and the PR will need to be merged manually.

It would be possible to update the action to also run on successful [`check_suite` events](https://docs.github.com/en/actions/reference/events-that-trigger-workflows#check_suite); this would require some changes to access the PR title, which is not included in the event payload for non-`pull_request` events.

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
