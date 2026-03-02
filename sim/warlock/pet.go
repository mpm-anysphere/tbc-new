package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	AutoCastAbilities []*core.Spell
}

var petBaseStats = map[proto.WarlockOptions_Summon]stats.Stats{
	proto.WarlockOptions_Imp: {
		stats.Health:    2800,
		stats.Mana:      2200,
		stats.Stamina:   101,
		stats.Strength:  145,
		stats.Agility:   38,
		stats.Intellect: 327,
		stats.Spirit:    263,
	},
	proto.WarlockOptions_Voidwalker: {
		stats.Health:    4200,
		stats.Mana:      1600,
		stats.Stamina:   340,
		stats.Strength:  185,
		stats.Agility:   72,
		stats.Intellect: 125,
		stats.Spirit:    135,
	},
	proto.WarlockOptions_Succubus: {
		stats.Health:      3400,
		stats.Mana:        1850,
		stats.Stamina:     280,
		stats.Strength:    153,
		stats.Agility:     108,
		stats.Intellect:   133,
		stats.Spirit:      122,
		stats.AttackPower: 20,
	},
	proto.WarlockOptions_Felhunter: {
		stats.Health:      3500,
		stats.Mana:        1850,
		stats.Stamina:     280,
		stats.Strength:    153,
		stats.Agility:     108,
		stats.Intellect:   133,
		stats.Spirit:      122,
		stats.AttackPower: 20,
	},
	proto.WarlockOptions_Felguard: {
		stats.Health:      4000,
		stats.Mana:        1800,
		stats.Stamina:     320,
		stats.Strength:    185,
		stats.Agility:     108,
		stats.Intellect:   133,
		stats.Spirit:      122,
		stats.AttackPower: 35,
	},
}

func (warlock *Warlock) petStatInheritance(ownerStats stats.Stats) stats.Stats {
	return stats.Stats{
		stats.Stamina:          ownerStats[stats.Stamina] * 0.3,
		stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
		stats.Armor:            ownerStats[stats.Armor] * 0.35,
		stats.AttackPower:      (ownerStats[stats.SpellDamage] + ownerStats[stats.ShadowDamage]) * 0.57,
		stats.SpellDamage:      (ownerStats[stats.SpellDamage] + ownerStats[stats.ShadowDamage]) * 0.15,
		stats.SpellPenetration: ownerStats[stats.SpellPenetration],
	}
}

func (warlock *Warlock) petMeleeAutoAttack(swingSpeed float64) *core.AutoAttackOptions {
	return &core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  83.4,
			BaseDamageMax:  123.4,
			SwingSpeed:     swingSpeed,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	}
}

func (warlock *Warlock) registerPets() {
	warlock.Imp = warlock.registerImp()
	warlock.Succubus = warlock.registerSuccubus()
	warlock.Felhunter = warlock.registerFelhunter()
	warlock.Voidwalker = warlock.registerVoidwalker()
	warlock.Felguard = warlock.registerFelguard()
}

func (warlock *Warlock) registerImp() *WarlockPet {
	enabledOnStart := warlock.Options.GetSummon() == proto.WarlockOptions_Imp
	return warlock.makePet(proto.WarlockOptions_Imp, enabledOnStart, false)
}

func (warlock *Warlock) registerSuccubus() *WarlockPet {
	enabledOnStart := warlock.Options.GetSummon() == proto.WarlockOptions_Succubus
	return warlock.makePet(proto.WarlockOptions_Succubus, enabledOnStart, true)
}

func (warlock *Warlock) registerFelhunter() *WarlockPet {
	enabledOnStart := warlock.Options.GetSummon() == proto.WarlockOptions_Felhunter
	return warlock.makePet(proto.WarlockOptions_Felhunter, enabledOnStart, true)
}

func (warlock *Warlock) registerVoidwalker() *WarlockPet {
	enabledOnStart := warlock.Options.GetSummon() == proto.WarlockOptions_Voidwalker
	return warlock.makePet(proto.WarlockOptions_Voidwalker, enabledOnStart, true)
}

func (warlock *Warlock) registerFelguard() *WarlockPet {
	enabledOnStart := warlock.Options.GetSummon() == proto.WarlockOptions_Felguard && warlock.Talents.GetSummonFelguard()
	return warlock.makePet(proto.WarlockOptions_Felguard, enabledOnStart, true)
}

func (warlock *Warlock) makePet(summonType proto.WarlockOptions_Summon, enabledOnStart bool, canMelee bool) *WarlockPet {
	name := proto.WarlockOptions_Summon_name[int32(summonType)]
	baseStats, ok := petBaseStats[summonType]
	if !ok {
		return nil
	}

	pet := &WarlockPet{
		Pet: core.NewPet(core.PetConfig{
			Name:                     name,
			Owner:                    &warlock.Character,
			BaseStats:                baseStats,
			NonHitExpStatInheritance: warlock.petStatInheritance,
			EnabledOnStart:           enabledOnStart,
		}),
	}

	pet.Class = pet.Owner.Class
	pet.EnableManaBar()
	pet.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	pet.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[proto.Class_ClassPaladin])
	pet.AddStatDependency(stats.Intellect, stats.SpellCritPercent, core.CritPerIntMaxLevel[proto.Class_ClassPaladin])

	if canMelee {
		pet.EnableAutoAttacks(pet, *warlock.petMeleeAutoAttack(2.0))
	}

	if enabledOnStart {
		warlock.RegisterResetEffect(func(sim *core.Simulation) {
			warlock.ActivePet = pet
		})
	}

	warlock.AddPet(pet)
	pet.registerPrimarySpell(summonType, warlock)

	return pet
}

func (pet *WarlockPet) registerPrimarySpell(summonType proto.WarlockOptions_Summon, owner *Warlock) {
	switch summonType {
	case proto.WarlockOptions_Imp:
		pet.registerFireboltSpell(owner)
	case proto.WarlockOptions_Succubus:
		pet.registerLashOfPainSpell(owner)
	case proto.WarlockOptions_Felhunter:
		pet.registerShadowBiteSpell(owner)
	case proto.WarlockOptions_Voidwalker:
		pet.registerTormentSpell()
	case proto.WarlockOptions_Felguard:
		pet.registerFelguardStrikeSpell()
	}
}

func (pet *WarlockPet) GetPet() *core.Pet {
	return &pet.Pet
}

func (pet *WarlockPet) Reset(_ *core.Simulation) {
}

func (pet *WarlockPet) OnEncounterStart(_ *core.Simulation) {
}

func (pet *WarlockPet) ExecuteCustomRotation(sim *core.Simulation) {
	if len(pet.AutoCastAbilities) == 0 {
		pet.WaitUntil(sim, sim.CurrentTime+time.Second)
		return
	}

	waitDuration := 500 * time.Millisecond
	for _, spell := range pet.AutoCastAbilities {
		if spell.CanCast(sim, pet.CurrentTarget) {
			spell.Cast(sim, pet.CurrentTarget)
			return
		}
		if spell.CD.Timer != nil {
			waitDuration = min(waitDuration, max(spell.TimeToReady(sim), 100*time.Millisecond))
		}
	}

	pet.WaitUntil(sim, sim.CurrentTime+waitDuration)
}
