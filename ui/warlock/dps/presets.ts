import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs, Stat, TristateEffect } from '../../core/proto/common';
import { Warlock_Options as WarlockOptions } from '../../core/proto/warlock';
import { defaultRaidBuffMajorDamageCooldowns } from '../../core/proto_utils/utils';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import DefaultAPL from './apls/default.apl.json';
import BlankGear from './gear_sets/blank.gear.json';
import PreBisGear from './gear_sets/prebis.gear.json';
import P1Gear from './gear_sets/p1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const WARLOCK_DEFAULT_APL = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);

export const BLANK_GEARSET = PresetUtils.makePresetGear('Blank', BlankGear);
export const PREBIS_GEAR = PresetUtils.makePresetGear('Pre-Raid', PreBisGear);
export const P1_BIS_GEAR = PresetUtils.makePresetGear('P1', P1Gear);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatIntellect]: 0.42,
			[Stat.StatSpirit]: 0.1,
			[Stat.StatSpellDamage]: 1.0,
			[Stat.StatSpellHitRating]: 1.2,
			[Stat.StatSpellCritRating]: 0.74,
			[Stat.StatSpellHasteRating]: 1.06,
		},
		{},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const Talents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '000000000000000000000-0000000000000000000000-000000000000000000000',
	}),
};

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		summon: 0,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	guardianElixirId: 32067, // Elixir of Draenic Wisdom
	battleElixirId: 28103, // Adept's Elixir
	foodId: 27657, // Blackened Basilisk
	mhImbueId: 25122, // Brilliant Wizard Oil
	potId: 22839, // Destruction Potion
	conjuredId: 12662, // Demonic Rune
	drumsId: 351355, // Greater Drums of Battle
});

export const OtherDefaults = {
	distanceFromTarget: 20,
};

export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(),
	divineSpirit: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	shadowProtection: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	manaTideTotems: 1,
	wrathOfAirTotem: 1,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	misery: true,
	curseOfElements: TristateEffect.TristateEffectImproved,
	improvedSealOfTheCrusader: true,
	judgementOfWisdom: true,
});
