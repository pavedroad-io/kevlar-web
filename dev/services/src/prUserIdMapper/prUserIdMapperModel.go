// prUserIdMapperModel
// data structure and db/message interfaces

package main

import (
	"database/sql"
	_ "encoding/json"
	"errors"
	_ "errors"
	"fmt"
	"log"
	_ "os"
)

// A GenericError is the default error message that is generated.
// For certain status codes there are more appropriate error structures.
//
// swagger:response genericError
type GenericError struct {
	// in: body
	Body struct {
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// Return a basic message as json in the body
//
// swagger:response statusResponse
type statusResponse struct {
	// in: body
  msg error `json:"message"`
}

// Return list of mappings
//
// swagger:response tokenList
type listResponse struct {
	// in: body
  mappingList []prUserIdMapper `json:"mappings"`
}

// prUserIdMapper data structure for mapper storage
// It is used to map a 3rd party credential to an interanl UUID
//
// swagger:model userIdMapper
type prUserIdMapper struct {
  // API verions
  //
  // required: true
	APIVersion string `json:"apiVersion"`

  // Object verions
  //
  // required: true
	ObjVersion string `json:"objVersion"`

  // Type of object
  //
  // required: true
	Kind       string `json:"kind"`

  // 3rd party credential: email, login, or key
  //
  // required: true
	Credential string `json:"login"`

  // PR user ID 
  //
  // required: true
	UserUUID   string `json:"userUUID"`

  // Number of times this user has logged in
  //
  // readonly
	LoginCount int    `json:"loginCount"`

  // Time this record was creaetd
  //
  // readonly
	Created    string `json:"created,ignoreempty"`

  // Time this record was last updated
  //
  // readonly
	Updated    string `json:"updated,ignoreempty"`

  // Is this an active user true or false, default true
  //
  // required: false
	Active     string `json:"active"`
}
// An Mapper response model
//
// This is used for returning a response with a single mapper as body
//
// swagger:response mapperResponse
type MapperResponse struct {
	// in: body
	response string `json:"order"`
}

// updateUserIdMapper in database
func (t *prUserIdMapper) updateUserIdMapper(db *sql.DB) error {
	update := `
	UPDATE pavedroad.pruseridmapper
  SET apiVersion = '%s', 
      objVersion = '%s',
      kind = '%s', 
      userUUID = '%s', 
      loginCount = '%d', 
      updated = '%s', 
      active = '%s'
  WHERE
      credential = '%s';
  `
	statement := fmt.Sprintf(update, t.APIVersion, t.ObjVersion, t.Kind,
		t.UserUUID, t.LoginCount, t.Updated, t.Active, t.Credential)

	_, er1 := db.Query(statement)

	if er1 != nil {
		log.Println("Update failed")
		return er1
	}

	return nil
}

// createUserIdMapper in database
func (t *prUserIdMapper) createUserIdMapper(db *sql.DB) error {

	statement := fmt.Sprintf("INSERT INTO pavedroad.pruseridmapper(apiVersion, objVersion, kind, credential, userUUID, loginCount, created, updated, active) VALUES('%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', '%s') RETURNING credential;",
		t.APIVersion, t.ObjVersion, t.Kind, t.Credential, t.UserUUID, t.LoginCount,
		t.Created, t.Updated, t.Active)

	rows, er1 := db.Query(statement)

	if er1 != nil {
		log.Printf("Insert failed for: %s", t.Credential)
		log.Printf("SQL Error: %s", er1)
		return er1
	}

	defer rows.Close()

	return nil
}

// getUserIdMappers: return a list of tokens
//
func (t *prUserIdMapper) getUserIdMappers(db *sql.DB, start, count int) ([]prUserIdMapper, error) {
	statement := fmt.Sprintf("SELECT apiVersion, objVersion, kind, credential, userUUID, loginCount, created, updated, active FROM pavedroad.pruseridmapper LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mappings := []prUserIdMapper{}

	for rows.Next() {
		var t prUserIdMapper
		err := rows.Scan(&t.APIVersion, &t.ObjVersion, &t.Kind, &t.Credential,
			&t.UserUUID, &t.LoginCount, &t.Created, &t.Updated, &t.Active)

		if err != nil {
			log.Printf("SQL rows.Scan failed: %s", err)
			return mappings, err
		}

		mappings = append(mappings, t)
	}

	return mappings, nil
}

// getUserIdMapper: return a token based on credential
//
func (t *prUserIdMapper) getUserIdMapper(db *sql.DB, key string) error {
	statement := fmt.Sprintf("SELECT apiVersion, objVersion, kind, credential, userUUID, loginCount, created, updated, active FROM pavedroad.pruseridmapper WHERE credential = '%s'", key)
	row := db.QueryRow(statement)

	// Fill in mapper
	switch err := row.Scan(&t.APIVersion, &t.ObjVersion, &t.Kind, &t.Credential,
		&t.UserUUID, &t.LoginCount, &t.Created, &t.Updated, &t.Active); err {

	case sql.ErrNoRows:
		m := fmt.Sprintf("credential %s does not exist", key)
		return errors.New(m)
	case nil:
		break
	default:
		//Some error to catch
		panic(err)
	}

  // increment the login count
  t.LoginCount +=1
	statement = fmt.Sprintf("UPDATE pavedroad.pruseridmapper SET loginCount = '%d' WHERE credential = '%s'", t.LoginCount, key)

	_, er1 := db.Query(statement)

	if er1 != nil {
		log.Println("Incrementing login count failed")
  }

	return nil
}

// deleteUserIdMapper: return a token based on UID
//
func (t *prUserIdMapper) deleteUserIdMapper(db *sql.DB, cred string) error {
	statement := fmt.Sprintf("DELETE FROM pavedroad.pruseridmapper WHERE credential = '%s'", cred)
	result, err := db.Exec(statement)
	c, e := result.RowsAffected()

	if e == nil && c == 0 {
		em := fmt.Sprintf("credential %s does not exist", cred)
		log.Println(em)
		log.Println(e)
		return errors.New(em)
	}

	return err
}
