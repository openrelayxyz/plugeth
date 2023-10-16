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
	if active, ok := PluginIs1559(plugins.DefaultPluginLoader, num); ok {
		return active
	}
	return c.IsLondon(num)
}


func PluginIs1559(pl *plugins.PluginLoader, num *big.Int) (bool, bool) {
	fn, ok := plugins.LookupOne[func(*big.Int) bool](pl, "Is1559")
	if !ok {
		return false, false
	}
	return fn(num), ok
}
