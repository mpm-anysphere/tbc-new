package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (pet *WarlockPet) registerShadowBiteSpell(_ *Warlock) {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 54049},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellFelHunterShadowBite,

		ManaCost: core.ManaCostOptions{FlatCost: 190},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: 6 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.38,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 150.0
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}))
}

func (pet *WarlockPet) registerFireboltSpell(owner *Warlock) {
	castTimeReduction := time.Millisecond * 250 * time.Duration(owner.Talents.GetImprovedFirebolt())
	damageMultiplier := 1 + (0.1 * float64(owner.Talents.GetImprovedImp()))

	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 27267},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellImpFireBolt,
		MissileSpeed:   16,

		ManaCost: core.ManaCostOptions{FlatCost: 190},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: max(0, 2*time.Second-castTimeReduction),
			},
		},

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.571,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := pet.CalcAndRollDamageRange(sim, 112, 127)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	}))
}

func (pet *WarlockPet) registerLashOfPainSpell(owner *Warlock) {
	cooldownReduction := time.Second * 3 * time.Duration(owner.Talents.GetImprovedLashOfPain())
	damageMultiplier := 1 + (0.1 * float64(owner.Talents.GetImprovedSayaad()))

	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 27274},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellSuccubusLashOfPain,

		ManaCost: core.ManaCostOptions{FlatCost: 190},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: max(0, 12*time.Second-cooldownReduction),
			},
		},

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.429,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 123, spell.OutcomeMagicHitAndCrit)
		},
	}))
}

func (pet *WarlockPet) registerTormentSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3716},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellVoidwalkerTorment,

		ManaCost: core.ManaCostOptions{FlatCost: 170},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: 5 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.3,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 150, spell.OutcomeMagicHitAndCrit)
		},
	}))
}

func (pet *WarlockPet) registerFelguardStrikeSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30213},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: WarlockSpellFelGuardLegionStrike,
		Flags:          core.SpellFlagMeleeMetrics,

		ManaCost: core.ManaCostOptions{FlatCost: 170},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: 6 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := pet.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower()) + 78
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	}))
}
