package turingpi

type ResultError struct {
	Reason string
}

func (m *ResultError) Error() string {
	return m.Reason
}
