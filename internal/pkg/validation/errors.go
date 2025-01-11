package validation

type Errors map[string][]string

type Error struct {
	Errors Errors `json:"errors"`
}

func NewError() *Error {
	return &Error{
		Errors: make(Errors),
	}
}

func (e *Error) Add(field string, msg string) {
	e.Errors[field] = append(e.Errors[field], msg)
}

func (e *Error) Get(field string) []string {
	errs := e.Errors[field]
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (e *Error) Count() int {
	return len(e.Errors)
}

func (e *Error) Error() string {
	return "Invalid input!"
}
