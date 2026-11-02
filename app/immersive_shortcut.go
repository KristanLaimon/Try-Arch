package main

type immersiveShortcut struct{ held bool }

func (key *immersiveShortcut) route(enabled, focused, control, alt, rightAlt, down bool) (toggle, swallow bool) {
	if key.held {
		if !down {
			key.held = false
		}
		return false, true
	}
	if enabled && focused && control && alt && !rightAlt && down {
		key.held = true
		return true, true
	}
	return false, false
}
