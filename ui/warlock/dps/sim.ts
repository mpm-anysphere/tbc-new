import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Faction, ItemSlot, PseudoStat, Race, Spec, Stat } from '../../core/proto/common';
import { DEFAULT_CASTER_GEM_STATS, Stats, UnitStat } from '../../core/proto_utils/stats';
import { TypedEvent } from '../../core/typed_event';
import * as WarlockInputs from './inputs';
import * as Presets from './presets';

const modifyDisplayStats = (player: Player<Spec.SpecWarlock>) => {
	let stats = new Stats();

	TypedEvent.freezeAllAndDo(() => {
		const currentStats = player.getCurrentStats().finalStats?.stats;
		if (currentStats === undefined) {
			return {};
		}

		// stats = stats.addStat(Stat.StatMP5, (currentStats[Stat.StatMP5] * currentStats[Stat.StatSpellHasteRating]) / HASTE_RATING_PER_HASTE_PERCENT / 100);
	});

	return {
		talents: stats,
	};
};

const SPEC_CONFIG = registerSpecConfig(Spec.SpecWarlock, {
	cssClass: 'warlock-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Warlock),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: ['Alpha preview: baseline rotation and presets are implemented, pet specialization modeling is still in progress.'],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellDamage, Stat.StatSpellHitRating, Stat.StatSpellCritRating, Stat.StatSpellHasteRating],
	// Reference stat against which to calculate EP. DPS classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellDamage,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatHealth,
			Stat.StatMana,
			Stat.StatStamina,
			Stat.StatIntellect,
			Stat.StatSpellDamage,
			Stat.StatMP5,
		],
		[PseudoStat.PseudoStatSpellHitPercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHastePercent],
	),
	gemStats: DEFAULT_CASTER_GEM_STATS,

	modifyDisplayStats,
	defaults: {
		// Default equipped gear.
		gear: Presets.PREBIS_GEAR.gear,

		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,

		// Default talents.
		talents: Presets.Talents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,

		// Default buffs and debuffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: Presets.DefaultPartyBuffs,
		individualBuffs: Presets.DefaultIndividualBuffs,
		debuffs: Presets.DefaultDebuffs,

		other: Presets.OtherDefaults,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [WarlockInputs.PetInput()],

	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	petConsumeInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.InputDelay, OtherInputs.DistanceFromTarget, OtherInputs.TankAssignment, OtherInputs.ChannelClipDelay],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.Talents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.WARLOCK_DEFAULT_APL],

		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PREBIS_GEAR, Presets.P1_BIS_GEAR],
		itemSwaps: [],
	},

	autoRotation: (_player: Player<Spec.SpecWarlock>): APLRotation => {
		return Presets.WARLOCK_DEFAULT_APL.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecWarlock,
			talents: Presets.Talents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.PREBIS_GEAR.gear,
				},
				[Faction.Horde]: {
					1: Presets.PREBIS_GEAR.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class WarlockSimUI extends IndividualSimUI<Spec.SpecWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarlock>) {
		super(parentElem, player, SPEC_CONFIG);

		// Migration guard: older saved Warlock settings may keep an empty scaffold APL,
		// which causes idle 0 DPS sims. Run after settings load and seed a default APL.
		this.sim.waitForInit().then(() => {
			const resolvedRotation = this.player.getResolvedAplRotation(true);
			const hasNoAplActions = (resolvedRotation.prepullActions?.length ?? 0) === 0 && (resolvedRotation.priorityList?.length ?? 0) === 0;
			if (hasNoAplActions) {
				this.player.setAplRotation(TypedEvent.nextEventID(), Presets.WARLOCK_DEFAULT_APL.rotation.rotation!);
			}
		});
	}
}
