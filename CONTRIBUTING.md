# Contributing to cbmultimanager

First of all, thanks for taking the time to contribute! 👍 The following documentation is a set of guidelines for contributing to cbmultimanager.

## Reporting Bugs and Requesting Enhancements

Issues and enhancements for cbmultimanager are tracked in the [CMOS project](https://issues.couchbase.com/browse/CMOS) on Jira. If you have found a bug or have an idea to improve cbmultimanager, please don't keep it to yourself, even if you don't intend to write the code for it yourself! File an issue and someone will look at it.

## Code Contribution Workflow

To start, follow the setup instructions in README.md.

Please use a descriptive Git branch name - ideally every change should be associated with a CMOS issue, so name your branch something like `CMOS-1234`.

Please follow this format for your commit messages, as this means they will automatically be associated with Jira tickets when you submit them for review:

```
CMOS-1234: Fix the frobnicator

A longer description of the background to this commit and the issue it fixes. Sometimes commit messages can be longer than the actual code changes, in which case having the additional context is invaluable.

Write your commit message in the imperative: "Fix bug" and not "Fixed bug" or "Fixes bug."  This convention matches up with commit messages generated by commands like git merge and git revert.
```

## Pre-Commit Checklist

Before committing your changes, it's a good idea to run through this checklist. All of this is also enforced by automated Commit Validation, so you will not be able to submit your changes unless this passes, but running it locally may save you time.

Back-end (`cluster-manager` and `agent`):

* Are all linters happy? (`make lint`, you may need to [install golangci-lint](https://golangci-lint.run/usage/install/#local-installation))
  * `golangci-lint run` might tell you to run gofumpt, run `go install mvdan.cc/gofumpt@latest` to install it, and then `gofumpt -w **/*.go` to run it.
* Do all tests pass? (`go test ./...`, check there are no failures)
* If you have added a new health checker:
  * Have you given it an ID and added it to the [spreadsheet](https://docs.google.com/spreadsheets/d/15ylgauQTtPyxklesgcmxsDJYUrQhGVcK3AaDtj1BpNQ/edit#gid=0)?
  * Have you added an entry for it in `docs/modules/ROOT/pages/checkers.adoc`? Follow the lead of existing documentation.
    * If there are also Prometheus/Loki rules in couchbaselabs/observability that can trigger that checker, instead document it there and add a comment with its ID to `checkers.adoc` (otherwise CV will fail)

Front-end (`ui`):

* Does it build? (`npm run build`)

Container:

* Does it build? (`docker build .`)
* Does it run (following the instructions in README.md), can you access the UI (localhost:7196), and monitor a cluster?

## Code Review

Don't push your code to master. Instead, get it reviewed through Gerrit.

To get it set up, follow the [Gerrit set-up guide](https://hub.internal.couchbase.com/confluence/display/CR/Contributing+Changes+via+Gerrit) on Confluence.

Once you have done so, add Gerrit as a Git remote:

```
git remote add gerrit ssh://review.couchbase.org:29418/cbmultimanager.git
```

And push your branch up to create a review:

```
git push gerrit HEAD:refs/for/master
```

The message it prints out will link to the newly created review. (If it gives you an error about no `Change-Id`, follow the instructions in the message to set it up.)

Once your code is on Gerrit, it will need to be tested by a robot and reviewed by a human (one of the other developers on this project). This is tracked through the `Verified` and `Code-Review` labels respectively. To allow you to "submit" (merge) your code, it will need a Verified +1, a Code-Review +2, and no Code-Review -2 labels (more information on the meanings of the labels is in the [Gerrit docs](https://gerrit-review.googlesource.com/Documentation/config-labels.html)). Once you have the sign-offs, simply click "Submit" to merge your code!

If you get a message requesting changes, don't panic, simply fix up your code, run `git commit --amend`, push it again, and repeat the process until it's ready for merge.

## Adding a checker to CMOS

Visit the [Checker spreadsheet](https://docs.google.com/spreadsheets/d/15ylgauQTtPyxklesgcmxsDJYUrQhGVcK3AaDtj1BpNQ/edit#gid=0) to find a checker to work on. It should be a high priority, not in CMOS, and one which either doesn't already have a Jira CMOS ticket, or has one but isn't being worked on. Create a CMOS ticket for this checker and assign it to yourself, link that ticket to the row in the spreadsheet that corresponds to the checker. Give the checker the next available ID.

Firstly, a quick bit of terminology. A checker examines the state of your cluster/node/bucket and if needed sends an alert. An alert is the notification that the user of CMOS will receive. Checkers are stored in cbmultimanager and metrics and log and metric (Loki and Prometheus) alerts in observability. The default application to write checkers in is typically cbmultimanager. Write in observability when Loki or Prometheus metrics are needed. It's currently impossible to access `memcached.log` without using Loki.

---

### Adding a checker to Cbmultimanager.

1. You need to determine if the checker is on a node, cluster, or bucket level. Call this x, and consider the pair of files:
	`x_checkers.go` and `x_checkers_test.go`


2. Write the checker in `status/x_checkers.go`. Each checker is a function that takes a struct representing x. If the checker determines that something is wrong, then it should set the checker status to the appropriate value. Possible checker statuses are:
	`GoodCheckerStatus` - This is the default.
	`WarnCheckerStatus` - issues that warrant review and have the potential to cause application degradation, but do not need an immediate response.
	`AlertCheckerStatus` - issues that need immediate actioning to prevent damage or application degradation.
	`InfoCheckerStatus` - issues that warrant review, or configuration that does not follow best practice.
	There is a bit of boiler plate in each checker that you can copy from the other checkers in `status/x_checkers.go`


3. Create a test for the checker in `status/x_checkers_test.go`. These tests create simulated environments to run your checker against. You want to create environments that test each possible outcome of your checker. For example, the Auto-Failover checker requires two tests, one where Auto-Failover is enabled, and one disabled.


4. Add a definition to `values/checker_defs.go`. Define a string at the bottom of the const for your checker. The name of the variable should be `Result.Name` from your checker. The value of the variable should make it easy to undestand what the checker checks. Then at the bottom of the document create a new definition, following the convention of the definitions previously added.


5. Add a key/value pair to the `allCheckerFns` map in `status/checkers.go.` The key should be values.y where y is Result.Name from your checker. The key should be the name of the checker function.


6. Add the relevant information to `docs/checkers.adoc`, get the checker ID from the checkers spreadsheet.


#### How to test your changes

By default, if you try to run CMOS, the master branch of cbmultimanager will be downloaded. This will mean that you can't see your changes. In order to see them you need to follow the steps outlined in our [cbmultimanager readme](https://github.com/couchbaselabs/workbench-prototype/blob/master/README.md).