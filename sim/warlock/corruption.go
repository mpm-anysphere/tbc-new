package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (warlock *Warlock) RegisterCorruption() *core.Spell {
	empoweredCorruptionBonus := 0.02 * float64(warlock.Talents.GetEmpoweredCorruption())
	bonusCoefficient := 0.156 + empoweredCorruptionBonus

	warlock.Corruption = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 172},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCorruption,

		ManaCost: core.ManaCostOptions{FlatCost: 370},
		Cast: core.CastConfig{DefaultCast: core.Cast{
			GCD:      core.GCDDefault,
			CastTime: 2 * time.Second,
		}},

		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Corruption",
			},
			NumberOfTicks:       6,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: false,
			BonusCoefficient:    bonusCoefficient,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, 150)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			dot := spell.Dot(target)
			if useSnapshot {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
				result.Damage /= dot.TickPeriod().Seconds()
				return result
			} else {
				result := spell.CalcPeriodicDamage(sim, target, 150, spell.OutcomeExpectedMagicAlwaysHit)
				result.Damage /= dot.TickPeriod().Seconds()
				return result
			}
		},
	})

	return warlock.Corruption
}
