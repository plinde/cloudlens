# cloudlens (plinde fork) — Feature Spec

Fork of `one2nc/cloudlens` (abandoned since Apr 2024). This fork is the de facto
continuation. All work happens in the `plinde/cloudlens` fork only; upstream is
read-only.

## Status

- **PR**: [#1](https://github.com/plinde/cloudlens/pull/1) — open, not yet merged
- **Branch**: `feature/asg-lifecycle-view`
- **Worktree**: `~/workspace/github.com/plinde/cloudlens--asg-lifecycle`

## Features Added

### ASG Lifecycle View (`asg`)
- **List**: 10-column table (Name, Min, Max, Desired, Status, AZs, Health Check, Cooldown, Launch Config, Termination Policies)
- **Edit** (`e`): Modal form to change min/max/desired capacity with validation (Min ≤ Desired ≤ Max). Calls `UpdateAutoScalingGroup` API.
- **Describe** (`d`): Formatted key-value detail view (replaces raw JSON)
- **Drill-down** (`Enter`): Shows ASG instances in a sub-view with Instance-Id, Type, AZ, Lifecycle State, Health Status, Protected columns. Title bar shows `ASG <name> (min=X max=Y desired=Z)` — persists across refreshes.
- **Command-line**: `cloudlens aws --view asg` starts directly in the ASG view
- **Shell wrappers**: `cloudlens-kalibrate-staging` passes `--view asg` by default

### Background Refresh & Footer
- **Periodic refresh**: Active view polls AWS every 30s (was disabled upstream)
- **Manual refresh** (`r`): Immediate API re-fetch, resets the 30s timer
- **Footer bar**: Replaces breadcrumbs — shows context-sensitive hotkeys + last-sync timestamp. Updates on every refresh.

### Navigation Hotkeys
- `q` / `esc`: Back (pop view stack). At root level, quits the app.
- `ctrl+c`: Hard quit from anywhere.
- `d`: Describe (k9s-style detail) — formatted key-value for ASG, live text for other views.

### Project Quality
- **AGENTS.md** + **CLAUDE.md**: Fork-only instructions, architecture reference, worktree discipline
- **Makefile**: `help` (default), `build`, `install`, `clean`, `test`, `test-v`, `cover`, `run`
- **Shell aliases**: `cl` → `cloudlens`, `cl-kalibrate-staging` → wrapper, `cl-kalibrate-prod`, etc. via `~/.zsh/claude-managed.zshrc`
- **Zsh completions**: `~/.zsh/completion/_cloudlens` (cobra-generated) + `_cl` wrapper completion
- **AWS SDK compatibility**: EC2 SDK aligned with the core SDK; middleware-stack regressions covered by a transport-level test
- **EC2 snapshot command**: canonical `ec2:s`; legacy `ec2:S` remains accepted

## Architecture

Follows the existing cloudlens layered pattern:

| Layer | Key Files |
|-------|-----------|
| **AWS SDK** | `internal/aws/asg.go` — `GetASGs`, `GetASGInstances`, `GetASGFormatted`, `UpdateASGSize`, `UpdateASGDesiredCapacity` |
| **Types** | `internal/aws/types.go` — `ASGResp` (16 fields), `ASGInstanceResp` (6 fields) |
| **DAO** | `internal/dao/asg.go`, `asg_instance.go` — `List`, `Describe`, `UpdateSize` |
| **Renderer** | `internal/render/asg.go`, `asg_instance.go` — column headers + row rendering |
| **View** | `internal/view/asg.go`, `asg_instance.go` — key bindings, edit form, drill-down |
| **Model** | `internal/model/table.go` — enabled `updater()` goroutine, added `ResetTimer()` |
| **UI** | `internal/ui/footer.go` — hotkey bar + last-sync; `internal/ui/table.go` — `SetCustomTitle()` |
| **Config** | `internal/config/cloud_config.go` — `StartView` field for `--view` flag |
| **CLI** | `cmd/aws.go` — `--view` flag |

## Shell Integration

- `cl` alias → `cloudlens` (in `~/.zsh/claude-managed.zshrc`)
- `cl-*` aliases → `cloudlens-*` per-account wrappers (kalibrate, topiq, management, shared-services)
- `~/.zsh/completion/_cloudlens` + `_cl` — tab completion for all subcommands and wrappers

## Known Limitations

- V1 SDK (`aws-sdk-go/v1`) still used alongside V2 — causes double credential-load
- No SSO/OIDC support — reads static keys from `~/.aws/credentials` only
- No `AWS_SESSION_TOKEN` in env-var auth path
- 48 pre-existing `go vet` issues (fmt.Errorf format strings)
- 10 pre-existing `gofmt` violations
- Test suite partially broken (4 packages fail to build, 3 assertion mismatches from upstream "Name column" commit)
- ASG edit form appears as modal dialog; no undo on failure (but API validates first)
- ASG instance drill-down fetches full ASG for instances only — could skip the extra call by caching

## Resource Command Convention

- `ec2` lists EC2 instances.
- `ec2:s` lists EC2 snapshots (`ec2:S` remains a compatibility alias).
- `asg` lists Auto Scaling groups.
- Selecting an ASG and pressing Enter opens its instance-membership view. `asg:i` is the internal resource identifier for that nested view, not a replacement for `asg`.

## Upcoming

- Add `d` (Describe) to all views, not just S3/ECS/ASG
- Fix pre-existing go vet + test issues
- Add SSO credential support
- Add `--view` flag to GCP subcommand
