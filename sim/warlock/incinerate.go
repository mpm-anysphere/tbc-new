package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (warlock *Warlock) registerIncinerate() {
	coeff := 0.714 + 0.04*float64(warlock.Talents.GetShadowAndFlame())

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 32231},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellIncinerate,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{FlatCost: 355},
		Cast: core.CastConfig{DefaultCast: core.Cast{
			GCD:      core.GCDDefault,
			CastTime: 2500 * time.Millisecond,
		}},

		DamageMultiplier: 1,
		CritMultiplier:   warlock.spellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: coeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warlock.CalcAndRollDamageRange(sim, 444, 514)
			if warlock.Immolate != nil && warlock.Immolate.RelatedDotSpell.Dot(target).IsActive() {
				baseDamage += 119.5
			}

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
