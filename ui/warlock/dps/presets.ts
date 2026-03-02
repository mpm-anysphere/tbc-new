import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, Profession, RaidBuffs, Stat, TristateEffect } from '../../core/proto/common';
import { Warlock_Options as WarlockOptions, WarlockOptions_Summon as WarlockSummon } from '../../core/proto/warlock';
import { defaultRaidBuffMajorDamageCooldowns } from '../../core/proto_utils/utils';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import AfflictionAPL from './apls/affliction.apl.json';
import DefaultAPL from './apls/default.apl.json';
import BlankGear from './gear_sets/blank.gear.json';
import PreBisGear from './gear_sets/prebis.gear.json';
import PreBisShadowGear from './gear_sets/prebis_shadow.gear.json';
import P1Gear from './gear_sets/p1.gear.json';
import P1CraftedGear from './gear_sets/p1_crafted.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const ROTATION_PRESET_DESTRUCTION = PresetUtils.makePresetAPLRotation('Destruction (DS/Ruin)', DefaultAPL);
export const ROTATION_PRESET_AFFLICTION_UA = PresetUtils.makePresetAPLRotation('Affliction (UA)', AfflictionAPL);
export const ROTATION_PRESET_AFFLICTION_RUIN = PresetUtils.makePresetAPLRotation('Affliction (Ruin)', AfflictionAPL);

// Keep old exported name for migration and saved settings compatibility.
export const WARLOCK_DEFAULT_APL = ROTATION_PRESET_DESTRUCTION;

export const BLANK_GEARSET = PresetUtils.makePresetGear('Blank', BlankGear);
export const PREBIS_GEAR = PresetUtils.makePresetGear('Pre-Raid (Fire)', PreBisGear);
export const PREBIS_SHADOW_GEAR = PresetUtils.makePresetGear('Pre-Raid (Shadow)', PreBisShadowGear);
export const P1_BIS_GEAR = PresetUtils.makePresetGear('P1 - T4 Core', P1Gear);
export const P1_CRAFTED_GEAR = PresetUtils.makePresetGear('P1 - Crafted Hybrid', P1CraftedGear);

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

export const DESTRUCTION_DS_RUIN_TALENTS = PresetUtils.makePresetTalents(
	'Destruction DS/Ruin (0/21/40)',
	SavedTalents.create({
		// 0/21/40 DS/Ruin
		talentsString: '-20500301332101-55500050220001053025',
	}),
);

export const AFFLICTION_UA_TALENTS = PresetUtils.makePresetTalents(
	'Affliction UA (41/0/20)',
	SavedTalents.create({
		// 41/0/20
		talentsString: '350020251023310550031--0550005122',
	}),
);

export const AFFLICTION_RUIN_TALENTS = PresetUtils.makePresetTalents(
	'Affliction Ruin (40/0/21)',
	SavedTalents.create({
		// 40/0/21
		talentsString: '35002025102331055003--05500051220001',
	}),
);

// Keep old export wired to a valid default raid build.
export const Talents = DESTRUCTION_DS_RUIN_TALENTS;

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		summon: WarlockSummon.Succubus,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22866, // Flask of Pure Death
	foodId: 27657, // Blackened Basilisk
	mhImbueId: 25122, // Brilliant Wizard Oil
	potId: 22839, // Destruction Potion
	conjuredId: 12662, // Demonic Rune
	drumsId: 351355, // Greater Drums of Battle
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
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

const DefaultBuildSettings = {
	name: 'Default',
	playerOptions: OtherDefaults,
	specOptions: DefaultOptions,
	consumables: DefaultConsumables,
	buffs: DefaultIndividualBuffs,
	raidBuffs: DefaultRaidBuffs,
	partyBuffs: DefaultPartyBuffs,
	debuffs: DefaultDebuffs,
};

export const PREBIS_DESTRUCTION_BUILD = PresetUtils.makePresetBuild('Pre-Raid - Destruction (DS/Ruin)', {
	gear: PREBIS_GEAR,
	talents: DESTRUCTION_DS_RUIN_TALENTS,
	rotation: ROTATION_PRESET_DESTRUCTION,
	epWeights: P1_EP_PRESET,
	settings: DefaultBuildSettings,
});

export const PREBIS_AFFLICTION_UA_BUILD = PresetUtils.makePresetBuild('Pre-Raid - Affliction (UA)', {
	gear: PREBIS_SHADOW_GEAR,
	talents: AFFLICTION_UA_TALENTS,
	rotation: ROTATION_PRESET_AFFLICTION_UA,
	epWeights: P1_EP_PRESET,
	settings: DefaultBuildSettings,
});

export const P1_DESTRUCTION_BUILD = PresetUtils.makePresetBuild('P1 - Destruction (DS/Ruin)', {
	gear: P1_CRAFTED_GEAR,
	talents: DESTRUCTION_DS_RUIN_TALENTS,
	rotation: ROTATION_PRESET_DESTRUCTION,
	epWeights: P1_EP_PRESET,
	settings: DefaultBuildSettings,
});

export const P1_AFFLICTION_RUIN_BUILD = PresetUtils.makePresetBuild('P1 - Affliction (Ruin)', {
	gear: P1_CRAFTED_GEAR,
	talents: AFFLICTION_RUIN_TALENTS,
	rotation: ROTATION_PRESET_AFFLICTION_RUIN,
	epWeights: P1_EP_PRESET,
	settings: DefaultBuildSettings,
});
