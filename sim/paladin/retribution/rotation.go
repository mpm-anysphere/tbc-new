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

	// Open with JoW for mana support, then proceed into SoB-based twisting.
	if !ret.JudgementOfWisdomAura.IsActive() {
		if ret.CurrentSeal != ret.SealOfWisdomAura {
			if ret.SealOfWisdom.CanCast(sim, target) {
				ret.SealOfWisdom.Cast(sim, target)
				return
			}
		} else if ret.CanJudgementOfWisdom(sim, target) {
			ret.JudgementOfWisdom.Cast(sim, target)
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
	if !ret.GCD.IsReady(sim) {
		ret.WaitUntil(sim, ret.GCD.ReadyAt())
		return
	}

	nextSwingAt := ret.AutoAttacks.NextAttackAt()
	timeTilNextSwing := nextSwingAt - sim.CurrentTime
	spellGCD := ret.SpellGCD()

	// Complete an active twist by restoring SoB within the twist window.
	if ret.SealOfCommandAura.IsActive() {
		if timeTilNextSwing <= paladin.TwistWindow+time.Millisecond*50 &&
			ret.SealOfBlood.CanCast(sim, target) {
			ret.SealOfBlood.Cast(sim, target)
			return
		}

		if timeTilNextSwing > paladin.TwistWindow {
			ret.WaitUntil(sim, nextSwingAt-paladin.TwistWindow)
		} else {
			ret.WaitUntil(sim, nextSwingAt)
		}
		return
	}

	// Start a new twist cycle by swapping to SoC shortly before the next swing.
	if ret.SealOfBloodAura.IsActive() &&
		timeTilNextSwing > spellGCD &&
		timeTilNextSwing <= spellGCD+paladin.TwistWindow &&
		ret.SealOfCommand.CanCast(sim, target) {
		ret.SealOfCommand.Cast(sim, target)
		return
	}

	if ret.CurrentSeal == nil || !ret.CurrentSeal.IsActive() {
		if ret.SealOfBlood.CanCast(sim, target) {
			ret.SealOfBlood.Cast(sim, target)
			return
		}
	}

	if ret.CanJudgementOfBlood(sim, target) {
		ret.JudgementOfBlood.Cast(sim, target)
		return
	}

	if ret.CrusaderStrike.CanCast(sim, target) {
		ret.CrusaderStrike.Cast(sim, target)
		return
	}

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
		return
	}

	nextEvent := minAtLeast(
		sim.CurrentTime+time.Millisecond*100,
		ret.NextGCDAt(),
		ret.CrusaderStrike.CD.ReadyAt(),
		ret.JudgementOfBlood.CD.ReadyAt(),
		nextSwingAt-paladin.TwistWindow,
		nextSwingAt,
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

