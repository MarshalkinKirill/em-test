package em

import (
	"emsrv/pkg/db"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type Filter struct {
	PersonID    uuid.UUID `json:"personId" pg:"personId"`
	Name        string    `json:"name" pg:"name" validate:"required,max=25"`
	Surname     string    `json:"surname" pg:"surname" validate:"required,max=25"`
	Patronymic  string    `json:"patronymic" pg:"patronymic" validate:"required,max=25"`
	Age         int       `json:"age" pg:"age" validate:"age=0"`
	Gender      string    `json:"gender" pg:"gender" validate:"oneof=male female other Male Female Other"`
	Nationality string    `json:"nationality" pg:"nationality"`
}

type Pager struct {
	Page  int
	Limit int
}
type FiltersParams struct {
	Filter Filter
	Pager  Pager
}

func (em EmService) GetPersonHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var filterParams FiltersParams

	if err := c.Bind(&filterParams); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	persons, err := em.getPersons(ctx, filterParams)
	if err != nil {
		return err
	}
	em.Printf("[AllPersonHandler] all person get sucess\n")
	return c.JSON(http.StatusOK, persons)
}

func (em EmService) CreatePersonHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var person db.Person

	if err := c.Bind(&person); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	fmt.Printf("[CreatePersonHandler] name is: %s\n", person.Name)
	if age, err := em.GetRelevantAge(c, person.Name); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	} else {
		person.Age = age
	}
	if gender, err := em.GetRelevantGender(c, person.Name); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	} else {
		person.Gender = gender
	}
	if nationality, err := em.GetRelevantNationality(c, person.Name); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	} else {
		person.Nationality = nationality
	}

	err := em.createPerson(ctx, &person)
	if err != nil {
		return err
	}
	em.Printf("[CreatePersonHandler] person with id: %d - create success\n", person.PersonID)
	return c.JSON(http.StatusOK, map[string]string{"message": "person create success"})
}

func (em EmService) UpdatePersonHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var person db.UpdatePerson
	if err := c.Bind(&person); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	id := c.Param("personId")
	err := em.updatePerson(ctx, &person, id)
	if err != nil {
		return err
	}
	em.Printf("[AllPersonHandler] person with id: %d - update success\n", id)
	return c.JSON(http.StatusOK, map[string]string{"message": "person create success"})
}

func (em EmService) DeletePersonHandler(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("personId")
	fmt.Printf("id: %s", id)
	err := em.deletePerson(ctx, id)
	if err != nil {
		return err
	}
	em.Printf("[AllPersonHandler] person with id: %d - get sucess\n", id)
	return c.JSON(http.StatusOK, map[string]string{"message": "Person deleted successfully"})
}

func (em EmService) GetRelevantAge(c echo.Context, name string) (int, error) {
	url := fmt.Sprintf("%s%s", em.cfg.AgeUrl, name)
	data, err := em.Call(url)
	if err != nil {
		return 0, err
	}
	if age, ok := data["age"].(float64); ok != false {
		fmt.Printf("[GetRelevantGender] age is: %d'\n", int(age))
		return int(age), nil
	} else {
		fmt.Printf("[GetRelevantGender] age is: %d\n", int(age))
		return 0, nil
	}
}

func (em EmService) GetRelevantGender(c echo.Context, name string) (string, error) {
	url := fmt.Sprintf("%s%s", em.cfg.GenderUrl, name)
	data, err := em.Call(url)
	if err != nil {
		return "", err
	}
	if gender, ok := data["gender"].(string); ok != false {
		fmt.Printf("[GetRelevantGender] gender is: %s\n", gender)
		return gender, nil
	} else {
		fmt.Printf("[GetRelevantGender] gender is: %s\n", gender)
		return "", nil
	}
}

func (em EmService) GetRelevantNationality(c echo.Context, name string) (string, error) {
	url := fmt.Sprintf("%s%s", em.cfg.NationalityUtl, name)
	fmt.Printf("%s", url)
	data, err := em.Call(url)
	if err != nil {
		return "", err
	}

	if countries, ok := data["country"].([]interface{}); ok {
		var maxProbability float64
		var CountryID string
		for _, country := range countries {
			if countryMap, isMap := country.(map[string]interface{}); isMap {
				// Получение значений country_id и probability
				countryID, idOK := countryMap["country_id"].(string)
				probability, probOK := countryMap["probability"].(float64)

				// Проверка успешности преобразований типов
				if idOK && probOK {
					// Сравнение с текущей максимальной вероятностью
					if probability > maxProbability {
						maxProbability = probability
						CountryID = countryID
					}
				}
			}
		}
		fmt.Printf("[GetRelevantNationality] relevant country is: %s\n", CountryID)
		return CountryID, nil
	} else {
		return "", nil
	}
}

func (em EmService) Call(url string) (map[string]interface{}, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("[GetRelevantAge] http request failed: %w\n", err)
		return nil, fmt.Errorf("[GetRelevantAge] http request failed: %w\n", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("[GetRelevantAge] request body read failed: %w\n", err)
		return nil, fmt.Errorf("[GetRelevantAge] request body read failed: %w\n", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("[GetRelevantAge] unmurshal body failed: %w\n", err)
		return nil, fmt.Errorf("[GetRelevantAge] unmurshal body failed: %w\n", err)
	}
	return data, nil
}
