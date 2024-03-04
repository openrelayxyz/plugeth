package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
)

// Is1559 returns whether num is either equal to the London fork block or greater, if the chain supports EIP1559
func (c *ChainConfig) Is1559(num *big.Int) bool {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting is1559, but default PluginLoader has not been initialized")
		return c.IsLondon(num)
	}
	if active, ok := PluginEIPCheck(plugins.DefaultPluginLoader, "Is1559", num); ok {
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
// IsEIP160 returns whether num is either equal to the EIP160 block or greater. 
// This defaults to same as 158, but some chains do it at a different block
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

// IsShanghai is modified here to return whether num is either equal to the Shanghai fork block or greater, if the chain supports Shanghai
// the foundation implementation has been commented out
func (c *ChainConfig) IsShanghai(num *big.Int, time uint64) bool {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting isPluginShanghai, but default PluginLoader has not been initialized")
		return c.IsLondon(num) && isTimestampForked(c.ShanghaiTime, time)
	}
	if active, ok := PluginEIPCheck(plugins.DefaultPluginLoader, "IsShanghai", num); ok {
		return active
	}
	return c.IsLondon(num) && isTimestampForked(c.ShanghaiTime, time)
}