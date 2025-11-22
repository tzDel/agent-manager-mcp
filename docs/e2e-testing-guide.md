# E2E Testing Guide: Worktree Use Cases

Quick reference for end-to-end testing of create and remove worktree workflows.

---

## Test Setup

Initialize a test repository:

```bash
# Create test repository
mkdir test-repo && cd test-repo
git init
echo "Initial commit" > README.md
git add . && git commit -m "Initial commit"
```

---

## Testing Create Worktree Workflow

**Test Case:** Create worktree for agent "copilot"

**Execute:**
```bash
# Via MCP server or direct use case call
# Creates: ./worktrees/agent-copilot
# Creates branch: agent-copilot
```

**Verify:**
```bash
# 1. Confirm worktree exists
git worktree list
# Expected: Shows worktrees/agent-copilot with branch agent-copilot

# 2. Confirm branch was created
git branch
# Expected: Shows agent-copilot branch

# 3. Verify worktree is on correct branch
git -C worktrees/agent-copilot branch --show-current
# Expected: agent-copilot

# 4. Check directory exists
ls -la worktrees/
# Expected: Shows agent-copilot directory

# 5. Verify worktree has repository content
ls worktrees/agent-copilot/
# Expected: Shows README.md and other repo files
```

**Edge Cases to Test:**
```bash
# Attempt duplicate creation (should fail)
# Create with invalid agent ID (should fail)
# Create when worktrees/ directory doesn't exist (should auto-create)
```

---

## Testing Remove Worktree Workflow

**Test Case:** Remove worktree for agent "copilot"

**Execute:**
```bash
# Via MCP server or direct use case call
# Removes: ./worktrees/agent-copilot
# Optionally deletes branch: agent-copilot
```

**Verify:**
```bash
# 1. Confirm worktree was removed
git worktree list
# Expected: Does NOT show agent-copilot worktree

# 2. Confirm branch was deleted (if applicable)
git branch
# Expected: Does NOT show agent-copilot branch

# 3. Verify directory was removed
ls worktrees/
# Expected: Does NOT show agent-copilot directory

# 4. Check for stale references
git worktree prune --dry-run
# Expected: No output (no stale worktrees to prune)
```

**Edge Cases to Test:**
```bash
# Remove with uncommitted changes (should fail or force)
git -C worktrees/agent-copilot status
# Add test file: echo "test" > worktrees/agent-copilot/test.txt

# Remove non-existent worktree (should fail gracefully)

# Remove while worktree is in use (locked files - should fail)
```

---

## Complete E2E Test Flow

Full workflow testing both operations:

```bash
# Setup
git init test-repo && cd test-repo
echo "test" > README.md && git add . && git commit -m "init"

# Test 1: Create
# [Execute create worktree for agent "copilot"]
git worktree list                                    # Verify exists
git branch                                            # Verify branch
git -C worktrees/agent-copilot status                # Verify clean state

# Test 2: Make changes in worktree
echo "agent work" > worktrees/agent-copilot/feature.txt
git -C worktrees/agent-copilot add .
git -C worktrees/agent-copilot commit -m "Add feature"
git -C worktrees/agent-copilot log --oneline -1      # Verify commit

# Test 3: Remove
# [Execute remove worktree for agent "copilot"]
git worktree list                                    # Verify removed
git branch                                            # Verify branch deleted
ls worktrees/                                         # Verify directory gone

# Cleanup
cd .. && rm -rf test-repo
```

---

## Quick Verification Checklist

**After Create Worktree:**
- [ ] `git worktree list` shows new worktree
- [ ] `git branch` shows new branch
- [ ] Directory exists at expected path
- [ ] Worktree contains repository files

**After Remove Worktree:**
- [ ] `git worktree list` does NOT show worktree
- [ ] `git branch` does NOT show branch (if deleted)
- [ ] Directory does NOT exist
- [ ] `git worktree prune --dry-run` shows no stale refs

---

## Key Git Commands Reference

```bash
git worktree list                        # List all worktrees
git worktree add <path> -b <branch>     # Create worktree
git worktree remove <path>              # Remove worktree
git worktree prune                      # Clean stale references
git -C <path> <command>                 # Run command in worktree
git branch                              # List branches
git branch -d <name>                    # Delete branch (safe)
git branch -D <name>                    # Delete branch (force)
```
