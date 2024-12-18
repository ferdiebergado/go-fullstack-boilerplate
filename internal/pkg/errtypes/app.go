package errtypes

type severity string

const (
	Low      severity = "Low"
	High     severity = "High"
	Critical severity = "Critical"
	Fatal    severity = "Fatal"
)

type AppError struct {
	Description string
	Err         error
	Severity    severity
}

func (a *AppError) Error() string {
	return a.Description
}
