package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (paladin *Paladin) ApplyTalents() {
	// Core TBC Ret talent-driven stats.
	paladin.AddStat(stats.PhysicalCritPercent, float64(paladin.Talents.Conviction))
	paladin.AddStat(stats.PhysicalCritPercent, float64(paladin.Talents.SanctifiedSeals))
	paladin.AddStat(stats.SpellCritPercent, float64(paladin.Talents.SanctifiedSeals))

	paladin.AddStat(stats.PhysicalHitPercent, float64(paladin.Talents.Precision))
	paladin.AddStat(stats.SpellHitPercent, float64(paladin.Talents.Precision))

	paladin.AddStat(stats.ParryRating, float64(paladin.Talents.Deflection)*core.ParryRatingPerParryPercent)

	if paladin.Talents.DivineStrength > 0 {
		bonus := 1 + 0.02*float64(paladin.Talents.DivineStrength)
		paladin.MultiplyStat(stats.Strength, bonus)
	}
	if paladin.Talents.DivineIntellect > 0 {
		bonus := 1 + 0.02*float64(paladin.Talents.DivineIntellect)
		paladin.MultiplyStat(stats.Intellect, bonus)
	}

	paladin.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= paladin.WeaponSpecializationMultiplier()
	paladin.applyCrusade()
	paladin.applyVengeance()
}

func (paladin *Paladin) applyCrusade() {
	multiplier := paladin.crusadeMultiplier()
	if multiplier != 1 {
		paladin.PseudoStats.DamageDealtMultiplier *= multiplier
	}
}

func (paladin *Paladin) applyVengeance() {
	if paladin.Talents.Vengeance == 0 {
		return
	}

	bonusPerStack := 0.01 * float64(paladin.Talents.Vengeance)

	vengeanceAura := paladin.RegisterAura(core.Aura{
		Label:     "Vengeance" + paladin.Label,
		ActionID:  core.ActionID{SpellID: 20049},
		Duration:  time.Second * 30,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			paladin.PseudoStats.DamageDealtMultiplier /= 1 + bonusPerStack*float64(oldStacks)
			paladin.PseudoStats.DamageDealtMultiplier *= 1 + bonusPerStack*float64(newStacks)
		},
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label: "Vengeance Trigger" + paladin.Label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeCrit) {
				vengeanceAura.Activate(sim)
				vengeanceAura.AddStack(sim)
			}
		},
	}))
}

func (paladin *Paladin) WeaponSpecializationMultiplier() float64 {
	mhWeapon := paladin.GetMHWeapon()
	if mhWeapon == nil {
		return 1
	}

	if mhWeapon.HandType == proto.HandType_HandTypeTwoHand {
		return 1 + 0.02*float64(paladin.Talents.TwoHandedWeaponSpecialization)
	}

	return 1 + 0.01*float64(paladin.Talents.OneHandedWeaponSpecialization)
}

func (paladin *Paladin) crusadeMultiplier() float64 {
	if paladin.CurrentTarget == nil || paladin.Talents.Crusade == 0 {
		return 1
	}

	switch paladin.CurrentTarget.MobType {
	case proto.MobType_MobTypeHumanoid, proto.MobType_MobTypeDemon, proto.MobType_MobTypeUndead, proto.MobType_MobTypeElemental:
		return 1 + 0.01*float64(paladin.Talents.Crusade)
	default:
		return 1
	}
}

func (paladin *Paladin) MeleeCritMultiplier() float64 {
	return paladin.DefaultMeleeCritMultiplier() * paladin.crusadeMultiplier()
}

func (paladin *Paladin) SpellCritMultiplier() float64 {
	return paladin.DefaultSpellCritMultiplier() * paladin.crusadeMultiplier()
}
