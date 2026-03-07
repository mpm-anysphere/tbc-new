import * as PresetUtils from '../../core/preset_utils';
import {
	Class,
	ConsumesSpec,
	Debuffs,
	IndividualBuffs,
	PartyBuffs,
	Profession,
	PseudoStat,
	Race,
	RaidBuffs,
	Stat,
	TristateEffect,
} from '../../core/proto/common';
import { HunterOptions_PetType as PetType, Hunter_Options as HunterOptions, HunterOptions_Ammo, HunterOptions_QuiverBonus } from '../../core/proto/hunter';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import { defaultRaidBuffMajorDamageCooldowns } from '../../core/proto_utils/utils';
import DefaultAPL from './apls/default.apl.json';
import P1_2HGear from './gear_sets/p1_2h.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DEFAULT_APL = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);

export const P1_2H_GEARSET = PresetUtils.makePresetGear('P1 - 2H', P1_2HGear);

export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 1,
			[Stat.StatStrength]: 0.06,
			[Stat.StatIntellect]: 0.01,
			[Stat.StatAttackPower]: 0.06,
			[Stat.StatRangedAttackPower]: 0.4,
			[Stat.StatMeleeHitRating]: 0.12,
			[Stat.StatMeleeCritRating]: 0.92,
			[Stat.StatMeleeHasteRating]: 0.788,
			[Stat.StatArmorPenetration]: 0.16,
		},
		{
			[PseudoStat.PseudoStatRangedDps]: 1.75,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const BMTalents = {
	name: 'BM',
	data: SavedTalents.create({
		talentsString: '512002005250122431051-0505201205',
	}),
};
export const SVTalents = {
	name: 'SV',
	data: SavedTalents.create({
		talentsString: '502-0550201205-333200022003223005103',
	}),
};

export const DefaultOptions = HunterOptions.create({
	classOptions: {
		ammo: HunterOptions_Ammo.WardensArrow,
		quiverBonus: HunterOptions_QuiverBonus.Speed15,
		petType: PetType.Ravager,
		petUptime: 1,
		petSingleAbility: false,
	},
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	unleashedRage: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	braidedEterniumChain: true,
	ferociousInspiration: 1,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	windfuryTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(Class.ClassWarrior),
	arcaneBrilliance: true,
	divineSpirit: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	shadowProtection: true,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	exposeWeaknessUptime: 0.9,
	exposeWeaknessHunterAgility: 1080,
	faerieFire: TristateEffect.TristateEffectImproved,
	giftOfArthas: true,
	huntersMark: TristateEffect.TristateEffectImproved,
	improvedSealOfTheCrusader: true,
	insectSwarm: true,
	judgementOfLight: true,
	judgementOfWisdom: true,
	mangle: true,
	misery: true,
	sunderArmor: true,
});

export const DefaultConsumables = ConsumesSpec.create({
	battleElixirId: 22831, // Elixir of Major Agility
	foodId: 27659, // Warp Burger
	potId: 22838, // Haste Potion
	conjuredId: 12662,
	explosiveId: 30217,
	drumsId: 351355,
	petFoodId: 33874, // Kibler's Bits
	petScrollAgi: true,
	petScrollStr: true,
	superSapper: true,
	goblinSapper: true,
	scrollAgi: true,
	scrollStr: true,
});

export const OtherDefaults = {
	distanceFromTarget: 7,
	iterationCount: 25000,
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	race: Race.RaceOrc,
};
