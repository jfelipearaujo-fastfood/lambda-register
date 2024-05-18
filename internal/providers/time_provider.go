package providers

import (
	"time"

	"github.com/jfelipearaujo-org/lambda-register/internal/providers/interfaces"
)

type TimeProvider struct {
	funcTime interfaces.FuncTime
}

func NewTimeProvider(funcTime interfaces.FuncTime) *TimeProvider {
	return &TimeProvider{
		funcTime: funcTime,
	}
}

func (p *TimeProvider) GetTime() time.Time {
	return p.funcTime()
}
