package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

const (
	JudgementManaCost = 147
	JudgementCD       = time.Second * 10
)

func (paladin *Paladin) judgementCost() int32 {
	return int32(float64(JudgementManaCost) * (1 - 0.03*float64(paladin.Talents.Benediction)))
}

func (paladin *Paladin) sanctifiedJudgement(sim *core.Simulation, sealCost float64) {
	if paladin.Talents.SanctifiedJudgement == 0 {
		return
	}

	procChance := 0.33 * float64(paladin.Talents.SanctifiedJudgement)
	if paladin.Talents.SanctifiedJudgement == 3 {
		procChance = 1
	}
	if sim.RandomFloat("Sanctified Judgement") < procChance {
		if paladin.SanctifiedJudgementMetrics == nil {
			paladin.SanctifiedJudgementMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 31930})
		}
		paladin.AddMana(sim, sealCost*0.8, paladin.SanctifiedJudgementMetrics)
	}
}

func (paladin *Paladin) canJudgement(sim *core.Simulation) bool {
	return paladin.CurrentSeal != nil && paladin.CurrentSeal.IsActive() && paladin.JudgementOfBlood.IsReady(sim)
}

func (paladin *Paladin) registerJudgementOfBloodSpell(cdTimer *core.Timer) {
	paladin.JudgementOfBlood = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31898},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ManaCost: core.ManaCostOptions{
			FlatCost: paladin.judgementCost(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: JudgementCD - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
		},

		DamageMultiplier: paladin.WeaponSpecializationMultiplier(),
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.43,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.CalcAndRollDamageRange(sim, 295, 0.10) + 0.43*spell.SpellDamage(target)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			spell.DealDamage(sim, result)

			if paladin.SpiritualAttunementMetrics != nil && result.Landed() {
				paladin.AddMana(sim, result.Damage*0.033, paladin.SpiritualAttunementMetrics)
			}

			if paladin.SealOfBloodAura != nil && paladin.SealOfBloodAura.IsActive() {
				paladin.sanctifiedJudgement(sim, retSealCost(paladin.SealOfBlood))
				paladin.SealOfBloodAura.Deactivate(sim)
				if paladin.CurrentSeal == paladin.SealOfBloodAura {
					paladin.CurrentSeal = nil
				}
			}
		},
	})
}

func (paladin *Paladin) CanJudgementOfBlood(sim *core.Simulation, target *core.Unit) bool {
	return paladin.canJudgement(sim) &&
		paladin.CurrentSeal == paladin.SealOfBloodAura &&
		paladin.JudgementOfBlood.CanCast(sim, target)
}

func (paladin *Paladin) registerJudgementOfWisdomSpell(cdTimer *core.Timer) {
	paladin.JudgementOfWisdomAura = core.JudgementOfWisdomAura(paladin.CurrentTarget)

	paladin.JudgementOfWisdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27164},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,

		ManaCost: core.ManaCostOptions{
			FlatCost: paladin.judgementCost(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: JudgementCD - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				paladin.JudgementOfWisdomAura.Activate(sim)
				paladin.CurrentJudgement = paladin.JudgementOfWisdomAura
			}
			spell.DealOutcome(sim, result)

			if paladin.SealOfWisdomAura != nil && paladin.SealOfWisdomAura.IsActive() {
				paladin.sanctifiedJudgement(sim, retSealCost(paladin.SealOfWisdom))
				paladin.SealOfWisdomAura.Deactivate(sim)
				if paladin.CurrentSeal == paladin.SealOfWisdomAura {
					paladin.CurrentSeal = nil
				}
			}
		},
	})
}

func (paladin *Paladin) CanJudgementOfWisdom(sim *core.Simulation, target *core.Unit) bool {
	return paladin.canJudgement(sim) &&
		paladin.CurrentSeal == paladin.SealOfWisdomAura &&
		paladin.JudgementOfWisdom.CanCast(sim, target)
}

func (paladin *Paladin) registerJudgementOfTheCrusaderSpell(cdTimer *core.Timer) {
	paladin.JudgementOfTheCrusaderAura = core.ImprovedSealOfTheCrusaderAura(paladin.CurrentTarget)

	paladin.JudgementOfTheCrusader = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27159},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,

		ManaCost: core.ManaCostOptions{
			FlatCost: paladin.judgementCost(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: JudgementCD - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				paladin.JudgementOfTheCrusaderAura.Activate(sim)
				paladin.CurrentJudgement = paladin.JudgementOfTheCrusaderAura
			}
			spell.DealOutcome(sim, result)

			if paladin.SealOfTheCrusaderAura != nil && paladin.SealOfTheCrusaderAura.IsActive() {
				paladin.sanctifiedJudgement(sim, retSealCost(paladin.SealOfTheCrusader))
				paladin.SealOfTheCrusaderAura.Deactivate(sim)
				if paladin.CurrentSeal == paladin.SealOfTheCrusaderAura {
					paladin.CurrentSeal = nil
				}
			}
		},
	})
}

func (paladin *Paladin) CanJudgementOfTheCrusader(sim *core.Simulation, target *core.Unit) bool {
	return paladin.canJudgement(sim) &&
		paladin.CurrentSeal == paladin.SealOfTheCrusaderAura &&
		paladin.JudgementOfTheCrusader.CanCast(sim, target)
}

func (paladin *Paladin) setupJudgementRefresh() {
	paladin.RegisterAura(core.Aura{
		Label:    "Judgement Refresh" + paladin.Label,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}
			if paladin.CurrentJudgement != nil && paladin.CurrentJudgement.IsActive() {
				paladin.CurrentJudgement.Refresh(sim)
			}
		},
	})
}

func (paladin *Paladin) registerJudgements() {
	cdTimer := paladin.NewTimer()
	paladin.registerJudgementOfBloodSpell(cdTimer)
	paladin.registerJudgementOfWisdomSpell(cdTimer)
	paladin.registerJudgementOfTheCrusaderSpell(cdTimer)
}

func retSealCost(spell *core.Spell) float64 {
	if spell == nil {
		return 0
	}
	return float64(spell.Cost.BaseCost)
}

