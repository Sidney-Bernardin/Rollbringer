package database

import (
	"database/sql"
	"database/sql/driver"

	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
)

type HstoreArray []map[string]string

func (hstoreArray *HstoreArray) Scan(src any) error {
	hsAryCopy := []hstore.Hstore{}
	if err := pq.Array(&hsAryCopy).Scan(src); err != nil {
		return err
	}

	*hstoreArray = make([]map[string]string, len(hsAryCopy))
	for i, hs := range hsAryCopy {
		(*hstoreArray)[i] = map[string]string{}
		for k, v := range hs.Map {
			(*hstoreArray)[i][k] = v.String
		}
	}

	return nil
}

func (hstoreArray HstoreArray) Value() (driver.Value, error) {
	hsAryCopy := make([]hstore.Hstore, len(hstoreArray))

	for i, m := range hstoreArray {
		hsAryCopy[i] = hstore.Hstore{Map: map[string]sql.NullString{}}
		for k, v := range m {
			hsAryCopy[i].Map[k] = sql.NullString{
				String: v,
				Valid:  true,
			}
		}
	}

	return pq.Array(hsAryCopy).Value()
}
