package retribution

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	_ "github.com/wowsims/tbc/sim/common" // imported to get item effects included.
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func TestRetribution(t *testing.T) {
	if !core.WITH_DB {
		t.Skip("requires -tags=with_db")
	}

	ensureRetributionRegistered()

	gearSets := []string{"preraid", "p1"}
	races := []proto.Race{proto.Race_RaceBloodElf, proto.Race_RaceHuman, proto.Race_RaceDraenei}
	rotation := core.GetAplRotation("../../../ui/paladin/retribution/apls", "default").Rotation

	for _, gearSet := range gearSets {
		gear := core.GetGearSet("../../../ui/paladin/retribution/gear_sets", gearSet).GearSet
		for _, race := range races {
			label := fmt.Sprintf("%s_%s", gearSet, race.String())
			t.Run(label, func(t *testing.T) {
				player := core.WithSpec(&proto.Player{
					Class:              proto.Class_ClassPaladin,
					Race:               race,
					Equipment:          gear,
					TalentsString:      DefaultTalents,
					Consumables:        DefaultConsumables,
					Buffs:              core.FullIndividualBuffs,
					Rotation:           rotation,
					DistanceFromTarget: 5,
					Profession1:        proto.Profession_Engineering,
					Profession2:        proto.Profession_Blacksmithing,
				}, DefaultOptions)

				request := &proto.RaidSimRequest{
					Raid:       core.SinglePlayerRaidProto(player, core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs),
					Encounter:  core.MakeSingleTargetEncounter(180),
					SimOptions: core.DefaultSimTestOptions,
				}

				result := core.RunRaidSim(request)
				if result.GetError() != nil && result.GetError().GetMessage() != "" {
					t.Fatalf("raid sim failed: %s", result.GetError().GetMessage())
				}
				if result.RaidMetrics.Dps.Avg <= 0 {
					t.Fatalf("expected positive raid DPS, got %.2f", result.RaidMetrics.Dps.Avg)
				}
			})
		}
	}
}

func TestRetributionSealTwisting(t *testing.T) {
	if !core.WITH_DB {
		t.Skip("requires -tags=with_db")
	}

	ensureRetributionRegistered()

	gear := core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "preraid").GearSet
	rotation := core.GetAplRotation("../../../ui/paladin/retribution/apls", "default").Rotation
	player := core.WithSpec(&proto.Player{
		Class:              proto.Class_ClassPaladin,
		Race:               proto.Race_RaceBloodElf,
		Equipment:          gear,
		TalentsString:      DefaultTalents,
		Consumables:        DefaultConsumables,
		Buffs:              core.FullIndividualBuffs,
		Rotation:           rotation,
		DistanceFromTarget: 5,
		Profession1:        proto.Profession_Engineering,
		Profession2:        proto.Profession_Blacksmithing,
	}, DefaultOptions)

	request := &proto.RaidSimRequest{
		Raid:      core.SinglePlayerRaidProto(player, core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs),
		Encounter: core.MakeSingleTargetEncounter(60),
		SimOptions: &proto.SimOptions{
			Iterations: 1,
			IsTest:     true,
			Debug:      true,
			RandomSeed: 101,
		},
	}

	result := core.RunRaidSim(request)
	if result.GetError() != nil && result.GetError().GetMessage() != "" {
		t.Fatalf("raid sim failed: %s", result.GetError().GetMessage())
	}

	logs := result.GetLogs()
	if !strings.Contains(logs, "{SpellID: 20375}") {
		t.Fatalf("expected Seal of Command casts in logs for twisting, logs:\n%s", logs)
	}
	if strings.Count(logs, "{SpellID: 31892}") < 2 {
		t.Fatalf("expected repeated Seal of Blood casts in logs for twisting, logs:\n%s", logs)
	}
}

var registerRetributionOnce sync.Once

func ensureRetributionRegistered() {
	registerRetributionOnce.Do(RegisterRetributionPaladin)
}

var DefaultOptions = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: &proto.RetributionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal: proto.PaladinSeal_Truth,
			},
		},
	},
}

var DefaultTalents = ""

var DefaultConsumables = &proto.ConsumesSpec{
	FlaskId:      22854,  // Flask of Relentless Assault
	FoodId:       27658,  // Roasted Clefthoof
	PotId:        22838,  // Haste Potion
	ConjuredId:   12662,  // Demonic Rune
	SuperSapper:  true,
	GoblinSapper: true,
	DrumsId:      351355, // Greater Drums of Battle
}

