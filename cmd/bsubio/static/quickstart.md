# bsubio Quickstart Guide

Welcome to bsubio! This guide will help you get started quickly.

## Setup

1. Configure your API key:

    bsubio config

2. List available job types:

    bsubio types

## Basic Workflow

### Submit a Job

Submit a job for processing:

    bsubio submit pdf/extract Simple.pdf

You'll receive a job ID like `job_abc123`.

### Check Job Status

Monitor your job:

    bsubio status job_abc123

### Wait for Completion

Wait for the job to finish:

    bsubio wait job_abc123

### Get Results

Retrieve the output:

    bsubio cat job_abc123

Or check the logs:

    bsubio logs job_abc123

## Quick Submit and Wait

Submit a job and automatically wait for results:

    bsubio submit -w pdf/extract Simple.pdf

Save output to a file:

    bsubio submit -w -o output.json pdf/extract Simple.pdf

## Managing Jobs

List recent jobs:

    bsubio jobs

Cancel a job:

    bsubio cancel job_abc123

Delete a job:

    bsubio rm job_abc123

## Getting Help

View help for any command:

    bsubio help <command>

For example:

    bsubio help submit
    bsubio help wait

## Next Steps

- Use `bsubio help <command>` to learn more about specific commands
- Check server version with `bsubio version`
- Explore available job types with `bsubio types`
