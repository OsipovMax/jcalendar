package manager

import (
	"context"
)

type RuleManager struct{}

func NewRuleManager(_ context.Context) *RuleManager {
	return &RuleManager{}
}
