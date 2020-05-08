package publisher

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"github.com/ezio1119/fishapp-relaylog/domain"
	"github.com/ezio1119/fishapp-relaylog/domain/event"
	"github.com/nats-io/stan.go"
	"google.golang.org/protobuf/encoding/protojson"
)

func StartEventPublishing(conn stan.Conn, ch chan domain.BinlogEvent) {
	for bEvent := range ch {
		fmt.Println("NumGoroutine: ", runtime.NumGoroutine())
		e := convEventProto(bEvent.Event)

		eventByte, err := protojson.Marshal(e)
		if err != nil {
			log.Fatal(err)
		}

		if err := conn.Publish(e.EventType, eventByte); err != nil {
			log.Println(err)
		}
		log.Printf("success: published event: %s", e.EventType)

		posByte, err := json.Marshal(bEvent.Pos)
		if err != nil {
			log.Fatal(err)
		}

		if err := conn.Publish(bEvent.PosSubject, posByte); err != nil {
			log.Println(err)
		}
		log.Printf("success: published position: %+v", bEvent.Pos)
	}
}

func convEventProto(raw []interface{}) *event.Event {
	e := &event.Event{}
	// var ok bool
	id, ok := raw[0].(int32)
	if !ok {
		log.Fatal("failed conv event proto")
	}
	e.Id = fmt.Sprint(id)

	e.EventType, ok = raw[1].(string)
	if !ok {
		log.Fatal("failed conv event proto")
	}

	e.EventData, ok = raw[2].([]byte)
	if !ok {
		log.Fatal("failed conv event proto")
	}

	e.AggregateId, ok = raw[3].(string)
	if !ok {
		log.Fatal("failed conv event proto")
	}

	e.AggregateType, ok = raw[4].(string)
	if !ok {
		log.Fatal("failed conv event proto")
	}

	return e
}
