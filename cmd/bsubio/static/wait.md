# bsubio wait

Wait for a job to complete

## Usage

```
bsubio wait [options] <jobid>
```

## Options

- `-v` - Verbose output
- `-t <seconds>` - Polling interval in seconds (default: 5)

## Arguments

- `jobid` - Job ID to wait for

## Examples

Wait for a job to complete:
```
bsubio wait job_abc123
```

Wait with verbose output:
```
bsubio wait -v job_abc123
```

Wait with custom polling interval:
```
bsubio wait -t 10 job_abc123
```
