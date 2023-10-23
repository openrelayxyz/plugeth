package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
)

// IsLondon returns whether num is either equal to the London fork block or greater.
func (c *ChainConfig) Is1559(num *big.Int) bool {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting is1559, but default PluginLoader has not been initialized")
		return c.IsLondon(num)
	}
	if active, ok := PluginEIPCheck(plugins.DefaultPluginLoader, "Is1559" num); ok {
		return active
	}
	return c.IsLondon(num)
}


func PluginEIPCheck(pl *plugins.PluginLoader, eipHookName string, num *big.Int) (bool, bool) {
	fn, ok := plugins.LookupOne[func(*big.Int) bool](pl, eipHookName)
	if !ok {
		return false, false
	}
	return fn(num), ok
}
// IsLondon returns whether num is either equal to the London fork block or greater.
func (c *ChainConfig) IsEIP160(num *big.Int) bool {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting is160, but default PluginLoader has not been initialized")
		return c.IsEIP158(num)
	}
	if active, ok := PluginEIPCheck(plugins.DefaultPluginLoader, "Is160", num); ok {
		return active
	}
	return c.IsEIP158(num)
}
