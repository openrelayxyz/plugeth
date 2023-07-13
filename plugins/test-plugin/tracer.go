package main 

import (
	"math/big"
	"time"

	"github.com/openrelayxyz/plugeth-utils/core"
)

 type TracerService struct {}

var Tracers = map[string]func(core.StateDB) core.TracerResult{
    "testTracer": func(core.StateDB) core.TracerResult {
        return &TracerService{}
    },
}

func (b *TracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	m := map[string]struct{}{
		"StandardCaptureStart": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	m := map[string]struct{}{
		"StandardCaptureState": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
	m := map[string]struct{}{
		"StandardCaptureFault": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	m := map[string]struct{}{
		"StandardCaptureEnd": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	m := map[string]struct{}{
		"StandardCaptureEnter": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureExit(output []byte, gasUsed uint64, err error) {
	m := map[string]struct{}{
		"StandardCaptureExit": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) Result() (interface{}, error) { 
	m := map[string]struct{}{
		"StandardTracerResult": struct{}{},
	}
	hookChan <- m
	return "", nil }