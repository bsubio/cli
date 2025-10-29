# bsubio rm

Delete a job or all jobs

## Usage

```
bsubio rm [options] [jobid]
```

## Options

- `-a` - Delete all jobs

## Arguments

- `jobid` - Job ID to delete (not required with -a)

## Examples

Delete a specific job:
```
bsubio rm job_abc123
```

Delete all jobs:
```
bsubio rm -a
```
