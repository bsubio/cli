# bsubio status

Show detailed job status

## Usage

```
bsubio status <jobid>
```

## Arguments

- `jobid` - Job ID

## Examples

Show job status:
```
bsubio status job_abc123
```

## Output

Displays detailed information including:
- Job ID
- Job Type
- Status (pending, claimed, finished, failed)
- Data Size
- Created At
- Claimed At (if applicable)
- Finished At (if applicable)
- Claimed By (if applicable)
- Error Message (if failed)
