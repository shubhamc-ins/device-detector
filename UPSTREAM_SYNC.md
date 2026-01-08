# Upstream Sync Information

This project syncs regexes and fixtures from the [Matomo device-detector](https://github.com/matomo-org/device-detector) PHP library.

## Last Sync

- **Commit:** f674b94c373590f15e0a861894d8c529ccc248fe
- **Date:** 2026-01-07 12:20:09 +0000
- **Message:** Improves version detection for Fire OS (#8208)
- **Synced on:** 2026-01-08 11:50:33 UTC

## What's Synced

- `regexes/` - All regex YAML files (20 files)
- `fixtures/` - Test fixtures from `Tests/fixtures/` (83 files)

## How to Sync

```bash
# Sync to latest
./scripts/sync-upstream.sh

# Sync to specific commit
./scripts/sync-upstream.sh <commit-hash>
```
