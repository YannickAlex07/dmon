# Releasing dmon

The release process of **dmon** is quite simple and mostly automated. Just create a release from the release tab that has a tag that starts with `v` (i.e. `v1.0.0`).

Then [`goreleaser`](https://goreleaser.com) will automatically take over to build the package and attaches it to the newly created release. Furthermore it will also fill the release overview with a changelog.