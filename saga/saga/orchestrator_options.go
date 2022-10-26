package saga

import "github.com/tikivn/s14e-backend-utils/logger"

// OrchestratorOption options for Orchestrator
type OrchestratorOption func(o *Orchestrator)

// WithOrchestratorLogger is an option to set the logger.Logger of the Orchestrator
func WithOrchestratorLogger(logger logger.Logger) OrchestratorOption {
	return func(o *Orchestrator) {
		o.logger = logger
	}
}
