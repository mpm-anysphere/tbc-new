package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (warlock *Warlock) registerImmolate() {
	coeff := 0.2 + 0.04*float64(warlock.Talents.GetShadowAndFlame())

	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 27215},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellImmolate,

		ManaCost: core.ManaCostOptions{FlatCost: 445},
		Cast: core.CastConfig{DefaultCast: core.Cast{
			GCD:      core.GCDDefault,
			CastTime: 2 * time.Second,
		}},

		DamageMultiplier: 1,
		CritMultiplier:   warlock.spellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: coeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 332, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.RelatedDotSpell.Cast(sim, target)
			}
			spell.DealDamage(sim, result)
		},
	})

	warlock.Immolate.RelatedDotSpell = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 27215}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellImmolateDot,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate",
			},
			NumberOfTicks:       5,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: false,
			BonusCoefficient:    0.13,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, 615.0/5.0)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}
