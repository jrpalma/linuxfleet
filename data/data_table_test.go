package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	os.Exit(m.Run())
}

func TestListByOwner(t *testing.T) {
	id1 := uuid.NewString()
	id2 := uuid.NewString()

	table, err := NewDataTables(db)
	assert.NoError(t, err)

	for _, tableName := range dataTableList() {
		query := sqlInsert(tableName)

		obj1 := Object{ID: id1, OwnerID: "owner1", Version: 1, Attributes: map[string]any{"attr1": "val1"}}
		obj2 := Object{ID: id2, OwnerID: "owner1", Version: 2, Attributes: map[string]any{"attr2": "val2"}}

		// Insert test objects
		_, err := db.Exec(query, obj1.ID, obj1.OwnerID, obj1.Version, jsonString(obj1.Attributes))
		assert.NoError(t, err)

		_, err = db.Exec(query, obj2.ID, obj2.OwnerID, obj2.Version, jsonString(obj2.Attributes))
		assert.NoError(t, err)

		objects, err := table.ListByOwner(tableName, "owner1")
		assert.NoError(t, err)

		assert.Len(t, objects, 2)
		assert.Equal(t, obj1, objects[0])
		assert.Equal(t, obj2, objects[1])
	}
}

func TestCreate(t *testing.T) {
	table, err := NewDataTables(db)
	assert.NoError(t, err)

	for _, tableName := range dataTableList() {
		query := sqlGetByID(tableName)
		objectID := uuid.NewString()

		obj := Object{ID: objectID, OwnerID: "owner1", Version: 1, Attributes: map[string]any{"attr1": "val1"}}

		err := table.Insert(tableName, obj)
		assert.NoError(t, err)

		var id string
		var ownerID string
		var version int
		var attrsJson string
		err = db.QueryRow(query, obj.ID).Scan(&id, &ownerID, &version, &attrsJson)
		assert.NoError(t, err)

		var retrievedAttrs map[string]any
		err = json.Unmarshal([]byte(attrsJson), &retrievedAttrs)
		assert.NoError(t, err)
		assert.Equal(t, obj, Object{ID: id, OwnerID: ownerID, Version: version, Attributes: retrievedAttrs})
	}

}

func TestDeleteByID(t *testing.T) {
	table, err := NewDataTables(db)
	assert.NoError(t, err)

	for _, tableName := range dataTableList() {
		insertQuery := sqlInsert(tableName)
		objectID := uuid.NewString()

		obj := Object{ID: objectID, OwnerID: "owner1", Version: 1, Attributes: map[string]any{"attr1": "val1"}}

		// Insert test object
		_, err := db.Exec(insertQuery, obj.ID, obj.OwnerID, obj.Version, jsonString(obj.Attributes))
		assert.NoError(t, err)

		err = table.DeleteByID(tableName, objectID)
		assert.NoError(t, err)

		var id string
		selectQuery := fmt.Sprintf("SELECT id FROM %s WHERE id = ?", tableName)
		err = db.QueryRow(selectQuery, obj.ID).Scan(&id)
		assert.EqualError(t, sql.ErrNoRows, err.Error())
	}

}

func TestUpdateByID(t *testing.T) {
	table, err := NewDataTables(db)
	assert.NoError(t, err)

	for _, tableName := range dataTableList() {
		inserQuery := sqlInsert(tableName)

		objectID := uuid.NewString()
		obj := Object{ID: objectID, OwnerID: "owner1", Version: 1, Attributes: map[string]any{"attr1": "val1"}}

		// Insert test object
		_, err := db.Exec(inserQuery, obj.ID, obj.OwnerID, obj.Version, jsonString(obj.Attributes))
		assert.NoError(t, err)

		obj.OwnerID = "updatedOwner"
		obj.Version = 2
		obj.Attributes = map[string]any{"attr1": "updatedVal", "newAttr": "newValue"}
		err = table.UpdateByID(tableName, objectID, obj)
		assert.NoError(t, err)

		var id string
		var ownerID string
		var version int
		var attrsJson string
		selectQuery := fmt.Sprintf("SELECT id, owner_id, version, attributes FROM %s WHERE id = ?", tableName)
		err = db.QueryRow(selectQuery, obj.ID).Scan(&id, &ownerID, &version, &attrsJson)
		assert.NoError(t, err)

		var retrievedAttrs map[string]any
		err = json.Unmarshal([]byte(attrsJson), &retrievedAttrs)
		assert.NoError(t, err)
		assert.Equal(t, obj, Object{ID: id, OwnerID: ownerID, Version: version, Attributes: retrievedAttrs})
	}

}

func TestGetByID(t *testing.T) {
	table, err := NewDataTables(db)
	assert.NoError(t, err)

	for _, tableName := range dataTableList() {
		query := sqlInsert(tableName)

		objectID := uuid.NewString()
		obj := Object{ID: objectID, OwnerID: "owner1", Version: 1, Attributes: map[string]any{"attr1": "val1"}}

		// Insert test object
		_, err := db.Exec(query, obj.ID, obj.OwnerID, obj.Version, jsonString(obj.Attributes))
		assert.NoError(t, err)

		retrievedObj, err := table.GetByID(tableName, objectID)
		assert.NoError(t, err)
		assert.Equal(t, obj, retrievedObj)
	}
}

func jsonString(attrs map[string]any) string {
	jsonBytes, _ := json.Marshal(attrs)
	return string(jsonBytes)
}
