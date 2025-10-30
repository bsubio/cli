# bsubio jobs

List recent jobs

## Usage

```
bsubio jobs [options]
```

## Options

- `--status <status>` - Filter by status (pending, claimed, finished, failed)
- `--limit <n>` - Limit number of results

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
