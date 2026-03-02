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

## Execution update snapshot (this branch)

1. Active Warlock path now includes TBC core filler/DoT spells (Shadow Bolt, Immolate, Incinerate, Curse of Agony, Unstable Affliction, Siphon Life) and no longer relies on Fel Flame as the default filler.
2. `ApplyTalents()` now applies TBC-relevant Affliction/Demonology/Destruction effects (including Nightfall and Ruin-critical multipliers) in the active path.
3. Active pet registration/summon handling is re-enabled for Imp/Voidwalker/Succubus/Felhunter/Felguard with executable autocast rotations.
4. Demonic Sacrifice support is active in the core path and validated with summon-selection tests.
5. `sim/warlock/items.go` now contains active TBC set/proc wiring (Voidheart/Corruptor/Malefic/Felshroud + Ashtongue proc).

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
| Register Affliction/Destruction/Demonology agents in global path | PARTIAL | TBC class guide + repo architecture | N/A (single-spec proto model) | Repo currently exposes only `SpecWarlock`; spec identity is talent/APL-driven under one active agent. |
| Curse of Elements (1490) | DONE-ACTIVE | TBC guide + `sim/warlock/curse_of_elements.go` | Indirect via DPS tests | Active and usable in current APL. |
| Corruption (172) | DONE-ACTIVE | TBC guide + `sim/warlock/corruption.go` | DPS + cast-profile tests | Active with TBC tick cadence and empowered coefficient scaling. |
| Life Tap (1454) mana economy | DONE-ACTIVE | TBC guides + `sim/warlock/lifetap.go` | DPS + cast-profile tests | Active with improved life tap scaling, spellpower contribution, and mana-feed transfer. |
| Shadow Bolt filler | DONE-ACTIVE | TBC guide + `sim/warlock/shadowbolt.go` | `TestWarlockRotationCastsShadowBoltFiller` | Active default filler in DS/Ruin profile. |
| Immolate/Incinerate/Conflagrate (Destruction core) | PARTIAL | TBC guides + active `sim/warlock/immolate.go`, `sim/warlock/incinerate.go` | DPS + cast-profile tests | Immolate/Incinerate active; Conflagrate still pending. |
| Unstable Affliction, Curse of Agony, Siphon Life (Affliction core) | DONE-ACTIVE | TBC guides + active spell files | `TestWarlockAfflictionCastProfile` | Active and talent/APL gated in current compile path. |
| Seed of Corruption (AoE identity) | MISSING (active) / PARKED | TBC guides + parked `_affliction/seed_of_corruption.go` | None | Needed for multi-target fidelity. |
| Fel Flame filler in default APL | DONE-ACTIVE | `ui/warlock/dps/apls/default.apl.json` | `TestWarlockRotationCastsShadowBoltFiller` | Removed from default TBC profile. |
| Talent parsing from talents string | DONE-ACTIVE | `proto/warlock.proto`, `sim/warlock/warlock.go` | Multiple warlock tests | Trees parse and active talent effects now drive spells/auras/modifiers. |
| TBC Affliction tree effects | PARTIAL | `proto/warlock.proto` + TBC references | Cast-profile + DPS tests | Core damage/economy talents active; full utility/tuning pass still pending. |
| TBC Demonology tree effects (including DS/sac utility) | PARTIAL | `proto/warlock.proto` + TBC references | Pet/sac tests + DPS tests | DS/sac and core pet hooks active; full demonology-depth modeling pending. |
| TBC Destruction tree effects (Ruin/Backlash/etc) | PARTIAL | `proto/warlock.proto` + TBC references | Cast-profile + DPS tests | Core DS/Ruin multipliers and cast-time/cost hooks active; remaining talents pending. |
| Active pet registration / summon selection from options | DONE-ACTIVE | `proto/warlock.proto`, `sim/warlock/pet.go`, UI pet input | `TestWarlockPetSummonSelection` | Active summon options now map to registered pets/autocast. |
| Pet AI/autocast priorities | DONE-ACTIVE | `sim/warlock/pet.go`, `sim/warlock/pet_abilities.go` | `TestWarlockPetSummonSelection` | Active priority-based autocast loop in compile path. |
| Pet stat scaling from owner/talents/buffs | PARTIAL | `sim/warlock/pet.go` + TBC guides | Pet summon tests + DPS tests | Active inheritance restored; deep validation/calibration still pending. |
| DS / Demonic Sacrifice build support | DONE-ACTIVE | `sim/warlock/talents.go` | `TestWarlockPetSummonSelection` | DS aura handling and pet suppression are active in current path. |
| Warlock-specific TBC set bonuses | PARTIAL | TBC item sets + `sim/warlock/items.go` | Indirect via warlock suite | Active T4/T6 + Ashtongue hooks implemented; full T5 behavior still pending. |
| Warlock trinket/proc timing interactions | PARTIAL | TBC trinkets + core proc patterns | Indirect via warlock suite | Ashtongue proc path is active; broader proc timing coverage still pending. |
| Warlock APL value support for TBC rotation logic | MISSING | APL architecture + TBC rotation needs | None | Active TBC-specific custom APL values are absent. |
| Preset APLs for DS/Ruin, Affliction UA, Affliction/Ruin | DONE-ACTIVE | `ui/warlock/dps/presets.ts`, `ui/warlock/dps/apls/*.json` | Cast-profile + DPS tests | Presets/APLs now route to TBC filler/DoT priorities. |
| Cast distribution / DoT uptime assertions | PARTIAL | Regression harness goal | `TestWarlockRotationCastsShadowBoltFiller`, `TestWarlockAfflictionCastProfile` | Cast-distribution assertions added; explicit uptime assertions still pending. |
| Multi-target and pet scenario regression tests | PARTIAL | Regression harness goal | `TestWarlockPetSummonSelection` | Pet regression coverage added; dedicated multi-target suite still pending. |

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
