//go:build !android && !ios

package system

func UseTouchControls() bool {
	return false
}
