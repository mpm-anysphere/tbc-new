package warlock

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var TalentTreeSizes = [3]int{21, 22, 21}

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	ShadowBolt           *core.Spell
	Corruption           *core.Spell
	CurseOfAgony         *core.Spell
	CurseOfElementsAuras core.AuraArray
	Immolate             *core.Spell
	Incinerate           *core.Spell
	UnstableAffliction   *core.Spell
	Hellfire             *core.Spell
	DrainLife            *core.Spell
	SiphonLife           *core.Spell
	LifeTap              *core.Spell

	NightfallProcAura  *core.Aura
	ImpShadowboltAuras core.AuraArray

	ActivePet  *WarlockPet
	Felhunter  *WarlockPet
	Felguard   *WarlockPet
	Imp        *WarlockPet
	Succubus   *WarlockPet
	Voidwalker *WarlockPet
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func RegisterWarlock() {
	core.RegisterAgentFactory(
		proto.Player_Warlock{},
		proto.Spec_SpecWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewWarlock(character, options, options.GetWarlock().Options.ClassOptions)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Warlock)
			if !ok {
				panic("Invalid spec value for Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func (warlock *Warlock) ApplyTalents() {
	warlock.applyTalents()
}

func (warlock *Warlock) Initialize() {
	warlock.registerCurseOfElements()
	warlock.registerCurseOfAgony()
	warlock.registerLifeTap()
	warlock.registerShadowBolt()
	warlock.registerImmolate()
	warlock.registerIncinerate()
	warlock.RegisterCorruption()
	warlock.registerUnstableAffliction()
	warlock.registerSiphonLife()
	warlock.RegisterDrainLife()
	warlock.RegisterHellfire()
	warlock.applyTalentsPostSpellRegistration()

	// Fel Armor passive.
	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label:    "Fel Armor",
			ActionID: core.ActionID{SpellID: 28189},
		}))
	warlock.MultiplyStat(stats.Stamina, 1.1)
	warlock.MultiplyStat(stats.Health, 1.1)

	// 5% Intellect passive.
	warlock.MultiplyStat(stats.Intellect, 1.05)
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {

}

func (warlock *Warlock) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	bloodPact := core.MakeTristateValue(
		warlock.Options.GetSummon() == proto.WarlockOptions_Imp && !warlock.Talents.GetDemonicSacrifice(),
		warlock.Talents.GetImprovedImp() == 3,
	)
	if bloodPact > partyBuffs.BloodPact {
		partyBuffs.BloodPact = bloodPact
	}
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
}

func (warlock *Warlock) OnEncounterStart(_ *core.Simulation) {
}

func NewWarlock(character *core.Character, options *proto.Player, warlockOptions *proto.WarlockOptions) *Warlock {
	if warlockOptions == nil {
		warlockOptions = &proto.WarlockOptions{}
	}

	warlock := &Warlock{
		Character: *character,
		Talents:   &proto.WarlockTalents{},
		Options:   warlockOptions,
	}
	core.FillTalentsProto(warlock.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	warlock.EnableManaBar()
	warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)

	warlock.registerPets()
	warlock.setupDemonicSacrifice()

	return warlock
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
}

const (
	WarlockSpellFlagNone    int64 = 0
	WarlockSpellConflagrate int64 = 1 << iota
	WarlockSpellFaBConflagrate
	WarlockSpellShadowBolt
	WarlockSpellChaosBolt
	WarlockSpellImmolate
	WarlockSpellImmolateDot
	WarlockSpellIncinerate
	WarlockSpellFaBIncinerate
	WarlockSpellSoulFire
	WarlockSpellShadowBurn
	WarlockSpellLifeTap
	WarlockSpellCorruption
	WarlockSpellHaunt
	WarlockSpellUnstableAffliction
	WarlockSpellCurseOfElements
	WarlockSpellAgony
	WarlockSpellDrainSoul
	WarlockSpellDrainLife
	WarlockSpellMetamorphosis
	WarlockSpellSeedOfCorruption
	WarlockSpellSeedOfCorruptionExposion
	WarlockSpellHandOfGuldan
	WarlockSpellHellfire
	WarlockSpellImmolationAura
	WarlockSpellSearingPain
	WarlockSpellSummonDoomguard
	WarlockSpellDoomguardDoomBolt
	WarlockSpellSummonFelguard
	WarlockSpellFelGuardLegionStrike
	WarlockSpellFelGuardFelstorm
	WarlockSpellSummonImp
	WarlockSpellImpFireBolt
	WarlockSpellSummonFelhunter
	WarlockSpellFelHunterShadowBite
	WarlockSpellSummonSuccubus
	WarlockSpellSuccubusLashOfPain
	WarlockSpellVoidwalkerTorment
	WarlockSpellSummonInfernal
	WarlockSpellDemonSoul
	WarlockSpellShadowflame
	WarlockSpellShadowflameDot
	WarlockSpellSoulBurn
	WarlockSpellFelFlame
	WarlockSpellBurningEmbers
	WarlockSpellEmberTap
	WarlockSpellRainOfFire
	WarlockSpellFireAndBrimstone
	WarlockSpellDarkSoulInsanity
	WarlockSpellDarkSoulKnowledge
	WarlockSpellDarkSoulMisery
	WarlockSpellMaleficGrasp
	WarlockSpellDemonicSlash
	WarlockSpellTouchOfChaos
	WarlockSpellChaosWave
	WarlockSpellCarrionSwarm
	WarlockSpellDoom
	WarlockSpellVoidray
	WarlockSpellSiphonLife
	WarlockSpellHavoc
	WarlockSpellAll int64 = 1<<iota - 1

	WarlockShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellHaunt |
		WarlockSpellDrainSoul | WarlockSpellDrainLife | WarlockSpellAgony |
		WarlockSpellShadowBolt | WarlockSpellSeedOfCorruptionExposion | WarlockSpellHandOfGuldan |
		WarlockSpellShadowflame | WarlockSpellFelFlame | WarlockSpellChaosBolt | WarlockSpellShadowBurn | WarlockSpellHavoc

	WarlockPeriodicShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellAgony

	WarlockFireDamage = WarlockSpellConflagrate | WarlockSpellImmolate | WarlockSpellIncinerate | WarlockSpellSoulFire |
		WarlockSpellHandOfGuldan | WarlockSpellSearingPain | WarlockSpellImmolateDot |
		WarlockSpellShadowflameDot | WarlockSpellFelFlame | WarlockSpellChaosBolt | WarlockSpellShadowBurn | WarlockSpellFaBConflagrate |
		WarlockSpellFaBIncinerate

	WarlockDoT = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellAgony | WarlockSpellImmolateDot |
		WarlockSpellShadowflameDot | WarlockSpellBurningEmbers

	WarlockSummonSpells = WarlockSpellSummonImp | WarlockSpellSummonSuccubus | WarlockSpellSummonFelhunter |
		WarlockSpellSummonFelguard

	WarlockDarkSoulSpell             = WarlockSpellDarkSoulInsanity | WarlockSpellDarkSoulKnowledge | WarlockSpellDarkSoulMisery
	WarlockAllSummons                = WarlockSummonSpells | WarlockSpellSummonInfernal | WarlockSpellSummonDoomguard
	WarlockSpellsChaoticEnergyDestro = WarlockSpellAll &^ WarlockAllSummons &^ WarlockSpellDrainLife
)
