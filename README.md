# Firebase Functions SDK for Go
Welcome to the Firebase Functions SDK for Go. This is an experimental
repository for a project that may never be officially supported.
It is in early stages, there is absolutely no guarantees for
API stability and breaking changes should be expected often.

Instructions for use will be granted upon joining the Firebase
trusted testers group.

## How it works

Developers write their own module using this package. The
program at `/support/codegen` analyzes a package and generates
an emulator binary combining the package and `/support/emulator`.

The `firebase-tools` CLI can use this emulator binary to
serve or analyize a package for local development or deploying
to production.
