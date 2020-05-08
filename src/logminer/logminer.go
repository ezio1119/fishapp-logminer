package logminer

import (
	"context"

	"github.com/ezio1119/fishapp-relaylog/conf"
	"github.com/ezio1119/fishapp-relaylog/domain"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

func StartPostDBLogMining(ctx context.Context, ch chan domain.BinlogEvent, pos mysql.Position) error {
	syncer := newPostDBSyncer()
	defer syncer.Close()

	streamer, err := syncer.StartSync(pos)
	if err != nil {
		return err
	}

	for {
		binlogEvent, err := streamer.GetEvent(ctx)
		if err != nil {
			return err
		}
		rowsEvent, ok := binlogEvent.Event.(*replication.RowsEvent)
		if ok {
			if string(rowsEvent.Table.Table) == "outbox" {
				ch <- domain.BinlogEvent{
					Pos:        syncer.GetNextPosition(),
					PosSubject: conf.C.Nats.Subject.PosPostDB,
					Event:      rowsEvent.Rows[0],
				}
			}
		}
	}
}

func newPostDBSyncer() *replication.BinlogSyncer {
	return replication.NewBinlogSyncer(replication.BinlogSyncerConfig{
		ServerID: 1,
		Flavor:   conf.C.PostDB.Dbms,
		Host:     conf.C.PostDB.Host,
		Port:     conf.C.PostDB.Port,
		User:     conf.C.PostDB.User,
		Password: conf.C.PostDB.Pass,
		Charset:  conf.C.PostDB.Charset,
	})
}
