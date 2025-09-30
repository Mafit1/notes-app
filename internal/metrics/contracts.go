package metrics

//go:generate mockgen -destination=mocks/mock_$GOFILE -package=mocks . NotesMetrics
type NotesMetrics interface {
	IncCreated()
	IncCreatedError()
}
