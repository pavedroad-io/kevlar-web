// prTokenModel
// data structure and db/message interfaces

package main

import (
	_ "database/sql"
	_ "fmt"
	"time"
)

type PrToken struct {
	ApiVersion   string    `json:"core.pavedroad.io/v1alpha1"`
	Kind         string    `json:"PrToken"`
    Metadata struct {
	    Name      string    `json:"defauult-name"`
	    Namespace string    `json:"defauult-name"`
	    Uid       string    `json:"defauult-name"`
	    Token     string    `json:"defauult-name"`
        Scope     []string
    }
	Created       time.Time
	Updated       time.Time
	Active        bool       `json:"True"`
}
