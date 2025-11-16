# bsubio jobs

List recent jobs with status, timestamps, and duration

## Usage

```
bsubio jobs [options]
```

## Options

- `--status <status>` - Filter by status (pending, claimed, finished, failed)
- `--limit <n>` - Limit number of results

## Output Columns

- JOB ID - Unique job identifier
- TYPE - Job processing type
- STATUS - Current job status
- CREATED AT - Job submission timestamp
- FINISHED AT - Completion timestamp (finished/failed jobs only)
- TOOK (s) - Duration in seconds from creation to finish

## Examples

List all recent jobs:
```
bsubio jobs
```

List only failed jobs:
```
bsubio jobs --status failed
```

List last 10 jobs:
```
bsubio jobs --limit 10
```
