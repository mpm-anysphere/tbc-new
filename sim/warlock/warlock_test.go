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
	player := core.WithSpec(
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
	raid := core.SinglePlayerRaidProto(player, nil, nil, nil)
	p1Player := googleProto.Clone(player).(*proto.Player)
	p1Player.Equipment = core.GetGearSet("../../ui/warlock/dps/gear_sets", "p1").GearSet
	p1Raid := core.SinglePlayerRaidProto(p1Player, nil, nil, nil)

	core.RunTestSuite(t, t.Name(), []core.TestGenerator{
		&core.SingleCharacterStatsTestGenerator{
			Name: "CharacterStats",
			Request: &proto.ComputeStatsRequest{
				Raid: raid,
			},
		},
		&core.SingleDpsTestGenerator{
			Name: "AverageDps",
			Request: &proto.RaidSimRequest{
				Raid:       googleProto.Clone(raid).(*proto.Raid),
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
	})
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

var DefaultTalents = "000000000000000000000-0000000000000000000000-000000000000000000000"

var DefaultConsumes = &proto.ConsumesSpec{
	GuardianElixirId: 32067, // Elixir of Draenic Wisdom
	BattleElixirId:   28103, // Adept's Elixir
	FoodId:           27657, // Blackened Basilisk
	MhImbueId:        25122, // Brilliant Wizard Oil
	PotId:            22839, // Destruction Potion
	ConjuredId:       12662, // Demonic Rune
	DrumsId:          351355,
}
