package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (warlock *Warlock) registerCurseOfAgony() {
	agonyBaseTick := (1356.0 / 12.0) * (1 + 0.02*float64(warlock.Talents.GetImprovedCurseOfAgony()))

	warlock.CurseOfAgony = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 27218},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellAgony,

		ManaCost: core.ManaCostOptions{FlatCost: 265},
		Cast: core.CastConfig{DefaultCast: core.Cast{
			GCD: core.GCDDefault,
		}},

		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Curse of Agony",
			},
			NumberOfTicks:       12,
			TickLength:          2 * time.Second,
			AffectedByCastSpeed: false,
			BonusCoefficient:    0.1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, agonyBaseTick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
