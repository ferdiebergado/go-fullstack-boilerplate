package html

import "fmt"

type TemplateNotFoundError struct {
	Template string
}

func (t *TemplateNotFoundError) Error() string {
	return fmt.Sprintf("template: %s not found.", t.Template)
}
