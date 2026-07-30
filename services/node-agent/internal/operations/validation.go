package operations

import "fmt"

type Validator interface {
	Validate(req Request) error
}

type validator struct {
	provider Provider
}

func NewValidator(provider Provider) Validator {
	return &validator{provider: provider}
}

func (v *validator) Validate(req Request) error {
	if !knownAction(req.Action) {
		return ValidationError{Message: fmt.Sprintf("unknown action: %s", req.Action)}
	}

	if !req.Confirm {
		return ValidationError{Message: "confirmation is required"}
	}

	if v.provider != nil && !v.provider.Supports(req.Action) {
		return ValidationError{Message: fmt.Sprintf("provider does not support action: %s", req.Action)}
	}

	return nil
}

func knownAction(action Action) bool {
	switch action {
	case ActionSleep, ActionRestart, ActionShutdown:
		return true
	default:
		return false
	}
}
