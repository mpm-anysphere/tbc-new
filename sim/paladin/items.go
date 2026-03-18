package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func init() {
	// Libram of Avengement: on Judgement of Blood cast, gain +53 melee and spell
	// crit rating for 5s. Legacy TBC behavior.
	core.NewItemEffect(27484, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura(
			"Libram of Avengement Proc",
			core.ActionID{SpellID: 34260},
			stats.Stats{
				stats.MeleeCritRating: 53,
				stats.SpellCritRating: 53,
			},
			time.Second*5,
		)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of Avengement",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if spell == paladin.JudgementOfBlood {
					procAura.Activate(sim)
				}
			},
		})
	})
}

