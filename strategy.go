package main

import "fmt"

// DefaultMergeStrategy for handling overwrite, skip, and skipped keys
type DefaultMergeStrategy struct {
	Overwrite bool
	Skip      bool
	SkipKeys  map[string]struct{} // Keys to skip
}

func (s *DefaultMergeStrategy) Merge(target, source map[string]any) error {
	for key, value := range source {
		// Skip keys explicitly listed in SkipKeys
		if _, shouldSkip := s.SkipKeys[key]; shouldSkip {
			if *verboseFlag {
				fmt.Printf("Skipping key: %s\n", key)
			}
			continue
		}

		if existingValue, exists := target[key]; exists {
			switch v := value.(type) {
			case map[string]any:
				if existingMap, ok := existingValue.(map[string]any); ok {
					// Recursive merge
					strategy := &DefaultMergeStrategy{
						Overwrite: s.Overwrite,
						Skip:      s.Skip,
						SkipKeys:  s.SkipKeys,
					}
					if err := strategy.Merge(existingMap, v); err != nil {
						return err
					}
					continue
				}
				return s.handleConflict(key, value, target)
			case []any:
				if existingArray, ok := existingValue.([]any); ok {
					// Merge arrays
					target[key] = append(existingArray, v...)
					continue
				}
				return s.handleConflict(key, value, target)
			default:
				return s.handleConflict(key, value, target)
			}
		} else {
			target[key] = value
		}
	}
	return nil
}

func (s *DefaultMergeStrategy) handleConflict(key string, value any, target map[string]any) error {
	if s.Overwrite {
		target[key] = value
		return nil
	}
	if s.Skip {
		return nil
	}
	return fmt.Errorf("conflicting types for key: '%s'", key)
}
