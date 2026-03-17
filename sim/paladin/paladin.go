package paladin

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var TalentTreeSizes = [3]int{20, 23, 22}

type Paladin struct {
	core.Character

	Seal proto.PaladinSeal

	Talents *proto.PaladinTalents

	CurrentSeal *core.Aura
	// Tracks the currently active judged debuff we are responsible for.
	CurrentJudgement *core.Aura

	AvengingWrathAura *core.Aura

	CrusaderStrike *core.Spell
	Consecration   *core.Spell
	Exorcism       *core.Spell

	JudgementOfBlood       *core.Spell
	JudgementOfWisdom      *core.Spell
	JudgementOfTheCrusader *core.Spell

	SealOfBlood         *core.Spell
	SealOfCommand       *core.Spell
	SealOfTheCrusader   *core.Spell
	SealOfWisdom        *core.Spell
	SealOfRighteousnessAura *core.Aura
	SealOfRighteousness     *core.Spell

	SealOfBloodAura      *core.Aura
	SealOfCommandAura    *core.Aura
	SealOfTheCrusaderAura *core.Aura
	SealOfWisdomAura     *core.Aura

	JudgementOfTheCrusaderAura *core.Aura
	JudgementOfWisdomAura      *core.Aura

	SpiritualAttunementMetrics *core.ResourceMetrics

	DefensiveCooldownAuras []*core.Aura
}

// Implemented by each Paladin spec.
type PaladinAgent interface {
	GetPaladin() *Paladin
}

func (paladin *Paladin) GetCharacter() *core.Character {
	return &paladin.Character
}

func (paladin *Paladin) GetPaladin() *Paladin {
	return paladin
}

func (paladin *Paladin) AddRaidBuffs(_ *proto.RaidBuffs) {
}

func (paladin *Paladin) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if paladin.Talents.SanctityAura {
		effect := core.MakeTristateValue(true, paladin.Talents.ImprovedSanctityAura == 2)
		if partyBuffs.SanctityAura < effect {
			partyBuffs.SanctityAura = effect
		}
		return
	}

	if paladin.Talents.ImprovedDevotionAura > 0 {
		effect := core.MakeTristateValue(true, paladin.Talents.ImprovedDevotionAura == 5)
		if partyBuffs.DevotionAura < effect {
			partyBuffs.DevotionAura = effect
		}
		return
	}

	if paladin.Talents.ImprovedRetributionAura > 0 {
		effect := core.MakeTristateValue(true, paladin.Talents.ImprovedRetributionAura == 2)
		if partyBuffs.RetributionAura < effect {
			partyBuffs.RetributionAura = effect
		}
	}
}

func (paladin *Paladin) Initialize() {
	paladin.registerSpells()
}

func (paladin *Paladin) registerSpells() {
	paladin.setupSealOfBlood()
	paladin.setupSealOfCommand()
	paladin.setupSealOfTheCrusader()
	paladin.setupSealOfWisdom()
	paladin.setupSealOfRighteousness()
	paladin.setupJudgementRefresh()

	paladin.registerCrusaderStrikeSpell()
	paladin.registerExorcismSpell()
	paladin.registerJudgements()
	paladin.registerSpiritualAttunement()
}

func (paladin *Paladin) Reset(sim *core.Simulation) {
	if paladin.CurrentSeal != nil {
		paladin.CurrentSeal.Deactivate(sim)
		paladin.CurrentSeal = nil
	}
	if paladin.CurrentJudgement != nil {
		paladin.CurrentJudgement.Deactivate(sim)
		paladin.CurrentJudgement = nil
	}

	switch paladin.Seal {
	case proto.PaladinSeal_Truth:
		paladin.CurrentSeal = paladin.SealOfBloodAura
		paladin.SealOfBloodAura.Activate(sim)
	case proto.PaladinSeal_Justice:
		paladin.CurrentSeal = paladin.SealOfCommandAura
		paladin.SealOfCommandAura.Activate(sim)
	case proto.PaladinSeal_Insight:
		paladin.CurrentSeal = paladin.SealOfWisdomAura
		paladin.SealOfWisdomAura.Activate(sim)
	case proto.PaladinSeal_Righteousness:
		paladin.CurrentSeal = paladin.SealOfRighteousnessAura
		paladin.SealOfRighteousnessAura.Activate(sim)
	}
}

func (paladin *Paladin) OnEncounterStart(sim *core.Simulation) {
}

func NewPaladin(character *core.Character, talentsStr string, options *proto.PaladinOptions) *Paladin {
	if options == nil {
		options = &proto.PaladinOptions{}
	}

	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},
		Seal:      options.Seal,
	}

	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	paladin.PseudoStats.CanParry = true
	paladin.PseudoStats.BaseDodgeChance += 0.0065
	paladin.PseudoStats.BaseParryChance += 0.05
	paladin.PseudoStats.BaseBlockChance += 0.05

	paladin.EnableManaBar()

	// Only retribution and holy are actually pets performing some kind of action
	// if paladin.Spec != proto.Spec_SpecProtectionPaladin {
	// 	paladin.AncientGuardian = paladin.NewAncientGuardian()
	// }

	paladin.EnableAutoAttacks(paladin, core.AutoAttackOptions{
		MainHand:       paladin.WeaponFromMainHand(paladin.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	paladin.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])
	paladin.AddStatDependency(stats.Intellect, stats.SpellCritPercent, core.CritPerIntMaxLevel[character.Class])
	paladin.AddStatDependency(stats.Agility, stats.DodgeRating, 1/25.0*core.DodgeRatingPerDodgePercent)

	// Bonus Armor and Armor are treated identically for Paladins
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	return paladin
}

func (paladin *Paladin) AddDefensiveCooldownAura(aura *core.Aura) {
	paladin.DefensiveCooldownAuras = append(paladin.DefensiveCooldownAuras, aura)
}

func (paladin *Paladin) AnyActiveDefensiveCooldown() bool {
	for _, aura := range paladin.DefensiveCooldownAuras {
		if aura.IsActive() {
			return true
		}
	}

	return false
}
