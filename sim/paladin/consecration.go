package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (paladin *Paladin) RegisterConsecrationSpell(rank int32) {
	var actionID int32
	var manaCost int32
	var baseTickDamage float64

	switch rank {
	case 6:
		actionID = 27173
		manaCost = 660
		baseTickDamage = 64
	case 4:
		actionID = 20923
		manaCost = 390
		baseTickDamage = 35
	case 1:
		actionID = 26573
		manaCost = 120
		baseTickDamage = 8
	default:
		panic("Unsupported Consecration rank")
	}

	paladin.Consecration = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: actionID},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,

		MaxRange: 8,
		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				ActionID: core.ActionID{SpellID: actionID},
				Label:    "Consecration" + paladin.Label,
			},
			NumberOfTicks: 8,
			TickLength:    time.Second,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Consecration uses dynamic snapshotless ticking.
				tickDamage := baseTickDamage + 0.12*dot.Spell.SpellDamage(target)
				dot.Spell.CalcPeriodicAoeDamage(sim, tickDamage, dot.Spell.OutcomeMagicHit)
				dot.Spell.DealBatchedPeriodicDamage(sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}

