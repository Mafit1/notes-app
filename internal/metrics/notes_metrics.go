package metrics

import "github.com/prometheus/client_golang/prometheus"

type notesMetrics struct {
	counter    prometheus.Counter
	errCounter prometheus.Counter
}

func NewNotesMetrics() NotesMetrics {
	m := &notesMetrics{
		counter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "notes_created_total",
			Help: "Amount of notes created in total",
		}),
		errCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "notes_created_errors",
			Help: "Amount of errors while creating notes",
		}),
	}
	prometheus.MustRegister(m.counter)
	prometheus.MustRegister(m.errCounter)

	return m
}

func (m *notesMetrics) IncCreated() {
	m.counter.Inc()
}

func (m *notesMetrics) IncCreatedError() {
	m.errCounter.Inc()
}
