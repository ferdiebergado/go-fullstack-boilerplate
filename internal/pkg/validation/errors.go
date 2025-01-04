package validation

type errors map[string][]string

func (e errors) Add(field string, msg string) {
	e[field] = append(e[field], msg)
}

func (e errors) Get(field string) []string {
	errs := e[field]
	if len(errs) == 0 {
		return nil
	}
	return errs
}

type InputError struct {
	Errors errors `json:"errors"`
}

func (e *InputError) Error() string {
	return "Invalid input!"
}
