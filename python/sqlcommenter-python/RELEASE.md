# Releases

This page describes the release process. The release will be done by the maintainers.

## How to cut an individual release

### Branch management and versioning strategy

We use [Semantic Versioning](https://semver.org/).

We maintain a separate branch for each  release, named `release-<major>.<minor>`, e.g. `release-1.1`, `release-2.0`.

The usual flow is to merge new features and changes into the master branch and to merge bug fixes into the latest release branch. Bug fixes are then merged into master from the latest release branch. The master branch should always contain all commits from the latest release branch. As long as master hasn't deviated from the release branch, new commits can also go to master, followed by merging master back into the release branch.

If a bug fix got accidentally merged into master after non-bug-fix changes in master, the bug-fix commits have to be cherry-picked into the release branch, which then have to be merged back into master. Try to avoid that situation.

Maintaining the release branches for older minor releases happens on a best effort basis.

### Prepare your release

This helps ongoing PRs to get their changes in the right place, and to consider whether they need cherry-picked.

1. Make a PR to update `CHANGELOG.md` on master
   - Add a new section for the new release so that "unreleased" is blank and at the top.
   - New section should say "## x.y.0 - release date".
2. This PR will get merged.

### Create release branch

To prepare release branch, first create new release branch (release-X.Y) in repository from master commit of your choice

### Publish a stable release

To publish a stable release:

1. Do not change the release branch directly; make a PR to the release-X.Y branch with VERSION and any CHANGELOG changes.
2. After merging your PR to release branch, `git tag` the new release from release branch.
3. Merge the release branch `release-x.y` to `master`