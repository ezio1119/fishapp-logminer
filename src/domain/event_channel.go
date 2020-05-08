package domain

import (
	"github.com/siddontang/go-mysql/mysql"
)

type BinlogEvent struct {
	Pos        mysql.Position
	PosSubject string
	Event      []interface{}
}
