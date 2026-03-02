package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (warlock *Warlock) registerShadowBolt() {
	coeff := 0.857 + 0.04*float64(warlock.Talents.GetShadowAndFlame())

	if warlock.Talents.GetImprovedShadowBolt() > 0 {
		points := warlock.Talents.GetImprovedShadowBolt()
		warlock.ImpShadowboltAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			return core.ImprovedShadowBoltAura(target, 0, points)
		})
	}

	warlock.ShadowBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 27209},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellShadowBolt,
		MissileSpeed:   20,

		ManaCost: core.ManaCostOptions{FlatCost: 420},
		Cast: core.CastConfig{DefaultCast: core.Cast{
			GCD:      core.GCDDefault,
			CastTime: 3 * time.Second,
		}, ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
			if warlock.NightfallProcAura != nil && warlock.NightfallProcAura.IsActive() {
				cast.CastTime = 0
			}
		}},

		DamageMultiplier: 1,
		CritMultiplier:   warlock.spellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: coeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warlock.CalcAndRollDamageRange(sim, 544, 607)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if warlock.ImpShadowboltAuras == nil || !result.Landed() || !result.DidCrit() {
					return
				}

				debuff := warlock.ImpShadowboltAuras.Get(result.Target)
				debuff.Activate(sim)
				debuff.SetStacks(sim, 4)
			})
		},

		RelatedAuraArrays: warlock.ImpShadowboltAuras.ToMap(),
	})
}
