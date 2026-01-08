#!/bin/bash
#
# Sync regexes and fixtures from upstream Matomo device-detector
# Usage: ./scripts/sync-upstream.sh [commit-hash]
#
# If no commit hash is provided, syncs to the latest master branch

set -e

UPSTREAM_REPO="https://github.com/matomo-org/device-detector.git"
TEMP_DIR=$(mktemp -d)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Matomo Device Detector Sync ===${NC}"
echo "Project root: $PROJECT_ROOT"
echo "Temp directory: $TEMP_DIR"

# Cleanup on exit
cleanup() {
    echo -e "${YELLOW}Cleaning up temp directory...${NC}"
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Clone upstream repo (shallow clone for speed)
echo -e "${GREEN}Cloning upstream repository...${NC}"
if [ -n "$1" ]; then
    # Clone full repo if specific commit requested
    git clone --quiet "$UPSTREAM_REPO" "$TEMP_DIR/device-detector"
    cd "$TEMP_DIR/device-detector"
    git checkout --quiet "$1"
    COMMIT_HASH="$1"
else
    # Shallow clone for latest
    git clone --depth 1 --quiet "$UPSTREAM_REPO" "$TEMP_DIR/device-detector"
    cd "$TEMP_DIR/device-detector"
    COMMIT_HASH=$(git rev-parse HEAD)
fi

COMMIT_DATE=$(git log -1 --format="%ci" HEAD)
COMMIT_MSG=$(git log -1 --format="%s" HEAD)

echo -e "${GREEN}Upstream commit:${NC} $COMMIT_HASH"
echo -e "${GREEN}Commit date:${NC} $COMMIT_DATE"
echo -e "${GREEN}Commit message:${NC} $COMMIT_MSG"

# Sync regexes
echo -e "${GREEN}Syncing regexes...${NC}"
rm -rf "$PROJECT_ROOT/regexes"
cp -r regexes "$PROJECT_ROOT/"

# Count regex files
REGEX_COUNT=$(find "$PROJECT_ROOT/regexes" -name "*.yml" | wc -l | tr -d ' ')
echo -e "  Copied ${YELLOW}$REGEX_COUNT${NC} regex files"

# Sync fixtures
echo -e "${GREEN}Syncing fixtures...${NC}"
rm -rf "$PROJECT_ROOT/fixtures"
mkdir -p "$PROJECT_ROOT/fixtures"
cp -r Tests/fixtures/*.yml "$PROJECT_ROOT/fixtures/"

# Count fixture files
FIXTURE_COUNT=$(find "$PROJECT_ROOT/fixtures" -name "*.yml" | wc -l | tr -d ' ')
echo -e "  Copied ${YELLOW}$FIXTURE_COUNT${NC} fixture files"

# Record the sync info
cat > "$PROJECT_ROOT/UPSTREAM_SYNC.md" << EOF
# Upstream Sync Information

This project syncs regexes and fixtures from the [Matomo device-detector](https://github.com/matomo-org/device-detector) PHP library.

## Last Sync

- **Commit:** ${COMMIT_HASH}
- **Date:** ${COMMIT_DATE}
- **Message:** ${COMMIT_MSG}
- **Synced on:** $(date -u +"%Y-%m-%d %H:%M:%S UTC")

## What's Synced

- \`regexes/\` - All regex YAML files (${REGEX_COUNT} files)
- \`fixtures/\` - Test fixtures from \`Tests/fixtures/\` (${FIXTURE_COUNT} files)

## How to Sync

\`\`\`bash
# Sync to latest
./scripts/sync-upstream.sh

# Sync to specific commit
./scripts/sync-upstream.sh <commit-hash>
\`\`\`
EOF

echo -e "${GREEN}=== Sync Complete ===${NC}"
echo ""
echo "Summary:"
echo "  - Regex files: $REGEX_COUNT"
echo "  - Fixture files: $FIXTURE_COUNT"
echo "  - Upstream commit: ${COMMIT_HASH:0:8}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Run tests: go test ./..."
echo "  2. Review changes: git diff"
echo "  3. Commit: git add -A && git commit -m 'Sync with upstream (${COMMIT_HASH:0:8})'"
