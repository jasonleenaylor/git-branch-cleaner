# Git Branch Cleaner

**Git Branch Cleaner** is a Go-based command-line tool designed to help developers efficiently clean up local git branches. With flexible options, it enables removing unnecessary branches while retaining key ones based on your preferences. Branches are deleted with prejudice, ignoring their merge status.

## Features

- **Remove all branches except the current branch** using the `--all` flag.
- **Standard cleanup** using `--standard`
Retain only the following branches:
  - `master`
  - `main`
  - `develop`
  - `release/**`
- **Exclude additional branches**: Use the `--exclude` flag (can be specified multiple times) to protect specific branches during cleanup.

## Installation

### Using `go install`

Install Git Branch Cleaner directly from the repository:

```bash
go install github.com/jasonleenaylor/git-branch-cleaner@latest

