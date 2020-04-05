package contacts

import (
	"errors"
	"fmt"

	"github.com/LucasFrezarini/go-contacts/contacts/email"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// A Service contains all the business logic related to the contact resource in the application
type Service struct {
	Logger             *zap.Logger
	ContactsRepository Repository
	EmailRepository    email.GenericRepository
}

// ProvideContactsService creates a new Service with the provided dependencies.
// Created especially for the use of Wire, who will inject the dependencies via DI
func ProvideContactsService(logger *zap.Logger, cr Repository, er email.GenericRepository) *Service {
	return &Service{logger.Named("ContactsService"), cr, er}
}

// FindAllContacts fetches all the contacts registered in the application, as well as its emails and phones
func (s *Service) FindAllContacts() ([]*Contact, error) {
	contacts, err := s.ContactsRepository.FindAll()
	if err != nil {
		msg := fmt.Sprintf("FindAllContacts() error while trying to fetch contacts: %v", err)
		s.Logger.Error(msg)
		return nil, errors.New(msg)
	}

	for _, c := range contacts {
		emails, err := s.EmailRepository.FindByContactID(c.ID)
		if err != nil {
			msg := fmt.Sprintf("FindAllContacts() error while trying to fetch contact's emails: %v", err)
			s.Logger.Error(msg)
			return nil, errors.New(msg)
		}

		c.Emails = emails
	}

	return contacts, nil
}

// ServiceSet is a wire set which contains all the bindings needed for creating a new service
var ServiceSet = wire.NewSet(ProvideContactsService, RepositorySet, email.EmailRepositorySet)