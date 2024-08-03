package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qbohyzmcx7c5w6q")
		if err != nil {
			return err
		}

		// add
		new_id_employee_responsible := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "wirn82wi",
			"name": "id_employee_responsible",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "bda9gk4pcxm2psb",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), new_id_employee_responsible); err != nil {
			return err
		}
		collection.Schema.AddField(new_id_employee_responsible)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qbohyzmcx7c5w6q")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("wirn82wi")

		return dao.SaveCollection(collection)
	})
}
