package warlock

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) registerLifeTap() {
	actionID := core.ActionID{SpellID: 1454}
	baseRestore := 582.0 * (1 + 0.1*float64(warlock.Talents.GetImprovedLifeTap()))
	petRestoreFraction := 0.3333 * float64(warlock.Talents.GetManaFeed())

	manaMetrics := warlock.NewManaMetrics(actionID)
	petManaMetrics := make(map[*core.Pet]*core.ResourceMetrics, len(warlock.Pets))
	if warlock.Talents.GetManaFeed() > 0 {
		for _, pet := range warlock.Pets {
			petManaMetrics[pet] = pet.NewManaMetrics(actionID)
		}
	}

	warlock.LifeTap = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellLifeTap,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			restore := baseRestore + (warlock.GetStat(stats.SpellDamage)+warlock.GetStat(stats.ShadowDamage))*0.8

			healthLoss := min(restore, max(warlock.CurrentHealth()-1, 0))
			warlock.RemoveHealth(sim, healthLoss)
			warlock.AddMana(sim, restore, manaMetrics)

			if warlock.Talents.GetManaFeed() == 0 {
				return
			}

			for _, pet := range warlock.Pets {
				if !pet.IsActive() {
					continue
				}
				pet.AddMana(sim, restore*petRestoreFraction, petManaMetrics[pet])
			}
		},
	})
}
