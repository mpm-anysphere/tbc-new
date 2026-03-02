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
	if prebisDps < prebisShadowDps {
		t.Fatalf("expected prebis fire preset to be top pre-raid set, got Fire=%0.2f Shadow=%0.2f", prebisDps, prebisShadowDps)
	}

	p1Dps := runWarlockPresetDps(t, "p1")
	p1CraftedDps := runWarlockPresetDps(t, "p1_crafted")
	t.Logf("P1 DPS: T4Core=%0.2f, CraftedHybrid=%0.2f", p1Dps, p1CraftedDps)
	if p1CraftedDps < p1Dps {
		t.Fatalf("expected p1 crafted preset to be top p1 set, got T4Core=%0.2f CraftedHybrid=%0.2f", p1Dps, p1CraftedDps)
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

var DefaultOptions = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options: &proto.Warlock_Options{
			ClassOptions: &proto.WarlockOptions{
				Summon: proto.WarlockOptions_NoSummon,
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
