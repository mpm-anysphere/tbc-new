package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (paladin *Paladin) RegisterAvengingWrathCD() {
	actionID := core.ActionID{SpellID: 31884}

	paladin.AvengingWrathAura = paladin.RegisterAura(core.Aura{
		Label:    "Avenging Wrath" + paladin.Label,
		ActionID: actionID,
		Duration: time.Second * 20,
	}).AttachMultiplicativePseudoStatBuff(&paladin.PseudoStats.DamageDealtMultiplier, 1.3)

	spell := paladin.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		ManaCost: core.ManaCostOptions{
			FlatCost: 236,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{NonEmpty: true},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.AvengingWrathAura.Activate(sim)
		},
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

