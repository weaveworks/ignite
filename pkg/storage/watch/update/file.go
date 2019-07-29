package update

// FileUpdate is used by watchers to
// signal the state change of a file.
type FileUpdate struct {
	Event Event
	Path  string
}
