package em

import (
	"context"
	"emsrv/pkg/db"
	"emsrv/pkg/embedlog"
	"fmt"
	"github.com/google/uuid"
)

type Config struct {
	AgeUrl         string
	GenderUrl      string
	NationalityUtl string
}

type EmService struct {
	embedlog.Logger
	emRepo EmsrvRepo
	cfg    Config
}

func NewEmService(logger embedlog.Logger, dbo db.DB, cfg Config) *EmService {
	return &EmService{
		emRepo: NewEmsrvRepo(dbo),
		Logger: logger,
		cfg:    cfg,
	}
}

func (em EmService) getPersons(c context.Context, params FiltersParams) ([]db.Person, error) {
	var persons []db.Person
	err := em.emRepo.GetPersons(&persons, params)
	if err != nil {
		em.Errorf("[EmService][getPerson] error log: %w", err)
		return nil, err
	}
	return persons, nil
}

func (em EmService) createPerson(c context.Context, person *db.Person) error {
	if err := em.isValid(c, person); err != nil {
		em.Errorf("[EmService][createPerson] error log: %w", err)
		return err
	}
	var emptyUUID uuid.UUID
	if person.PersonID == emptyUUID {
		person.PersonID = uuid.New()
	}
	err := em.emRepo.CreatePerson(person)
	if err != nil {
		em.Errorf("[EmService][createPerson] error log: %w", err)
		return err
	}
	return nil
}

func (em EmService) updatePerson(c context.Context, person *db.UpdatePerson, id string) error {
	if err := em.isValidUpdate(c, person); err != nil {
		return err
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		em.Errorf("[EmService][updatePerson] error log: %w", err)
		return fmt.Errorf("[updatePerson] id parse failed: %w", err)
	}
	person.PersonID = uuid
	//fmt.Printf("", person)
	err = em.emRepo.UpdatePerson(person)
	if err != nil {
		em.Errorf("[EmService][updatePerson] error log: %w", err)
		return err
	}
	return nil
}

func (em EmService) deletePerson(c context.Context, id string) error {
	//var persons db.Person
	uuid, err := uuid.Parse(id)
	if err != nil {
		em.Errorf("[EmService][deletePerson] error log: %w", err)
		return fmt.Errorf("[deletePerson] id parse failed: %w", err)
	}
	err = em.emRepo.DeletePerson(uuid)
	if err != nil {
		em.Errorf("[EmService][deletePerson] error log: %w", err)
		return err
	}
	return nil
}
