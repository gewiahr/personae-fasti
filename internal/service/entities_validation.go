package service

import (
	"strings"

	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/validation"
)

const (
	entityNameMaxLength        = 100
	entityTitleMaxLength       = 200
	entityDescriptionMaxLength = 20_000
)

func normalizeAndValidateEntity(name, title *string, description string) error {
	*name = strings.TrimSpace(*name)
	*title = strings.TrimSpace(*title)

	fields := make(map[string]string)
	if validation.IsBlank(*name) {
		fields["name"] = "Введите название"
	} else if validation.CharacterCount(*name) > entityNameMaxLength {
		fields["name"] = "Название не может быть длиннее 100 символов"
	}
	if validation.CharacterCount(*title) > entityTitleMaxLength {
		fields["title"] = "Заглавие не может быть длиннее 200 символов"
	}
	if validation.CharacterCount(description) > entityDescriptionMaxLength {
		fields["description"] = "Описание не может быть длиннее 20000 символов"
	}

	if len(fields) > 0 {
		return e.NewFieldValidationError(fields)
	}
	return nil
}
