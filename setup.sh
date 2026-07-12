#!/bin/bash
# setup.sh - Build nsbox and symlink trae-sandbox -> nsbox
# Safe to re-run after trae updates: it will re-compile and re-link.
#
# Usage:  bash setup.sh
#   Flags:
#     --revert    Restore original trae-sandbox from backup
#     --dry-run   Show what would be done without making changes

set -euo pipefail

# ── paths ──────────────────────────────────────────────────────
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
SANDBOX_DIR="/usr/share/trae-cn/resources/app/modules/sandbox"
NSBOX_SRC="${PROJECT_DIR}/nsbox.go"
NSBOX_BIN="${PROJECT_DIR}/nsbox"
TARGET="${SANDBOX_DIR}/trae-sandbox"
BACKUP="${SANDBOX_DIR}/trae-sandbox.bak"

REVERT=false
DRY_RUN=false

for arg in "$@"; do
    case "$arg" in
        --revert)  REVERT=true  ;;
        --dry-run) DRY_RUN=true ;;
    esac
done

# ── helpers ────────────────────────────────────────────────────
info()  { echo "[info]  $*"; }
warn()  { echo "[warn]  $*" >&2; }
die()   { echo "[error] $*" >&2; exit 1; }

# ── revert mode ────────────────────────────────────────────────
if $REVERT; then
    if [ -f "$BACKUP" ]; then
        if $DRY_RUN; then
            info "[dry-run] Would remove $TARGET (symlink) and restore $BACKUP -> $TARGET"
        else
            info "Removing $TARGET (symlink)"
            rm -v "$TARGET"
            info "Restoring $BACKUP -> $TARGET"
            mv -v "$BACKUP" "$TARGET"
            info "Reverted. trae-sandbox is back to original binary."
        fi
    else
        die "No backup found at $BACKUP. Nothing to revert."
    fi
    exit 0
fi

# ── normal mode: build + link ─────────────────────────────────

# 1. Check source
[ -f "$NSBOX_SRC" ] || die "nsbox.go not found at $NSBOX_SRC"

# 2. Compile
info "Compiling nsbox.go -> nsbox ..."
if $DRY_RUN; then
    info "[dry-run] Would run: go build -o $NSBOX_BIN $NSBOX_SRC"
else
    go build -o "$NSBOX_BIN" "$NSBOX_SRC" || die "Compilation failed"
    info "Compiled: $NSBOX_BIN ($(du -h "$NSBOX_BIN" | cut -f1))"
fi

# 3. Ensure target exists
if [ ! -e "$TARGET" ] && [ ! -L "$TARGET" ]; then
    die "trae-sandbox not found at $TARGET"
fi

# 4. Handle existing state
if [ -L "$TARGET" ]; then
    LINK_TARGET="$(readlink -f "$TARGET")"
    if [ "$LINK_TARGET" = "$NSBOX_BIN" ]; then
        info "$TARGET already points to nsbox. Link is up to date."
        exit 0
    fi
    if $DRY_RUN; then
        info "[dry-run] Would remove existing symlink $TARGET"
    else
        info "Removing old symlink $TARGET -> $LINK_TARGET"
        rm -v "$TARGET"
    fi
fi

# 5. Backup original file if still present
if [ -f "$TARGET" ]; then
    if $DRY_RUN; then
        info "[dry-run] Would move $TARGET -> $BACKUP"
    else
        info "Backing up original: $TARGET -> $BACKUP"
        mv -v "$TARGET" "$BACKUP"
    fi
fi

# 6. Create symlink
if $DRY_RUN; then
    info "[dry-run] Would create symlink $TARGET -> $NSBOX_BIN"
else
    info "Creating symlink: $TARGET -> $NSBOX_BIN"
    ln -sv "$NSBOX_BIN" "$TARGET"
    info "Done. trae-sandbox now routes through nsbox (no sandbox)."
fi
