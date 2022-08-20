package app

import (
	"github.com/passwall/passwall-server/internal/storage"
	"github.com/passwall/passwall-server/model"
	"github.com/passwall/passwall-server/pkg/logger"
)

// FindAllEmails finds all logins
func FindAllEmails(s storage.Store, schema string) ([]model.Email, error) {
	list, err := s.Emails().All(schema)
	if err != nil {
		return nil, err
	}

	// Decrypt server side encrypted fields
	for i := range list {
		m, err := DecryptModel(&list[i])
		if err != nil {
			logger.Errorf("Error while decrypting credit card: %v", err)
			continue
		}
		list[i] = *m.(*model.Email)
	}

	return list, nil
}

// CreateEmail creates a new bank account and saves it to the store
func CreateEmail(s storage.Store, dto *model.EmailDTO, schema string) (*model.Email, error) {
	rawModel := model.ToEmail(dto)
	encModel := EncryptModel(rawModel)

	createdEmail, err := s.Emails().Create(encModel.(*model.Email), schema)
	if err != nil {
		return nil, err
	}

	return createdEmail, nil
}

// UpdateEmail updates the account with the dto and applies the changes in the store
func UpdateEmail(s storage.Store, email *model.Email, dto *model.EmailDTO, schema string) (*model.Email, error) {
	rawModel := model.ToEmail(dto)
	encModel := EncryptModel(rawModel).(*model.Email)

	email.Title = encModel.Title
	email.Email = encModel.Email
	email.Password = encModel.Password

	updatedEmail, err := s.Emails().Update(email, schema)
	if err != nil {
		return nil, err
	}

	return updatedEmail, nil
}
