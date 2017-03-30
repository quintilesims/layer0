# Contributing to Layer0

This page contains information about reporting issues as well as some tips and
guidelines for potential open source contributors. Make sure you read our [community guidelines](#layer0-community-guidelines) before participating.

If you're interested in hacking on Layer0, please also see our [README.md](README.md),
which details the various components of Layer0.

Our contributor documentation is based largely on the good work at
[Docker](https://github.com/docker/docker). Thanks to the Docker team for helping
to establish a solid OSS community!

## Topics

* [Reporting Security Issues](#reporting-security-issues)
* [Reporting Issues](#reporting-other-issues)
* [Quick Contribution Tips and Guidelines](#quick-contribution-tips-and-guidelines)
* [Community Guidelines](#layer0-community-guidelines)

## Reporting security issues

Please **DO NOT** file a public issue, instead send your report(s) privately to
[carbon@us.imshealth.com](mailto:carbon@us.imshealth.com).

## Reporting other issues

The easiest way to contribute to Layer0 is to send us a detailed report when you
encounter an issue.

Check that [our issue database](https://github.com/quintilesims/layer0/issues)
doesn't already include your problem or suggestion before submitting an issue.
If you find a match, you can use the "subscribe" button to get notified on
updates.

When reporting issues, include (where applicable):

* The output of `l0 --version`.
* The output of `l0 admin debug`.
* The output of `l0-setup --version`

Also include the steps required to reproduce the problem if possible.
Exceptionally long log output should be posted as a gist (https://gist.github.com).
Don't forget to remove sensitive data from your logfiles before posting (you can
replace those parts with "REDACTED").

## Quick contribution tips and guidelines

This section gives the experienced contributor some tips and guidelines.

### Pull requests are being accepted 🆒

The maintainers of Layer0 **love** pull requests. If you make a good one that
follows our brief guidelines, you'll be forever immortalized in our commit history.

### Conventions

Fork the repository and make changes on your fork in a feature branch:

* If it's a bug fix branch, name it XXXX-something where XXXX is the number of
	the issue.
* If it's a feature branch, create an enhancement issue to announce
	your intentions, and name it XXXX-something where XXXX is the number of the
	issue.

Depending on the changes you're proposing, you may need to update the follow tests:

* [Smoketests](/tests/smoke/README.md)
* Unit tests
* System tests

[Update the documentation](https://github.com/quintilesims/layer0/tree/develop/docs-src)
when creating or modifying features. Test your documentation changes for
clarity, concision, and correctness, as well as a clean documentation build via
`make build`.

Pull request descriptions should be as clear as possible and include a reference
to all the issues that they address.

### Successful Changes

* Make sure that your PR is directly related to an existing issue. If an issue doesn't
already exist, make a new one.

* Keep PRs small and concise. If you must make a large changeset, we can discuss
how to proceed in the PRs' associated issue.

## Layer0 community guidelines

* Be respectful. We appreciate courteous and polite community members; snarkiness
and soapboxing are highly discouraged. We are all here to make Layer0 better.

* Don't break the law by posting another company's assets, unlicensed Cat gifs, etc.

* Stay on topic.

Violating the community guidelines will result in being blocked from the QuintilesIMS
Github organization.
