package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func (paladin *Paladin) registerExorcismSpell() {
	paladin.Exorcism = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 10314},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,

		MaxRange: 30,

		ManaCost: core.ManaCostOptions{
			FlatCost: 295,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.SpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.CalcAndRollDamageRange(sim, 521, 0.11) + spell.SpellDamage(target)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}

func (paladin *Paladin) CanExorcism(target *core.Unit) bool {
	return target.MobType == proto.MobType_MobTypeUndead || target.MobType == proto.MobType_MobTypeDemon
}
