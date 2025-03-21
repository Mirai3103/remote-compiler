package executor

import (
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"go.uber.org/zap"
)

var (
	IsolateStrategy = "isolate"
	RiskStrategy    = "risk"
)

func NewExecutor(logger *zap.Logger, cfg config.ExecutorConfig) Executor {
	strategy := cfg.Strategy
	if strategy == RiskStrategy {
		logger.Warn("Using risk executor - this is not recommended for production environments")
	}
	switch strategy {
	case IsolateStrategy:
		return newIsolateExecutor(logger, cfg)
	case RiskStrategy:
		return newRiskExecutor(logger, cfg)
	default:
		return newRiskExecutor(logger, cfg)
	}
}
