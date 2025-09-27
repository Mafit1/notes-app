package metrics

type NotesMetrics interface {
	IncCreated()
	IncCreatedError()
}
