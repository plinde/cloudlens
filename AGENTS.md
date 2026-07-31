# AGENTS.md — plinde/cloudlens

Agent instructions for this repository. `CLAUDE.md` is a symlink to this file.

`cloudlens` is a terminal UI for AWS/GCP ("k9s for cloud"). This checkout is
**`plinde/cloudlens`, a personal fork** of `one2nc/cloudlens`.

## 🚫 Work on the fork only — never touch upstream

**`origin` is `plinde/cloudlens`. `upstream` is `one2nc/cloudlens`, and it is read-only.**

Allowed against `upstream`: `git fetch upstream`, and reading its issues/PRs/code.

**Never**, unless the user explicitly asks for that specific action in that
message:

- `git push upstream ...` — anything at all
- `gh pr create --repo one2nc/cloudlens` (or with `--repo` omitted while `upstream`
  is the default, which is the easy way to do this by accident)
- `gh issue create` / `gh issue comment` / any comment or review on an upstream
  issue or PR
- editing or closing anything in the upstream repo

Every PR belongs to `plinde/cloudlens` — always pass `--repo plinde/cloudlens` and
`--base main` explicitly. `gh` resolves a bare `gh pr create` against the
*upstream* of a fork, so the explicit flags are what keep an accidental
upstream PR from happening. Merging a PR on this fork is fine when asked.

If something genuinely looks worth upstreaming, **say so and stop** — the user
decides whether anything is ever sent to `one2nc/cloudlens`.

## All work must be done in linked worktrees

**No work happens in the main checkout.** All changes, commits, builds, and
validation occur in a linked worktree.

```
~/workspace/github.com/plinde/cloudlens/                 ← main checkout, keep clean
~/workspace/github.com/plinde/cloudlens--<description>/  ← linked worktrees
```

The main checkout is for `git fetch`, `git merge`/fast-forward, and creating
worktrees — nothing else. Never edit files, run builds, or commit in it.

Branch off `origin/main` after a fetch, and pull the main checkout forward
(`git pull --ff-only`) once a PR merges.

## Keep fork-local changes isolated in their own commits

Fork-local preferences must not be mixed into commits that would otherwise be
clean upstream changes. Anything that is *only* a local preference gets its own
commit, so the rest stays cherry-pickable if the user ever chooses to offer it.

## Build, test, install — validate before PR

Before marking a feature ready, run the full pipeline from the **worktree**:

```bash
go build ./... && go vet ./... && gofmt -l . && go test ./...
make install
```

`make install` builds to `execs/cloudlens` and copies the binary to
`~/.local/bin/cloudlens`, including macOS re-codesign. This lets the user run
`cloudlens` to validate the feature works before the PR is merged.

## Functionality Review Handoff

When asking the user to review functionality, always run `make install` from the
linked worktree immediately before the handoff. Report the installed binary's
version/commit and the worktree used. Never ask the user to review against an
uninstalled or stale `~/.local/bin/cloudlens` binary.

## Architecture

cloudlens follows a layered architecture for each AWS service view:

| Layer | Path | Purpose |
|-------|------|---------|
| AWS SDK calls | `internal/aws/<svc>.go` | Raw API calls returning `[]RespType` |
| Response types | `internal/aws/types.go` | Plain structs of displayable fields |
| DAO | `internal/dao/<svc>.go` | Bridges SDK calls to model layer, embeds `Accessor` |
| Renderer | `internal/render/<svc>.go` | `Header() + Render()` column definitions |
| View | `internal/view/<svc>.go` | Key bindings, enter behavior, embeds `Browser`→`Table` |
| Registry | `internal/model/registry.go` | Maps resource string → `{DAO, Renderer}` |
| View Registrar | `internal/view/registrar.go` | Maps resource string → `MetaViewer{viewerFn}` |
| Aliases | `internal/config/alias.go` | Default shortcuts (`a.declare(...)`) |
| Constants | `internal/constants.go` | Resource string constants + context keys |

## Adding a new service view checklist

Files to touch (8-9 locations):

1. `internal/constants.go` — Add `LowercaseXxx`/`UppercaseXxx` + context keys
2. `internal/aws/types.go` — Add response struct
3. `internal/aws/<svc>.go` — **Create** — SDK API call function
4. `internal/dao/<svc>.go` — **Create** — `List()`, optional `Describe()`/mutation methods
5. `internal/render/<svc>.go` — **Create** — `Header()` + `Render()`
6. `internal/model/registry.go` — Add entry
7. `internal/view/<svc>.go` — **Create** — View struct with `bindKeys()` + `enterCmd()`
8. `internal/view/registrar.go` — Add `MetaViewer` entry
9. `internal/config/alias.go` — Add `a.declare(...)` in `loadDefaultAliases`
10. `go.mod` — Add SDK dependency if new service

## Local context

The user's `/cloudlens` skill (`~/.agents/skills/cloudlens/SKILL.md`) documents the
per-AWS-account `cloudlens-*` shell wrappers built on this binary.
