# Warlock TBC Fidelity - Phase 1 Baseline and Gap Audit

## Scope

This document defines what is currently implemented vs missing for TBC Warlock fidelity, and sets a phased implementation order aligned with:

- Spells
- Talents
- Pets
- Set bonuses / item effects
- APL value support
- Test coverage and regression safety

## Baseline findings (current branch)

1. Only the base Warlock package is registered in global sim bootstrapping (`RegisterAll()` calls `warlock.RegisterWarlock()` only).  
2. Active base Warlock behavior currently registers a minimal subset (`Curse of Elements`, `Corruption`, `Life Tap`, `Fel Flame`, `Drain Life`, `Hellfire`) with a minimal default APL.  
3. Pet fields and registration calls are commented out in the active Warlock type/construction path.  
4. Most richer Warlock implementation is parked in underscore-prefixed files/directories (for example `_talents.go`, `_pets.go`, `_affliction`, `_destruction`, `_demonology`), which are not part of normal `go list ./...` package discovery.  
5. Warlock set bonus file exists, but all item set logic is commented and targets non-TBC tier sets.  
6. Existing Warlock tests are baseline DPS/stat checks; spec-level tests are placeholders.

## External references used for expected TBC behavior

- Wowhead TBC Warlock talent calculator (canonical tree/point layout): `https://www.wowhead.com/tbc/talent-calc/warlock`
- Wowhead TBC Warlock class guide (spell/talent/rotation expectations): `https://www.wowhead.com/tbc/guide/classes/warlock`
- Icy Veins TBC Warlock DPS guide (cross-check for spec play patterns and pet usage): `https://www.icy-veins.com/tbc-classic/warlock-dps-pve-guide`

## Existing architecture patterns to follow

- Talent wiring pattern: class-level `ApplyTalents()` dispatch to `register*` handlers (for example `sim/rogue/rogue.go` + `sim/rogue/talents_*.go`).
- Set bonus wiring pattern: `core.NewItemSet` + `AttachSpellMod` / `AttachProcTrigger` (for example `sim/rogue/items.go`).
- APL value wiring pattern: `NewAPLValue` in sim + UI registration in `ui/core/components/individual_sim_ui/apl_values.ts`.

## Tracking matrix (feature / status / source / test coverage)

Status legend:

- `DONE-ACTIVE`: implemented and reachable in active sim path
- `PARKED`: implementation exists but is not on active compile/runtime path
- `MISSING`: no implementation in active or parked path
- `PARTIAL`: implemented but fidelity-incomplete for TBC expectations

