package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) spellCritMultiplier() float64 {
	return warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.GetRuin(), 1, 0))
}

func (warlock *Warlock) applyTalents() {
	if warlock.Talents.GetDemonicEmbrace() > 0 {
		staminaMultiplier := 1 + 0.03*float64(warlock.Talents.GetDemonicEmbrace())
		spiritMultiplier := 1 - 0.01*float64(warlock.Talents.GetDemonicEmbrace())
		warlock.MultiplyStat(stats.Stamina, staminaMultiplier)
		warlock.MultiplyStat(stats.Spirit, spiritMultiplier)
	}

	if warlock.Talents.GetDemonicTactics() > 0 {
		warlock.AddStat(stats.SpellCritRating, float64(warlock.Talents.GetDemonicTactics())*core.SpellCritRatingPerCritPercent)
	}

	destructionClassMask := WarlockSpellShadowBolt | WarlockSpellIncinerate | WarlockSpellImmolate
	afflictionClassMask := WarlockSpellCorruption | WarlockSpellAgony | WarlockSpellUnstableAffliction | WarlockSpellSiphonLife

	if warlock.Talents.GetCataclysm() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_PowerCost_Pct,
			ClassMask:  destructionClassMask,
			FloatValue: -0.01 * float64(warlock.Talents.GetCataclysm()),
		})
	}

	if warlock.Talents.GetBane() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_CastTime_Flat,
			ClassMask: WarlockSpellShadowBolt | WarlockSpellImmolate,
			TimeValue: -100 * time.Millisecond * time.Duration(warlock.Talents.GetBane()),
		})
	}

	if warlock.Talents.GetImprovedCorruption() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_CastTime_Flat,
			ClassMask: WarlockSpellCorruption,
			TimeValue: -400 * time.Millisecond * time.Duration(warlock.Talents.GetImprovedCorruption()),
		})
	}

	if warlock.Talents.GetDevastation() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  destructionClassMask,
			FloatValue: float64(warlock.Talents.GetDevastation()),
		})
	}

	if warlock.Talents.GetImprovedImmolate() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  WarlockSpellImmolate,
			FloatValue: 0.05 * float64(warlock.Talents.GetImprovedImmolate()),
		})
	}

	if warlock.Talents.GetEmberstorm() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			School:     core.SpellSchoolFire,
			FloatValue: 0.02 * float64(warlock.Talents.GetEmberstorm()),
		})
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_CastTime_Pct,
			ClassMask:  WarlockSpellIncinerate,
			FloatValue: -0.02 * float64(warlock.Talents.GetEmberstorm()),
		})
	}

	if warlock.Talents.GetShadowMastery() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			School:     core.SpellSchoolShadow,
			FloatValue: 0.02 * float64(warlock.Talents.GetShadowMastery()),
		})
	}

	if warlock.Talents.GetContagion() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  afflictionClassMask,
			FloatValue: 0.01 * float64(warlock.Talents.GetContagion()),
		})
	}

	if warlock.Talents.GetSuppression() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusHit_Percent,
			ClassMask:  afflictionClassMask,
			FloatValue: 2 * float64(warlock.Talents.GetSuppression()),
		})
	}

	if warlock.Talents.GetDestructiveReach() > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_ThreatMultiplier_Pct,
			ClassMask:  destructionClassMask,
			FloatValue: -0.05 * float64(warlock.Talents.GetDestructiveReach()),
		})
	}

	if warlock.Talents.GetSoulLink() {
		warlock.PseudoStats.DamageDealtMultiplier *= 1.05
	}
}

func (warlock *Warlock) applyTalentsPostSpellRegistration() {
	warlock.setupNightfall()
}

func (warlock *Warlock) setupNightfall() {
	if warlock.Talents.GetNightfall() == 0 || warlock.ShadowBolt == nil {
		return
	}

	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
		Label:    "Nightfall",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: 10 * time.Second,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != warlock.ShadowBolt || spell.CurCast.CastTime != 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Nightfall Trigger",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if !spell.Matches(WarlockSpellCorruption) {
				return
			}
			if sim.Proc(0.04, "Nightfall") {
				warlock.NightfallProcAura.Activate(sim)
			}
		},
	}))
}

func (warlock *Warlock) setupDemonicSacrifice() {
	if !warlock.Talents.GetDemonicSacrifice() {
		return
	}

	for _, pet := range []*WarlockPet{warlock.Imp, warlock.Voidwalker, warlock.Succubus, warlock.Felhunter, warlock.Felguard} {
		if pet != nil {
			pet.DisableOnStart()
		}
	}
	warlock.ActivePet = nil

	switch warlock.Options.GetSummon() {
	case proto.WarlockOptions_Imp:
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.15
	case proto.WarlockOptions_Succubus:
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.15
	case proto.WarlockOptions_Felhunter, proto.WarlockOptions_Voidwalker, proto.WarlockOptions_Felguard:
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.10
	}
}
