import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Profession, PseudoStat, Race, Stat } from '../../core/proto/common.js';
import { PaladinSeal, RetributionPaladin_Options as RetributionPaladinOptions } from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1_Gear from './gear_sets/p1.gear.json';
import Preraid_Gear from './gear_sets/preraid.gear.json';

export const P1_GEAR_PRESET = PresetUtils.makePresetGear('P1', P1_Gear);
export const PRERAID_GEAR_PRESET = PresetUtils.makePresetGear('Pre-raid', Preraid_Gear);

export const APL_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.0,
			[Stat.StatAgility]: 0.55,
			[Stat.StatAttackPower]: 0.42,
			[Stat.StatMeleeHitRating]: 0.95,
			[Stat.StatMeleeCritRating]: 0.78,
			[Stat.StatMeleeHasteRating]: 0.58,
			[Stat.StatArmorPenetration]: 0.30,
			[Stat.StatExpertiseRating]: 0.75,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.0,
		},
	),
);

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '5-503201-0523005130033125231051',
	}),
};

export const DefaultOptions = RetributionPaladinOptions.create({
	classOptions: {
		seal: PaladinSeal.Truth,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22854, // Flask of Relentless Assault
	foodId: 27658, // Roasted Clefthoof
	potId: 22838, // Haste Potion
	conjuredId: 12662, // Demonic Rune
	superSapper: true,
	goblinSapper: true,
	drumsId: 351355, // Greater Drums of Battle
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
	race: Race.RaceBloodElf,
};
