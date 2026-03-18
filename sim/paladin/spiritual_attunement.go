package paladin

import "github.com/wowsims/tbc/sim/core"

func (paladin *Paladin) registerSpiritualAttunement() {
	paladin.SpiritualAttunementMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 33776})
}

