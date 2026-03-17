package retribution

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/paladin"
)

func (ret *RetributionPaladin) ExecuteCustomRotation(sim *core.Simulation) {
	if ret.CurrentTarget == nil {
		return
	}

	if !ret.openerCompleted {
		ret.openingRotation(sim)
		return
	}

	ret.mainRotation(sim)
}

func (ret *RetributionPaladin) openingRotation(sim *core.Simulation) {
	target := ret.CurrentTarget

	// Open with Judgement of the Crusader for throughput parity.
	if !ret.JudgementOfTheCrusaderAura.IsActive() {
		if ret.CurrentSeal != ret.SealOfTheCrusaderAura {
			if ret.SealOfTheCrusader.CanCast(sim, target) {
				ret.SealOfTheCrusader.Cast(sim, target)
				return
			}
		} else if ret.CanJudgementOfTheCrusader(sim, target) {
			ret.JudgementOfTheCrusader.Cast(sim, target)
			return
		}
	}

	// Cast Seal of Command first, then Seal of Blood to begin twist cadence.
	if !ret.SealOfCommandAura.IsActive() {
		if ret.SealOfCommand.CanCast(sim, target) {
			ret.SealOfCommand.Cast(sim, target)
			return
		}
	}

	if !ret.SealOfBloodAura.IsActive() {
		if ret.SealOfBlood.CanCast(sim, target) {
			ret.SealOfBlood.Cast(sim, target)
			return
		}
	}

	ret.AutoAttacks.EnableAutoSwing(sim)
	ret.openerCompleted = true
}

func (ret *RetributionPaladin) mainRotation(sim *core.Simulation) {
	target := ret.CurrentTarget
	socActive := ret.SealOfCommandAura.IsActive()
	if ret.CurrentMana() <= 1000 && !socActive {
		ret.lowManaRotation(sim)
		return
	}

	gcdCD := ret.GCD.TimeToReady(sim)
	crusaderStrikeCD := ret.CrusaderStrike.TimeToReady(sim)
	nextCrusaderStrikeAt := ret.CrusaderStrike.CD.ReadyAt()
	judgementCD := ret.JudgementOfBlood.TimeToReady(sim)
	nextJudgementAt := ret.JudgementOfBlood.CD.ReadyAt()

	nextSwingAt := ret.AutoAttacks.NextAttackAt()
	timeTilNextSwing := nextSwingAt - sim.CurrentTime
	spellGCD := ret.SpellGCD()

	sobActive := ret.SealOfBloodAura.IsActive()
	inTwistWindow := sim.CurrentTime >= nextSwingAt-paladin.TwistWindow && sim.CurrentTime < nextSwingAt
	latestTwistStart := nextSwingAt - spellGCD
	possibleTwist := timeTilNextSwing > spellGCD+gcdCD
	willTwist := possibleTwist && (nextSwingAt+spellGCD <= nextCrusaderStrikeAt)

	if judgementCD == 0 && sobActive && willTwist {
		ret.JudgementOfBlood.Cast(sim, target)
		sobActive = false
	}

	if gcdCD == 0 {
		if socActive && inTwistWindow {
			if ret.SealOfBlood.CanCast(sim, target) {
				ret.SealOfBlood.Cast(sim, target)
				return
			}
		} else if crusaderStrikeCD == 0 && !willTwist && (sobActive || spellGCD < timeTilNextSwing) {
			ret.CrusaderStrike.Cast(sim, target)
			return
		} else if willTwist && !socActive && nextJudgementAt > latestTwistStart {
			if ret.SealOfCommand.CanCast(sim, target) {
				ret.SealOfCommand.Cast(sim, target)
				return
			}
		} else if !sobActive && !socActive && !willTwist {
			if ret.SealOfBlood.CanCast(sim, target) {
				ret.SealOfBlood.Cast(sim, target)
				return
			}
		} else if !willTwist && !socActive &&
			timeTilNextSwing+ret.AutoAttacks.MainhandSwingSpeed() > spellGCD*2 &&
			spellGCD < crusaderStrikeCD {
			ret.useFillers(sim, target)
			return
		}
	}

	nextEvent := minAtLeast(
		sim.CurrentTime+time.Millisecond*50,
		nextSwingAt,
		latestTwistStart,
		ret.GCD.ReadyAt(),
		nextJudgementAt,
		nextCrusaderStrikeAt,
	)
	ret.WaitUntil(sim, nextEvent)
}

func (ret *RetributionPaladin) useFillers(sim *core.Simulation, target *core.Unit) {
	if ret.Exorcism.CanCast(sim, target) &&
		ret.CanExorcism(target) &&
		ret.CurrentManaPercent() > 0.4 {
		ret.Exorcism.Cast(sim, target)
		return
	}

	if ret.Consecration != nil &&
		ret.Consecration.CanCast(sim, target) &&
		ret.CurrentManaPercent() > 0.6 {
		ret.Consecration.Cast(sim, target)
	}
}

func (ret *RetributionPaladin) lowManaRotation(sim *core.Simulation) {
	target := ret.CurrentTarget
	sobExpiration := ret.SealOfBloodAura.ExpiresAt()
	nextSwingAt := ret.AutoAttacks.NextAttackAt()

	manaRegenAt := core.NeverExpires
	if sim.CurrentTime+time.Second >= sobExpiration {
		sobAndJudgeCost := ret.SealOfBlood.Cost.GetCurrentCost() + ret.JudgementOfBlood.Cost.GetCurrentCost()
		if ret.CanJudgementOfBlood(sim, target) && ret.CurrentMana() >= sobAndJudgeCost {
			ret.JudgementOfBlood.Cast(sim, target)
		}

		if ret.GCD.IsReady(sim) {
			if !ret.SealOfBlood.Cast(sim, target) {
				manaRegenAt = sim.CurrentTime + ret.TimeUntilManaRegen(ret.SealOfBlood.CurCast.Cost)
			}
		}
	} else if ret.GCD.IsReady(sim) && ret.CrusaderStrike.CD.IsReady(sim) {
		spellGCD := ret.SpellGCD()
		sobAndCSCost := ret.SealOfBlood.Cost.GetCurrentCost() + ret.CrusaderStrike.Cost.GetCurrentCost()

		if !(spellGCD+sim.CurrentTime > nextSwingAt && sobExpiration < nextSwingAt) &&
			(ret.CurrentMana() >= sobAndCSCost) {
			ret.CrusaderStrike.Cast(sim, target)
		}
	}

	nextEvent := minAtLeast(
		sim.CurrentTime+time.Millisecond*100,
		ret.GCD.ReadyAt(),
		ret.CrusaderStrike.CD.ReadyAt(),
		manaRegenAt,
		sobExpiration-time.Second,
	)
	ret.WaitUntil(sim, nextEvent)
}

func minAtLeast(base time.Duration, values ...time.Duration) time.Duration {
	next := time.Duration(1<<63 - 1)
	for _, v := range values {
		if v > base && v < next {
			next = v
		}
	}
	if next == time.Duration(1<<63-1) {
		return base
	}
	return next
}

