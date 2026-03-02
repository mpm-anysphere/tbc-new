package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ItemSetGladiatorsFelshroud = core.NewItemSet(core.ItemSet{
	ID:   615,
	Name: "Gladiator's Felshroud",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			agent.GetCharacter().AddStat(stats.SpellDamage, 29)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			agent.GetCharacter().AddStat(stats.SpellDamage, 88)
		},
	},
})

var ItemSetVoidheartRaiment = core.NewItemSet(core.ItemSet{
	ID:   645,
	Name: "Voidheart Raiment",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()
			shadowflameAura := warlock.NewTemporaryStatsAura(
				"Voidheart Shadowflame",
				core.ActionID{SpellID: 37377},
				stats.Stats{stats.SpellDamage: 135},
				time.Second*15,
			)

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Voidheart 2pc Proc",
				Callback:   core.CallbackOnCastComplete,
				ProcChance: 0.05,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if !spell.SpellSchool.Matches(core.SpellSchoolShadow | core.SpellSchoolFire) {
						return
					}
					shadowflameAura.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_DotNumberOfTicks_Flat,
				ClassMask: WarlockSpellImmolateDot,
				IntValue:  1,
			})
		},
	},
})

var ItemSetCorruptorRaiment = core.NewItemSet(core.ItemSet{
	ID:   646,
	Name: "Corruptor Raiment",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, _ *core.Aura) {},
		4: func(agent core.Agent, _ *core.Aura) {},
	},
})

var ItemSetMaleficRaiment = core.NewItemSet(core.ItemSet{
	ID:   670,
	Name: "Malefic Raiment",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, _ *core.Aura) {},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  WarlockSpellShadowBolt | WarlockSpellIncinerate,
				FloatValue: 0.06,
			})
		},
	},
})

func init() {
	core.NewItemEffect(32493, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()
		procAura := warlock.NewTemporaryStatsAura(
			"Ashtongue Talisman of Shadows",
			core.ActionID{SpellID: 40478},
			stats.Stats{stats.SpellDamage: 220},
			time.Second*5,
		)

		warlock.MakeProcTriggerAura(core.ProcTrigger{
			Name:           "Ashtongue Talisman Proc",
			Callback:       core.CallbackOnPeriodicDamageDealt,
			ClassSpellMask: WarlockSpellCorruption,
			ProcChance:     0.2,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})
}
