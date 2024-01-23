package em

import (
	"emsrv/pkg/db"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
	"reflect"
)

type EmsrvRepo struct {
	db   orm.DB
	join map[string][]string
}

// NewEmsrvRepo returns new repository
func NewEmsrvRepo(db orm.DB) EmsrvRepo {
	return EmsrvRepo{
		db: db,
	}
}

// WithTransaction is a function that wraps EmsrvRepo with pg.Tx transaction.
func (em EmsrvRepo) WithTransaction(tx *pg.Tx) EmsrvRepo {
	em.db = tx
	return em
}

func (em EmsrvRepo) GetPersons(persons *[]db.Person, params FiltersParams) error {
	query := em.db.Model(&db.Person{})

	filtersValue := reflect.ValueOf(params.Filter)
	for i := 0; i < filtersValue.NumField(); i++ {
		field := filtersValue.Field(i)
		if !reflect.ValueOf(field.Interface()).IsZero() {
			query = query.Where(fmt.Sprintf("%s = ?", filtersValue.Type().Field(i).Tag.Get("pg")), field.Interface())
		}
	}

	query = query.Limit(params.Pager.Limit).Offset(params.Pager.Page)

	err := query.Select(persons)
	if err != nil {
		return fmt.Errorf("[GetPersons] select person failed: %w\n", err)
	}
	return nil
}

func (em EmsrvRepo) GetPersonByID(person *db.Person, personID uuid.UUID) error {
	return em.db.Model(person).Where("\"personId\" = ?", personID).Select()
}

func (em EmsrvRepo) CreatePerson(person *db.Person) error {
	_, err := em.db.Model(person).Insert()
	if err != nil {
		return fmt.Errorf("[CreatePerson] create person failed: %w\n", err)
	}
	return nil
}

func (em EmsrvRepo) UpdatePerson(person *db.UpdatePerson) error {
	_, err := em.db.Model(person).Where("\"personId\" = ?", person.PersonID).UpdateNotZero()
	if err != nil {
		return fmt.Errorf("[UpdatePerson] update person failed: %w\n", err)
	}
	return nil
}

func (em EmsrvRepo) DeletePerson(personID uuid.UUID) error {
	result, err := em.db.Model((*db.Person)(nil)).Where("\"personId\" = ?", personID).Delete()
	if err != nil {
		return fmt.Errorf("[DeletePerson] delete person failed: %w\n", err)
	}
	if result.RowsAffected() < 1 {
		return fmt.Errorf("[DeletePerson] object does not exist")
	}
	return nil
}
