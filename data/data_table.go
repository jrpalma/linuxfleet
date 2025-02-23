package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Object struct {
	ID         string         `json:"id"`
	CreatedAt  Timestamp      `json:"created_at"`
	UpdatedAt  Timestamp      `json:"updated_at"`
	OwnerID    string         `json:"owner_id"`
	Version    int            `json:"version"`
	Attributes map[string]any `json:"attributes"`
}

type Tables struct {
	db *sql.DB
}

// NewTables creates a new data tables object from the sql DB.
func NewTables(db *sql.DB) (*Tables, error) {
	for _, tableName := range dataTableList() {
		query := sqlCreatTable(tableName)
		_, err := db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	return &Tables{db: db}, nil
}

// ListByOwner retrieves a list of objects by owner ID and object type from the database.
func (table *Tables) ListByOwner(tableName string, ownerID string) ([]Object, error) {
	query := sqlListByOwner(tableName)
	rows, err := table.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []Object
	for rows.Next() {
		var obj Object
		var attrsJson string
		err := rows.Scan(&obj.ID, &obj.CreatedAt, &obj.UpdatedAt, &obj.OwnerID, &obj.Version, &attrsJson)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(attrsJson), &obj.Attributes)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, rows.Err()
}

// Insert inserts a new object into the specified table in the database.
func (table *Tables) Insert(tableName string, obj Object) error {
	attrsJson, err := json.Marshal(obj.Attributes)
	if err != nil {
		return err
	}
	query := sqlInsert(tableName)
	obj.CreatedAt = NowTimestamp()
	_, err = table.db.Exec(query, obj.ID, obj.CreatedAt, obj.UpdatedAt, obj.OwnerID, obj.Version, attrsJson)
	return err
}

// DeleteByID deletes an object from the specified table in the database by its ID.
func (table *Tables) DeleteByID(tableName string, id string) error {
	query := sqlDeleteByID(tableName)
	_, err := table.db.Exec(query, id)
	return err
}

// UpdateByID updates an existing object in the specified table in the database by its ID.
func (table *Tables) UpdateByID(tableName string, id string, obj Object) error {
	attrsJson, err := json.Marshal(obj.Attributes)
	if err != nil {
		return err
	}
	query := sqlUpdateByID(tableName)
	obj.UpdatedAt = NowTimestamp()
	_, err = table.db.Exec(query, obj.UpdatedAt, obj.OwnerID, obj.Version, attrsJson, id)
	return err
}

// GetByID retrieves an object from the specified table in the database by its ID.
func (table *Tables) GetByID(tableName, id string) (Object, error) {
	query := sqlGetByID(tableName)
	var obj Object
	var attrsJson string
	err := table.db.QueryRow(query, id).Scan(&obj.ID, &obj.CreatedAt, &obj.UpdatedAt, &obj.OwnerID, &obj.Version, &attrsJson)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal([]byte(attrsJson), &obj.Attributes)
	return obj, err
}

// sqlGetByID constructs the SQL query to retrieve an object by its ID from the specified table.
func sqlGetByID(tableName string) string {
	query := `SELECT id, created_at, updated_at, owner_id, version, attributes FROM %s WHERE id = ?`
	return fmt.Sprintf(query, tableName)
}

// sqlUpdateByID constructs the SQL query to update an object by its ID in the specified table.
func sqlUpdateByID(tableName string) string {
	query := `UPDATE %s SET updated_at = ?, owner_id = ?, version = ?, attributes = ? WHERE id = ?`
	return fmt.Sprintf(query, tableName)
}

// sqlDeleteByID constructs the SQL query to delete an object by its ID from the specified table.
func sqlDeleteByID(tableName string) string {
	query := `DELETE FROM %s WHERE id = ?`
	return fmt.Sprintf(query, tableName)
}

// sqlInsert constructs the SQL query to insert a new object into the specified table.
func sqlInsert(tableName string) string {
	query := `INSERT INTO %s (id, created_at, updated_at, owner_id, version, attributes) VALUES (?, ?, ?, ?, ?, ?)`
	return fmt.Sprintf(query, tableName)
}

// sqlListByOwner constructs the SQL query to list objects by their owner ID from the specified table.
func sqlListByOwner(tableName string) string {
	query := `SELECT id, created_at, updated_at, owner_id, version, attributes FROM %s WHERE owner_id = ?`
	return fmt.Sprintf(query, tableName)
}

func sqlCreatTable(tableName string) string {
	query := `
	CREATE TABLE IF NOT EXISTS %s(
		id TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		owner_id TEXT NOT NULL,
		version TEXT,
		attributes TEXT,
		PRIMARY KEY (id)
	);
	CREATE INDEX IF NOT EXISTS idx_owner_id ON %s(owner_id);
	`
	return fmt.Sprintf(query, tableName, tableName)
}

func dataTableList() []string {
	return []string{
		"admins",
		"users",
		"devices",
		"signup",
	}
}
