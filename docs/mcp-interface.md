# orchestragent-mcp MCP Interface

## Connect
- Server metadata: name `orchestragent-mcp`, version `0.1.0`
- Transport: `stdio`
- Command: `.\bin\orchestragent-mcp.exe -repo <path-to-git-repo> [-db <database-directory>]`
- Defaults: repo = current working directory; db directory = current working directory, database file created as `.orchestragent-mcp.db`
- Repo assumptions: base branch `main`; worktrees live under `.worktrees/`; session branches are `session-<sessionId>`
- Example registration (Codex CLI): `codex mcp add orchestragent-mcp -- ".\bin\orchestragent-mcp.exe" -repo C:\path\to\repo`

## Tools
### `create_worktree`
- Purpose: Create an isolated git worktree and branch for a session.
- Params: `sessionId` (string, required) – 2–50 chars, lowercase letters/numbers/hyphens, must start/end with alphanumeric.
- Success result body:
  - `sessionId` (string)
  - `worktreePath` (string)
  - `branchName` (string, `session-<sessionId>`)
  - `status` (string: `open`)
- Notes: Fails if session already exists or branch already exists.

Example call payload:
```json
{ "name": "create_worktree", "arguments": { "sessionId": "abc-123" } }
```
Example success content text: `Successfully created worktree for session 'abc-123' at '<path>' on branch 'session-abc-123'.`

### `remove_session`
- Purpose: Remove a session’s worktree and branch.
- Params:
  - `sessionId` (string, required)
  - `force` (boolean, optional, default `false`) – skip safety checks.
- Result body:
  - `sessionId` (string)
  - `removedAt` (RFC3339 string, omitted if not deleted)
  - `hasUnmergedChanges` (bool)
  - `unmergedCommits` (int)
  - `uncommittedFiles` (int)
  - `warning` (string, optional)
- Behavior: If `force=false` and there are uncommitted files or unpushed commits, the call returns with `hasUnmergedChanges=true` and a warning; the worktree is **not** removed. Set `force=true` to delete anyway.

Example call:
```json
{ "name": "remove_session", "arguments": { "sessionId": "abc-123", "force": false } }
```
Example warning content: `WARNING: Session 'abc-123' has unmerged changes ... Call with force=true to remove anyway.`

### `get_sessions`
- Purpose: List all tracked sessions with git diff stats vs base branch.
- Params: none.
- Result body:
  - `sessions` (array of):
    - `sessionId` (string)
    - `worktreePath` (string)
    - `branchName` (string)
    - `status` (string: `open` | `reviewed` | `merged`)
    - `linesAdded` (int)
    - `linesRemoved` (int)
- Example content text: `Found 2 session(s)`.

Example call:
```json
{ "name": "get_sessions", "arguments": {} }
```

## Error/response conventions
- Text responses are returned in `content` as plain text; `IsError=true` when a tool fails.
- Common failure reasons: invalid `sessionId` format, session not found, git errors, branch/worktree already exists.
- If `remove_session` finds unmerged work and `force=false`, it returns `IsError=false` but `hasUnmergedChanges=true` to prompt the client to confirm with `force=true`.

## Client usage hints
- Always send lowercased, hyphen-safe `sessionId` values (2–50 chars).
- Before calling `remove_session` with `force=true`, surface `warning` to the user.
- `get_sessions` diff stats fall back to zeros if git diff fails, so treat zeros as “unknown” if an error is likely.