| Feature | Status | Source | Test coverage | Notes |
|---|---|---|---|---|
| Register base Warlock agent | DONE-ACTIVE | `sim/register_all.go`, `sim/warlock/warlock.go` | `sim/warlock/warlock_test.go` | Active spec entrypoint is only base Warlock. |
| Register Affliction/Destruction/Demonology agents in global path | MISSING | TBC class guide + repo architecture | None | Spec packages exist but are not globally registered. |
| Curse of Elements (1490) | DONE-ACTIVE | TBC guide + `sim/warlock/curse_of_elements.go` | Indirect via DPS tests | Active and usable in current APL. |
| Corruption (172) | DONE-ACTIVE | TBC guide + `sim/warlock/corruption.go` | Indirect via DPS tests | Snapshot logic exists, coefficients need TBC verification pass. |
| Life Tap (1454) mana economy | PARTIAL | TBC guides + `sim/warlock/lifetap.go` | Indirect via DPS tests | No dedicated sustain/economy tests yet. |
| Shadow Bolt filler | MISSING (active) / PARKED | TBC guide + parked spec code | None | Core TBC filler absent in active base path. |
| Immolate/Incinerate/Conflagrate (Destruction core) | MISSING (active) / PARKED | TBC guides + parked `_destruction` code | None | Not reachable in active package flow. |
| Unstable Affliction, Curse of Agony, Siphon Life (Affliction core) | MISSING (active) / PARKED | TBC guides + parked `_affliction` code | None | Not reachable in active package flow. |
| Seed of Corruption (AoE identity) | MISSING (active) / PARKED | TBC guides + parked `_affliction/seed_of_corruption.go` | None | Needed for multi-target fidelity. |
| Fel Flame filler in default APL | PARTIAL / out-of-era | Current `ui/warlock/dps/apls/default.apl.json` | Indirect via DPS tests | Not a TBC baseline rotational filler. |
| Talent parsing from talents string | DONE-ACTIVE | `proto/warlock.proto`, `sim/warlock/warlock.go` | Indirect | Trees parse, effects mostly not applied. |
| TBC Affliction tree effects | MISSING | `proto/warlock.proto` + TBC references | None | No active `register*` talent effects for TBC tree. |
| TBC Demonology tree effects (including DS/sac utility) | MISSING | `proto/warlock.proto` + TBC references | None | No active implementation in base package path. |
| TBC Destruction tree effects (Ruin/Backlash/etc) | MISSING | `proto/warlock.proto` + TBC references | None | No active implementation in base package path. |
| Active pet registration / summon selection from options | MISSING (active) / PARKED | `proto/warlock.proto`, UI pet input | None | Pet fields/calls commented out in active `warlock.go`. |
| Pet AI/autocast priorities | PARKED | `_pets.go` | None | Exists in parked code only. |
| Pet stat scaling from owner/talents/buffs | PARKED | `_pets.go` + TBC guides | None | Needs TBC validation and activation. |
| DS / Demonic Sacrifice build support | MISSING (active) | TBC guides + talent list | None | Critical for DS/Ruin fidelity. |
| Warlock-specific TBC set bonuses | MISSING | TBC item sets + sim itemset pattern | None | `sim/warlock/items.go` is commented and non-TBC oriented. |
| Warlock trinket/proc timing interactions | PARTIAL | TBC trinkets + core proc patterns | No targeted warlock proc tests | Needs deterministic trigger coverage on warlock spell events. |
| Warlock APL value support for TBC rotation logic | MISSING | APL architecture + TBC rotation needs | None | Active TBC-specific custom APL values are absent. |
| Preset APLs for DS/Ruin, Affliction UA, Affliction/Ruin | PARTIAL | `ui/warlock/dps/presets.ts`, `ui/warlock/dps/apls/*.json` | Indirect via DPS tests | Presets exist, but behavior is currently baseline/alpha-level. |
| Cast distribution / DoT uptime assertions | MISSING | Regression harness goal | None | Needed for calibration confidence. |
| Multi-target and pet scenario regression tests | MISSING | Regression harness goal | None | Current suite is mostly single-target baseline DPS. |

## Phase-to-PR implementation order (signed-off proposal)

### Gate 0 (required before Phase 2)

1. Decide structural strategy for parked Warlock code:
   - Option A: migrate needed code out of underscore-prefixed files/directories into active compile paths.
   - Option B: re-implement TBC-only subset directly in active package and treat parked code as reference only.
2. Freeze naming and package boundaries so tests and APLs do not churn repeatedly.

### PR1 - Gap audit + foundational tests (Phase 1 deliverable)

- Commit this matrix.
- Add focused failing/placeholder tests for:
  - core spell expected outputs,
  - talent hooks,
  - pet summon selection,
  - set/proc trigger sanity.
- Exit criteria: agreed red/green map for all rows above.

### PR2 - Core spell/talent fidelity (Phase 2)

- Implement active-path TBC core spells and coefficients by spec identity.
- Implement TBC talent effects in `ApplyTalents()` plus spell hooks.
- Add targeted tests for snapshot/refresh, crit/multipliers, hit/resist, mana sustain.

### PR3 - Pet system completion (Phase 3)

- Re-enable active pet registration and summon handling.
- Implement pet AI priorities and stat scaling.
- Validate DS/sac and pet-dependent builds in tests.

### PR4 - Item/set bonus + proc integration (Phase 4)

- Implement TBC Warlock set bonuses (not non-TBC placeholders).
- Add deterministic tests for proc triggers and uptime windows.

### PR5 - Spec-quality APL profiles + calibration (Phase 5 + first half of Phase 6)

- Build production APLs for:
  - Destruction DS/Ruin
  - Affliction UA
  - Affliction/Ruin
- Add cast-count and uptime checks and baseline DPS calibration outputs.

### PR6 - Final defaults + regression hardening (end of Phase 6)

- Lock default presets after calibration.
- Expand regression suite (single-target, multi-target, pet, mana pressure).
- Document known caveats and expected metric envelopes.

## Definition of done for Phase 1

- Gap list exists and is traceable to source files and TBC references.
- Each feature row has status, source, and test coverage column.
- Implementation order is sequenced and gated.
- Ready for sign-off and execution in PR1-PR6 order.

## Sign-off

- Engineering audit prepared: ✅
- Execution order proposed: ✅
- Awaiting product/maintainer sign-off on matrix scope and Gate 0 strategy: ☐
