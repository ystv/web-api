// Package campus, this is a collection of useful endpoints relating to what campus
// will hopefully offer: current term, are you on a campus net, are you
// on the ystv net
package campus

import "github.com/jmoiron/sqlx"

type Campuser struct {
	db sqlx.DB
}
