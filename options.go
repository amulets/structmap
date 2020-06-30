package structmap

// WithDebug enable stack trace when have an panic
func WithDebug(sm *StructMap) {
	sm.debug = true
}

// WithBehaviors defines all logic of behaviors
func WithBehaviors(behaviors ...Behavior) OptionFunc {
	return func(sm *StructMap) {
		sm.behaviors = append(sm.behaviors, behaviors...)
	}
}
