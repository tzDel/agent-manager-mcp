# orchestrAIgent

An MCP (Model Context Protocol) server that manages isolated git worktrees for AI coding agents.

## What it does
- Creates `.worktrees/session-<id>` directories plus `session-<id>` branches on demand (base branch is `main`).
- Persists session data in a SQLite database (`.orchestrAIgent.db` by default). This ensures that session information survives server restarts.
- Lists active sessions with worktree paths and diff stats versus the base branch.
- Removes sessions with safety checks for uncommitted files and unpushed commits (override with `force=true`).
- Speaks MCP over stdio so you can register it with any MCP-compatible client.

## Architecture

See [docs/architecture.md](docs/architecture.md) for detailed architecture documentation.

## MCP Tools

- **Create Worktree**: Create isolated git worktree and branch for an agent
  ```
  Example: "Create a worktree for session abc123"
  ```

- **Remove Session**: Remove agent worktree and branch
  ```
  Example: "Remove all sessions" || "Remove session abc123"
  ```

- **List Sessions**: Retrieve status, worktree path, and branch information
  ```
  Example: "List all sessions" || "List session abc123"
  ```

## Configure with MCP Clients

### Claude Code CLI
```shell
claude mcp add --transport stdio orchestrAIgent -- "<path-to-binary>" -repo <path-to-repo>
claude mcp list
```

### Codex CLI
```shell
codex mcp add orchestrAIgent -- "<path-to-binary>" -repo <path-to-repo>
codex mcp list
```

## Project Status

Currently in development. See [docs/concept.md](docs/concept.md) for vision and roadmap.
