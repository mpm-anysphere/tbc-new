package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

const (
	SealDuration = time.Second * 30
	TwistWindow  = time.Millisecond * 400
)

func (paladin *Paladin) setupSealOfBlood() {
	procActionID := core.ActionID{SpellID: 31893}

	sealProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    procActionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeProc,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,

		DamageMultiplier: 0.35 * paladin.WeaponSpecializationMultiplier(),
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
			spell.DealDamage(sim, result)

			if result.Landed() && paladin.SpiritualAttunementMetrics != nil {
				// Seal of Blood inflicts self-damage and Spiritual Attunement returns a
				// fraction of that as mana. Approximation used for alpha tuning.
				paladin.AddMana(sim, result.Damage*0.01, paladin.SpiritualAttunementMetrics)
			}
		},
	})

	paladin.SealOfBloodAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Blood" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 31892},
		Duration: SealDuration,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.CurrentSeal == aura {
				paladin.CurrentSeal = nil
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			sealProc.Cast(sim, result.Target)
		},
	})

	paladin.SealOfBlood = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31892},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagHelpful,
		ManaCost: core.ManaCostOptions{
			FlatCost: 210,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.UpdateSeal(sim, paladin.SealOfBloodAura)
		},
	})
}

func (paladin *Paladin) setupSealOfCommand() {
	procActionID := core.ActionID{SpellID: 20424}
	dpm := paladin.NewLegacyPPMManager(7.0, core.ProcMaskMeleeMH)
	icd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second,
	}

	sealProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    procActionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeProc,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,

		DamageMultiplier: 0.7 * paladin.WeaponSpecializationMultiplier(),
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.29,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.MHWeaponDamage(sim, spell.MeleeAttackPower(target)) + 0.29*spell.SpellDamage(target)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})

	paladin.SealOfCommandAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Command" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 20375},
		Duration: SealDuration,
		Dpm:      dpm,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.CurrentSeal == aura {
				paladin.CurrentSeal = nil
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			if !dpm.Proc(sim, spell.ProcMask, "Seal of Command"+paladin.Label) {
				return
			}

			icd.Use(sim)
			sealProc.Cast(sim, result.Target)
		},
	})

	paladin.SealOfCommand = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20375},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagHelpful,
		ManaCost: core.ManaCostOptions{
			FlatCost: 65,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.UpdateSeal(sim, paladin.SealOfCommandAura)
		},
	})
}

func (paladin *Paladin) setupSealOfTheCrusader() {
	paladin.SealOfTheCrusaderAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of the Crusader" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 27158},
		Duration: SealDuration,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.CurrentSeal == aura {
				paladin.CurrentSeal = nil
			}
		},
	})

	paladin.SealOfTheCrusader = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27158},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagHelpful,
		ManaCost: core.ManaCostOptions{
			FlatCost: 210,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.UpdateSeal(sim, paladin.SealOfTheCrusaderAura)
		},
	})
}

func (paladin *Paladin) setupSealOfWisdom() {
	paladin.SealOfWisdomAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Wisdom" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 27166},
		Duration: SealDuration,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.CurrentSeal == aura {
				paladin.CurrentSeal = nil
			}
		},
	})

	paladin.SealOfWisdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27166},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagHelpful,
		ManaCost: core.ManaCostOptions{
			FlatCost: 270,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.UpdateSeal(sim, paladin.SealOfWisdomAura)
		},
	})
}

func (paladin *Paladin) setupSealOfRighteousness() {
	procActionID := core.ActionID{SpellID: 27156}

	sealProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    procActionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeProc,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,

		DamageMultiplier: paladin.WeaponSpecializationMultiplier(),
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.25*paladin.MHWeaponDamage(sim, spell.MeleeAttackPower(target)) + 0.1*spell.SpellDamage(target)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
			spell.DealDamage(sim, result)
		},
	})

	paladin.SealOfRighteousnessAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Righteousness" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 27155},
		Duration: SealDuration,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.CurrentSeal == aura {
				paladin.CurrentSeal = nil
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			sealProc.Cast(sim, result.Target)
		},
	})

	paladin.SealOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27155},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagHelpful,
		ManaCost: core.ManaCostOptions{
			FlatCost: 260,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.UpdateSeal(sim, paladin.SealOfRighteousnessAura)
		},
	})
}

func (paladin *Paladin) UpdateSeal(sim *core.Simulation, newSeal *core.Aura) {
	if paladin.CurrentSeal == paladin.SealOfCommandAura && paladin.SealOfCommandAura.IsActive() {
		// Seal twisting behavior: while swapping from Seal of Command, keep it
		// alive for the twist window so both seals can interact on the next swing.
		paladin.SealOfCommandAura.Duration = TwistWindow
		paladin.SealOfCommandAura.Refresh(sim)
		paladin.SealOfCommandAura.Duration = SealDuration
	} else if paladin.CurrentSeal != nil {
		paladin.CurrentSeal.Deactivate(sim)
	}

	paladin.CurrentSeal = newSeal
	if newSeal != nil {
		newSeal.Activate(sim)
	}
}

