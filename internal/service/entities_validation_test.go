package service

import (
	"errors"
	"strings"
	"testing"

	e "personae-fasti/internal/pkg/errorutils"
)

func TestNormalizeAndValidateEntity(t *testing.T) {
	tests := []struct {
		name        string
		entityName  string
		title       string
		description string
		wantField   string
	}{
		{name: "valid", entityName: "  Элли  ", title: "  Следопыт  "},
		{name: "blank name", entityName: " \n\t ", wantField: "name"},
		{name: "long unicode name", entityName: strings.Repeat("я", entityNameMaxLength+1), wantField: "name"},
		{name: "long title", entityName: "Элли", title: strings.Repeat("я", entityTitleMaxLength+1), wantField: "title"},
		{name: "long description", entityName: "Элли", description: strings.Repeat("я", entityDescriptionMaxLength+1), wantField: "description"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, title := test.entityName, test.title
			err := normalizeAndValidateEntity(&name, &title, test.description)
			if test.wantField == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if name != "Элли" || title != "Следопыт" {
					t.Fatalf("plain fields were not trimmed: name=%q title=%q", name, title)
				}
				return
			}

			var appErr *e.AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected AppError, got %T: %v", err, err)
			}
			if appErr.Fields[test.wantField] == "" {
				t.Fatalf("expected validation error for %q, got %#v", test.wantField, appErr.Fields)
			}
		})
	}
}
