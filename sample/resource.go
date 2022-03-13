package sample

import (
	"time"

	"github.com/antoniohueso/gplan"
)

// Resource Contiene información de un recurso
type Resource struct {
	*gplan.ResourceBase `bson:",inline"`
	Worked              uint        `json:"worked"`
	Holidays            []*Holidays `json:"holidays"`
}

func (s *Resource) Base() *gplan.ResourceBase {
	return s.ResourceBase
}

func (s *Resource) GetHolidays() []gplan.IHolidays {
	newArr := make([]gplan.IHolidays, len(s.Holidays))
	for i := range s.Holidays {
		newArr[i] = s.Holidays[i]
	}
	return newArr
}

// NewResource crea un nuevo recurso
func NewResource(id gplan.ResourceID, description string, resourceType string, availableFrom time.Time, holidays []*Holidays) *Resource {

	if holidays == nil {
		holidays = []*Holidays{}
	}

	return &Resource{
		ResourceBase: gplan.NewResourceBase(id, description, resourceType, availableFrom),
		Holidays:     holidays,
	}
}
