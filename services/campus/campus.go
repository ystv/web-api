// Package campus collection of useful endpoints relating to campus,
// will hopefully offer: current term, are you on campus net, are you
// on the ystv net
package campus

import "github.com/jackc/pgx"

type Campuser struct {
	db pgx.Conn
}
