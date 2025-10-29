# bsubio cancel

Cancel a job or all jobs

## Usage

```
bsubio cancel [options] [jobid]
```

## Options

- `-a` - Cancel all pending/claimed jobs

## Arguments

- `jobid` - Job ID to cancel (not required with -a)

## Examples

Cancel a specific job:
```
bsubio cancel job_abc123
```

Cancel all pending/claimed jobs:
```
bsubio cancel -a
```
