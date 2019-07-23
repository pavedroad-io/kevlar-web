// prUserModel
// data structure and db/message interfaces

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	_ "os"
)

// prUser data structure
type prUser struct {
	APIVersion string                `json:"apiVersion"`
	ObjVersion string                `json:"objVersion"`
	Kind       string                `json:"kind"`
	Metadata   userMetadata          `json: "metadata"`
  Labels     map[string](string)   `json: "labels"`
	Created    string                `json:"created,ignoreempty"`
	Updated    string                `json:"updated,ignoreempty"`
	Active     string                `json:"active"`
}

// user metadata
type userMetadata struct {
	Name      string                `json:"name"`
	Namespace string                `json:"namespace"`
	UUID      string                `json:"uuid"`
	Profiles  map[string](string)   `json:"profiles,ignoreempty"`
}

// userList of selected attributes
type userList struct {
	UUID   string
	Active string
	Name   string
}

// updateUser in database
func (t *prUser) updateUser(db *sql.DB, key string) error {
	update := `UPDATE pavedroad.pruser
  SET pruser = '%s'
  WHERE UUID = '%s';`

	jb, err := json.Marshal(t)
	if err != nil {
		log.Println("marshall failed")
		panic(err)
	}

	statement := fmt.Sprintf(update, jb, key)
	_, er1 := db.Query(statement)

	if er1 != nil {
		log.Println("Update failed")
		return er1
	}

	return nil
}

// createUser in database
func (t *prUser) createUser(db *sql.DB) (string, error) {
	jb, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	statement := fmt.Sprintf("INSERT INTO pavedroad.pruser(pruser) VALUES('%s') RETURNING uuid", jb)
	rows, er1 := db.Query(statement)

	if er1 != nil {
		log.Printf("Insert failed for: %s", t.Metadata.Name)
		log.Printf("SQL Error: %s", er1)
		return "", er1
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&t.Metadata.UUID)
		if err != nil {
			return "", err
		}
	}

	return t.Metadata.UUID, nil
}

// getUsers: return a list of users
//
func (t *prUser) getUsers(db *sql.DB, start, count int) ([]userList, error) {
	qry := `select uuid, 
          pruser ->> 'active' as active,  
          pruser -> 'Metadata' ->> 'name' as name 
          from pavedroad.pruser LIMIT %d OFFSET %d;`
	statement := fmt.Sprintf(qry, count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ul := []userList{}

	for rows.Next() {
		var t userList
		err := rows.Scan(&t.UUID, &t.Active, &t.Name)

		if err != nil {
			log.Printf("SQL rows.Scan failed: %s", err)
			return ul, err
		}

		ul = append(ul, t)
	}

	return ul, nil
}

// getUser: return a user
// method is either uuid, the default, or name
func (t *prUser) getUser(db *sql.DB, key string, method int) error {

	var statement string

	switch method {
	case UUID:
		statement = fmt.Sprintf(`
  SELECT uuid, pruser
  FROM pavedroad.pruser 
  WHERE uuid = '%s';`, key)
	case NAME:
		statement = fmt.Sprintf(`
  SELECT uuid, pruser
  FROM pavedroad.pruser 
  WHERE pruser -> 'Metadata' ->> 'name' = '%s';`, key)
	}
	row := db.QueryRow(statement)

	// Fill in mapper
	var jb []byte
	var uid string
	switch err := row.Scan(&uid, &jb); err {

	case sql.ErrNoRows:
		m := fmt.Sprintf("name %s does not exist", key)
		return errors.New(m)
	case nil:
		err = json.Unmarshal(jb, t)
		if err != nil {
			m := fmt.Sprintf("unmarshal failed %s", key)
			return errors.New(m)
		}
		t.Metadata.UUID = uid
		break
	default:
		//Some error to catch
		panic(err)
	}

	return nil
}

// deleteUser: return a user based on UID
//
func (t *prUser) deleteUser(db *sql.DB, uuid string) error {
	statement := fmt.Sprintf("DELETE FROM pavedroad.pruser WHERE uuid = '%s'", uuid)
	result, err := db.Exec(statement)
	c, e := result.RowsAffected()

	if e == nil && c == 0 {
		em := fmt.Sprintf("UUID %s does not exist", uuid)
		log.Println(em)
		log.Println(e)
		return errors.New(em)
	}

	return err
}
