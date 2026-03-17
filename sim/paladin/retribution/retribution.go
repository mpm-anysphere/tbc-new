package retribution

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/paladin"
)

func RegisterRetributionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_RetributionPaladin{},
		proto.Spec_SpecRetributionPaladin,
		func(character *core.Character, options *proto.Player, _ *proto.Raid) core.Agent {
			return NewRetributionPaladin(character, options)
		},
		func(player *proto.Player, spec any) {
			playerSpec, ok := spec.(*proto.Player_RetributionPaladin)
			if !ok {
				panic("Invalid spec value for Retribution Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRetributionPaladin(character *core.Character, options *proto.Player) *RetributionPaladin {
	retOptions := options.GetRetributionPaladin()
	classOptions := &proto.PaladinOptions{}
	if retOptions != nil && retOptions.Options != nil && retOptions.Options.ClassOptions != nil {
		classOptions = retOptions.Options.ClassOptions
	}

	ret := &RetributionPaladin{
		Paladin: paladin.NewPaladin(character, options.TalentsString, classOptions),
	}

	// Backward-compatible safety: if no talents are provided, use a standard
	// TBC Ret baseline so alpha defaults remain functional.
	if options.TalentsString == "" {
		applyDefaultRetTalents(ret.Talents)
	}

	return ret
}

type RetributionPaladin struct {
	*paladin.Paladin

	openerCompleted bool
}

func (ret *RetributionPaladin) GetPaladin() *paladin.Paladin {
	return ret.Paladin
}

func (ret *RetributionPaladin) Initialize() {
	ret.Paladin.Initialize()
	ret.RegisterConsecrationSpell(6)
	ret.RegisterAvengingWrathCD()
}

func (ret *RetributionPaladin) ApplyTalents() {
	ret.Paladin.ApplyTalents()
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)

	// Legacy parity opener: start with Seal of the Crusader active so we can
	// immediately cast Judgement of the Crusader at pull.
	if ret.SealOfTheCrusaderAura != nil {
		if ret.CurrentSeal != nil {
			ret.CurrentSeal.Deactivate(sim)
		}
		ret.CurrentSeal = ret.SealOfTheCrusaderAura
		ret.SealOfTheCrusaderAura.Activate(sim)
	}

	ret.AutoAttacks.CancelAutoSwing(sim)
	ret.openerCompleted = false
}

func applyDefaultRetTalents(talents *proto.PaladinTalents) {
	talents.DivineStrength = 5
	talents.Benediction = 5
	talents.ImprovedJudgement = 2
	talents.ImprovedSealOfTheCrusader = 3
	talents.Conviction = 5
	talents.SealOfCommand = true
	talents.Crusade = 3
	talents.TwoHandedWeaponSpecialization = 3
	talents.SanctityAura = true
	talents.ImprovedSanctityAura = 2
	talents.Vengeance = 5
	talents.SanctifiedJudgement = 2
	talents.SanctifiedSeals = 3
	talents.Fanaticism = 5
	talents.CrusaderStrike = true
	talents.Precision = 3
}
