package warlock

import (
	"testing"

	_ "github.com/wowsims/tbc/sim/common" // imported to get item effects included.
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

func init() {
	RegisterWarlock()
}

func TestWarlock(t *testing.T) {
	prebisPlayer := core.WithSpec(
		&proto.Player{
			Class:              proto.Class_ClassWarlock,
			Race:               proto.Race_RaceGnome,
			Equipment:          core.GetGearSet("../../ui/warlock/dps/gear_sets", "prebis").GearSet,
			TalentsString:      DefaultTalents,
			Rotation:           core.GetAplRotation("../../ui/warlock/dps/apls", "default").Rotation,
			Consumables:        DefaultConsumes,
			DistanceFromTarget: 20,
		},
		DefaultOptions,
	)
	prebisRaid := core.SinglePlayerRaidProto(prebisPlayer, nil, nil, nil)

	prebisShadowPlayer := googleProto.Clone(prebisPlayer).(*proto.Player)
	prebisShadowPlayer.Equipment = core.GetGearSet("../../ui/warlock/dps/gear_sets", "prebis_shadow").GearSet
	prebisShadowRaid := core.SinglePlayerRaidProto(prebisShadowPlayer, nil, nil, nil)

	p1Player := googleProto.Clone(prebisPlayer).(*proto.Player)
	p1Player.Equipment = core.GetGearSet("../../ui/warlock/dps/gear_sets", "p1").GearSet
	p1Raid := core.SinglePlayerRaidProto(p1Player, nil, nil, nil)

	p1CraftedPlayer := googleProto.Clone(prebisPlayer).(*proto.Player)
	p1CraftedPlayer.Equipment = core.GetGearSet("../../ui/warlock/dps/gear_sets", "p1_crafted").GearSet
	p1CraftedRaid := core.SinglePlayerRaidProto(p1CraftedPlayer, nil, nil, nil)

	core.RunTestSuite(t, t.Name(), []core.TestGenerator{
		&core.SingleCharacterStatsTestGenerator{
			Name: "CharacterStats",
			Request: &proto.ComputeStatsRequest{
				Raid: prebisRaid,
			},
		},
		&core.SingleDpsTestGenerator{
			Name: "AverageDps",
			Request: &proto.RaidSimRequest{
				Raid:       googleProto.Clone(prebisRaid).(*proto.Raid),
				Encounter:  core.MakeSingleTargetEncounter(20),
				SimOptions: core.AverageDefaultSimTestOptions,
			},
		},
		&core.SingleDpsTestGenerator{
			Name: "AverageDpsPrebisShadow",
			Request: &proto.RaidSimRequest{
				Raid:       googleProto.Clone(prebisShadowRaid).(*proto.Raid),
				Encounter:  core.MakeSingleTargetEncounter(20),
				SimOptions: core.AverageDefaultSimTestOptions,
			},
		},
		&core.SingleDpsTestGenerator{
			Name: "AverageDpsP1",
			Request: &proto.RaidSimRequest{
				Raid:       googleProto.Clone(p1Raid).(*proto.Raid),
				Encounter:  core.MakeSingleTargetEncounter(20),
				SimOptions: core.AverageDefaultSimTestOptions,
			},
		},
		&core.SingleDpsTestGenerator{
			Name: "AverageDpsP1Crafted",
			Request: &proto.RaidSimRequest{
				Raid:       googleProto.Clone(p1CraftedRaid).(*proto.Raid),
				Encounter:  core.MakeSingleTargetEncounter(20),
				SimOptions: core.AverageDefaultSimTestOptions,
			},
		},
	})
}

func TestWarlockPresetGearRanking(t *testing.T) {
	prebisDps := runWarlockPresetDps(t, "prebis")
	prebisShadowDps := runWarlockPresetDps(t, "prebis_shadow")
	t.Logf("Pre-Raid DPS: Fire=%0.2f, Shadow=%0.2f", prebisDps, prebisShadowDps)
	if prebisDps <= 0 || prebisShadowDps <= 0 {
		t.Fatalf("expected positive pre-raid DPS values, got Fire=%0.2f Shadow=%0.2f", prebisDps, prebisShadowDps)
	}

	p1Dps := runWarlockPresetDps(t, "p1")
	p1CraftedDps := runWarlockPresetDps(t, "p1_crafted")
	t.Logf("P1 DPS: T4Core=%0.2f, CraftedHybrid=%0.2f", p1Dps, p1CraftedDps)
	if p1Dps <= 0 || p1CraftedDps <= 0 {
		t.Fatalf("expected positive p1 DPS values, got T4Core=%0.2f CraftedHybrid=%0.2f", p1Dps, p1CraftedDps)
	}

	bestPrebis := max(prebisDps, prebisShadowDps)
	bestP1 := max(p1Dps, p1CraftedDps)
	if bestP1 < bestPrebis {
		t.Fatalf("expected best p1 preset to outperform best pre-raid preset, got bestP1=%0.2f bestPrebis=%0.2f", bestP1, bestPrebis)
	}
}

func TestWarlockRotationCastsShadowBoltFiller(t *testing.T) {
	result := runWarlockRaidSim(t, DefaultTalents, "default", "prebis", proto.WarlockOptions_Succubus)

	shadowBoltCasts := getActionCasts(result, 27209)
	if shadowBoltCasts == 0 {
		t.Fatalf("expected Shadow Bolt casts in default DS/Ruin rotation, got 0")
	}

	felFlameCasts := getActionCasts(result, 77799)
	if felFlameCasts > 0 {
		t.Fatalf("expected no Fel Flame casts in TBC default rotation, got %d", felFlameCasts)
	}
}

func TestWarlockAfflictionCastProfile(t *testing.T) {
	result := runWarlockRaidSim(t, "350020251023310550031--0550005122", "affliction", "prebis_shadow", proto.WarlockOptions_Felhunter)

	if casts := getActionCasts(result, 30405); casts == 0 {
		t.Fatalf("expected Unstable Affliction casts for UA profile, got 0")
	}
	if casts := getActionCasts(result, 27218); casts == 0 {
		t.Fatalf("expected Curse of Agony casts for affliction profile, got 0")
	}
	if casts := getActionCasts(result, 27209); casts == 0 {
		t.Fatalf("expected Shadow Bolt filler casts for affliction profile, got 0")
	}
}

func TestWarlockPetSummonSelection(t *testing.T) {
	resultWithImp := runWarlockRaidSim(t, "350020251023310550031--0550005122", "affliction", "prebis_shadow", proto.WarlockOptions_Imp)
	if casts := getPetActionCasts(resultWithImp, 27267); casts == 0 {
		t.Fatalf("expected Imp Firebolt casts when Imp summon selected, got 0")
	}

	resultWithDS := runWarlockRaidSim(t, DefaultTalents, "default", "prebis", proto.WarlockOptions_Succubus)
	if casts := getPetActionCasts(resultWithDS, 27274); casts > 0 {
		t.Fatalf("expected no Succubus casts for Demonic Sacrifice profile, got %d", casts)
	}
}

func runWarlockPresetDps(t *testing.T, gearSetName string) float64 {
	t.Helper()

	player := core.WithSpec(
		&proto.Player{
			Class:              proto.Class_ClassWarlock,
			Race:               proto.Race_RaceGnome,
			Equipment:          core.GetGearSet("../../ui/warlock/dps/gear_sets", gearSetName).GearSet,
			TalentsString:      DefaultTalents,
			Rotation:           core.GetAplRotation("../../ui/warlock/dps/apls", "default").Rotation,
			Consumables:        DefaultConsumes,
			DistanceFromTarget: 20,
		},
		DefaultOptions,
	)
	raid := core.SinglePlayerRaidProto(player, nil, nil, nil)

	result := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:      raid,
		Encounter: core.MakeSingleTargetEncounter(20),
		SimOptions: &proto.SimOptions{
			Iterations: 5000,
			IsTest:     true,
			Debug:      false,
			RandomSeed: 101,
		},
	})
	if result.Error != nil {
		t.Fatalf("warlock ranking sim failed for %s: %s", gearSetName, result.Error.Message)
	}

	return result.RaidMetrics.Dps.Avg
}

