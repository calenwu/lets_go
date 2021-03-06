package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form {
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields []string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MinLength(fields []string, d int) {
	for _, field := range fields {
		value := f.Get(field)
		if value == "" {
			return
		}
		if utf8.RuneCountInString(value) < d {
			f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d)", d))
		}
	}
}

func (f *Form) MatchesPattern(fields []string, pattern *regexp.Regexp) {
	for _, field := range fields {
		value := f.Get(field)
		if !pattern.MatchString(value) {
			f.Errors.Add(field, "This field is invalid")
		}
	}
}

func (f *Form) IsEmail(fields []string) {
	for _, field := range fields {
		_, err := mail.ParseAddress(f.Get(field))
		if err != nil {
			f.Errors.Add(field, "This field is not an email")
		}
	}
}

func (f *Form) MaxLength(fields []string, d int) {
	for _, field := range fields {
		value := f.Get(field)
		if utf8.RuneCountInString(value) > d {
			f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d)", d))
		}
	}
}

func (f *Form) PermittedValues(fields []string, opts ...string) {
	for _, field := range fields {
		value := f.Get(field)
		for _, opt := range opts {
			if value == opt {
				return
			}
		}
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}