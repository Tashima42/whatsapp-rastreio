package routines

import (
	"database/sql"
	"fmt"
	"github.com/tashima42/whatsapp-rastreio/data"
	"github.com/tashima42/whatsapp-rastreio/helpers"
	"github.com/tashima42/whatsapp-rastreio/providers"
)

type ObjectsRoutines struct {
	DB               *sql.DB
	CorreiosProvider providers.CorreiosProvider
}

func (or *ObjectsRoutines) UpdateObjectsEvents() {
	fmt.Println("update")
	events, objects, err := data.GetObjectsByLastUpdated(or.DB, helpers.NowMinusMinutes(5))
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	if len(events) != len(objects) {
		fmt.Println("events and objects must have the same length")
		return
	}
	if len(events) <= 0 {
		fmt.Println("There must be at least one event")
		return
	}
	go or.updateEvents(events, objects)
	// separate them in chunks
	// run in go routines
	// decide to whether updated or not the record
}

func (or *ObjectsRoutines) updateEvents(events []data.Event, objects []data.Object) {
	objectsCodes := filterObjectsCodes(objects)
	correiosObjects, err := or.CorreiosProvider.GetCorreiosObjects(objectsCodes)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	for i, correiosObject := range correiosObjects {
		event := correiosObject.Event[0].ToEvent()
		if events[i].Hash != event.Hash {
			events[i].Body = event.Body
			events[i].Hash = event.Hash
			err = events[i].Update(or.DB)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				return
			}
		}
	}
}

func filterObjectsCodes(objects []data.Object) []string {
	var objectsCodes []string
	for _, object := range objects {
		objectsCodes = append(objectsCodes, object.Code)
	}
	return objectsCodes
}
