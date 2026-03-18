package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (paladin *Paladin) registerCrusaderStrikeSpell() {
	if !paladin.Talents.CrusaderStrike {
		return
	}

	paladin.CrusaderStrike = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 35395},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		MaxRange: core.MaxMeleeRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: 236,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1.1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}
