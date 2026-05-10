package validator

import (
	"fmt"
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	msgs := []string{}
	for _, err := range e {
		msgs = append(msgs, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(msgs, ", ")
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateRegister(email, password, name string) ValidationErrors {
	var errs ValidationErrors

	// Email
	if strings.TrimSpace(email) == "" {
		errs = append(errs, ValidationError{Field: "email", Message: "email is required"})
	} else if !emailRegex.MatchString(email) {
		errs = append(errs, ValidationError{Field: "email", Message: "invalid email format"})
	}

	// Password
	if strings.TrimSpace(password) == "" {
		errs = append(errs, ValidationError{Field: "password", Message: "password is required"})
	} else if len(password) < 8 {
		errs = append(errs, ValidationError{Field: "password", Message: "password must be at least 8 characters"})
	}

	// Name
	if strings.TrimSpace(name) == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "name is required"})
	} else if len(name) < 2 {
		errs = append(errs, ValidationError{Field: "name", Message: "name must be at least 2 characters"})
	}

	return errs
}

func ValidateLogin(email, password string) ValidationErrors {
	var errs ValidationErrors

	if strings.TrimSpace(email) == "" {
		errs = append(errs, ValidationError{Field: "email", Message: "email is required"})
	}

	if strings.TrimSpace(password) == "" {
		errs = append(errs, ValidationError{Field: "password", Message: "password is required"})
	}

	return errs
}

func ValidateCreateTask(title string) ValidationErrors {
	var errs ValidationErrors

	if strings.TrimSpace(title) == "" {
		errs = append(errs, ValidationError{Field: "title", Message: "title is required"})
	} else if len(title) < 3 {
		errs = append(errs, ValidationError{Field: "title", Message: "title must be at least 3 characters"})
	} else if len(title) > 255 {
		errs = append(errs, ValidationError{Field: "title", Message: "title must be less than 255 characters"})
	}

	return errs
}

func ValidateUpdateTask(status *string) ValidationErrors {
	var errs ValidationErrors

	if status != nil {
		validStatuses := map[string]bool{
			"todo":        true,
			"in_progress": true,
			"done":        true,
		}
		if !validStatuses[*status] {
			errs = append(errs, ValidationError{Field: "status", Message: "status must be todo, in_progress or done"})
		}
	}

	return errs
}