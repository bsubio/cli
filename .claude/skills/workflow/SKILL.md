---
name: proper workflow
description: Follows project workflow for GitHub tickets, issues and branch naming, as well as tags, and tying issue with PRs
allowedTools:
  - "Bash(gh:*)"
  - "Bash(git checkout:*)"
  - "Bash(git push:*)"
  - "Bash(git branch:*)"
  - "Bash(git status:*)"
---

# Project GitHub Workflow skill

## Instructions

This skill instructs how to properly create GitHub issues, branches and pull requests for the project.

When asked to create the ticket, feed the issue title and the body to the following command:

    gh issue create -b "<body>" --title "<title>"

where:

    <body>
    <title>

are body of the issue and title, respectively.

When work on the project ticket, pull the GitHub issue with this CLI command:

    gh issue list

When asked to work on the issue, make an appropriate branch with the following command:

    git checkout -b issue/<issue_number>
    git push -u origin issue/<issue_number>

When asked to make PR:

    gh pr create --title "<title>" --body "<body>"

there:

    <title>

is the "Feat: " or "Fix: " + explanation of what was done and "(Fixes: #issuenum)" added, and

    <body>

being explanation of what we made.