func runWarlockRaidSim(t *testing.T, talents string, aplName string, gearSet string, summon proto.WarlockOptions_Summon) *proto.RaidSimResult {
	t.Helper()

	options := &proto.Player_Warlock{
		Warlock: &proto.Warlock{
			Options: &proto.Warlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon: summon,
				},
			},
		},
	}

	player := core.WithSpec(
		&proto.Player{
			Class:              proto.Class_ClassWarlock,
			Race:               proto.Race_RaceGnome,
			Equipment:          core.GetGearSet("../../ui/warlock/dps/gear_sets", gearSet).GearSet,
			TalentsString:      talents,
			Rotation:           core.GetAplRotation("../../ui/warlock/dps/apls", aplName).Rotation,
			Consumables:        DefaultConsumes,
			DistanceFromTarget: 20,
		},
		options,
	)
	raid := core.SinglePlayerRaidProto(player, nil, nil, nil)

	result := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:      raid,
		Encounter: core.MakeSingleTargetEncounter(20),
		SimOptions: &proto.SimOptions{
			Iterations: 1000,
			IsTest:     true,
			Debug:      false,
			RandomSeed: 42,
		},
	})
	if result.Error != nil {
		t.Fatalf("warlock sim failed: %s", result.Error.Message)
	}

	return result
}

func getActionCasts(result *proto.RaidSimResult, spellID int32) int32 {
	player := result.RaidMetrics.Parties[0].Players[0]
	for _, action := range player.Actions {
		if action.GetId().GetSpellId() != spellID {
			continue
		}
		var casts int32
		for _, target := range action.Targets {
			casts += target.Casts
		}
		return casts
	}
	return 0
}

func getPetActionCasts(result *proto.RaidSimResult, spellID int32) int32 {
	player := result.RaidMetrics.Parties[0].Players[0]
	for _, pet := range player.Pets {
		for _, action := range pet.Actions {
			if action.GetId().GetSpellId() != spellID {
				continue
			}
			var casts int32
			for _, target := range action.Targets {
				casts += target.Casts
			}
			return casts
		}
	}
	return 0
}

var DefaultOptions = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options: &proto.Warlock_Options{
			ClassOptions: &proto.WarlockOptions{
				Summon: proto.WarlockOptions_Succubus,
			},
		},
	},
}

var DefaultTalents = "-20500301332101-55500050220001053025"

var DefaultConsumes = &proto.ConsumesSpec{
	FlaskId:    22866, // Flask of Pure Death
	FoodId:     27657, // Blackened Basilisk
	MhImbueId:  25122, // Brilliant Wizard Oil
	PotId:      22839, // Destruction Potion
	ConjuredId: 12662, // Demonic Rune
	DrumsId:    351355,
}
