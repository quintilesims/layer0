# Layer0 Release Process

This document describes the release process for Layer0.


## Release Branching Worklow

Layer0 follows the branching workflow as described by the [Git Flow](http://danielkummer.github.io/git-flow-cheatsheet/) model. This model ensures that the `master` branch is always in a production-ready state.

The basic branching workflow is:
```
<feature branches> -> develop -> release -> master
```


## Pre-release Checklist

Before starting the release process, make sure to check the following items:

* Have all the necessary [pull requests](https://github.com/quintilesims/layer0/pulls) been merged?
* Does a [release](https://github.com/quintilesims/layer0/branches) branch already exist? It should have been removed after the last release.
* Do you have the [layer0](https://github.com/quintilesims/layer0) repository cloned locally?


# Release Process

The following section describes the entire workflow for releasing a new Layer0 version.


## Create the `release` branch

Run these commands in the your local `layer0` repository to create the `release` branch:

```
# Checkout latest version of the develop branch
git fetch origin
git checkout remotes/origin/develop

# Create the release branch and merge it with master
git checkout -b release
git merge remotes/origin/master

# Push the release branch to gitlab
git push -u origin release
```

## Merge Release into Master

With the release notes and documentation updated, [create a pull request](https://github.com/quintilesims/layer0/compare) from the `release` branch targeting the `master` branch (_base: master, compare: release_). 
Merge using **squash and merge.**
Wait for the unit tests and smoketests to pass before merging. 
Once the pull request has finished, make sure to delete the `release` branch. 


## Add Version Tag

Once the pull request has been merged, add a new version tag:

```
# Fetch the updated master branch
git fetch origin
git checkout remotes/origin/master

# Add and push the version tag
git tag -a vX.X.X -m "<some message about the version>"
git push origin --tag
```

## Add Release Notes

Add release notes to the release with this format.
```
### Features
* Features here.

###Fixes
* Bug fixes here.
```
This can be done through the [Github UI](https://github.com/quintilesims/layer0/releases) by clicking the edit button on the latest release.
Please follow the existing formatting when adding release notes.


# Announce the release

Once the release is available for download, send a message to
[#carbon](https://ims-dev.slack.com/messages/carbon) with a link and a one-line
summary of the release contents. For example:
```
@here Layer0 has released v0.9.0 https://github.com/quintilesims/layer0/releases - which aims to be a stable version of our load balancer support.
```


# Merge Master into Develop

To bring the `develop` branch up-to-date with `master`, create a [create a pull request](https://github.com/quintilesims/layer0/compare) from the `master` branch targeting the `develop` branch (_base: develop, compare: master_). **Please rebase and merge** (squash and merge if not possible.)
