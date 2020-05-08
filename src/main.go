package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ezio1119/fishapp-relaylog/conf"
	"github.com/ezio1119/fishapp-relaylog/domain"
	"github.com/ezio1119/fishapp-relaylog/logminer"
	"github.com/ezio1119/fishapp-relaylog/publisher"
	"github.com/nats-io/stan.go"
	"github.com/siddontang/go-mysql/mysql"
)

func main() {
	ctx := context.Background()
	conn, err := newNatsConn()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	pos, err := getLastPosPostDB(conn)
	if err != nil {
		log.Fatal(err)
	}

	eventChan := make(chan domain.BinlogEvent)

	for i := 0; i < conf.C.Nats.PublisherNum; i++ {
		go publisher.StartEventPublishing(conn, eventChan)
	}

	if err := logminer.StartPostDBLogMining(ctx, eventChan, pos); err != nil {
		log.Printf("error: %s", err)
	}
}

func newNatsConn() (stan.Conn, error) {
	return stan.Connect(conf.C.Nats.ClusterID, conf.C.Nats.ClientID, stan.NatsURL(conf.C.Nats.URL))
}

func getLastPosPostDB(conn stan.Conn) (mysql.Position, error) {
	ch := make(chan mysql.Position)

	if _, err := conn.Subscribe(conf.C.Nats.Subject.PosPostDB, func(msg *stan.Msg) {

		pos := mysql.Position{}
		if err := json.Unmarshal(msg.MsgProto.Data, &pos); err != nil {
			log.Fatal(err)
		}

		if err := msg.Sub.Close(); err != nil {
			log.Fatal(err)
		}

		ch <- pos
	}, stan.StartWithLastReceived()); err != nil {
		log.Printf("error: %s", err)
	}

	return <-ch, nil
}
