package emitter

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/go-redis/redis/v8"
	"context"
)

const (
	FlagJson = "json"
	FlagVolatile = "volatile"
	FlagBroadcast = "broadcast"

	NormalEvent = 2
	BinaryEvent = 5

	DefaultRedisPrefix = "socket.io"
	DefaultUid = "sample_uid"
	DefaultNsp =  "/"
)

type Emitter struct {
	Redis *redis.Client
	Prefix string
	EventType int
	Nsp string
	Uid string
	Rooms  []string
	Flags map[string]bool
}

func (emitter Emitter) Of(Nsp string) Emitter {
	emitter.Nsp = Nsp
	return emitter
}

func (emitter Emitter) WithFlag(flag string, val bool) Emitter {
	flags := emitter.Flags
	if (flags == nil) {
		flags = map[string]bool{}
	}
	flags[flag] = val
	emitter.Flags = flags
	return emitter
}

func (emitter Emitter) In(room string) Emitter {
	emitter.Rooms = append(emitter.Rooms, room)
	return emitter
}

func (emitter Emitter) setDefaults() Emitter {
	if (emitter.Nsp == "") {
		emitter.Nsp = DefaultNsp
	}

	if (emitter.Uid == "") {
		emitter.Uid = DefaultUid
	}

	if (emitter.Prefix == "") {
		emitter.Prefix = DefaultRedisPrefix
	}

	if (emitter.EventType == 0) {
		emitter.EventType = NormalEvent
	}

	return emitter
}

func (emitter Emitter) Emit(args ...interface{}) {
	emitter = emitter.setDefaults()

	packet := make(map[string]interface{})
	packet["type"] = emitter.EventType
	packet["data"] = args
	packet["nsp"] = emitter.Nsp

	options := make(map[string]interface{})
	options["flags"] = emitter.Flags

	channleName := emitter.Prefix + "#" + emitter.Nsp + "#"
	if (len(emitter.Rooms) > 0){
		for _, room := range emitter.Rooms {
			options["rooms"] = []string{room}
			packedMsg, _ := msgpack.Marshal([]interface{}{emitter.Uid, packet, options})

			roomChannelName := channleName + room + "#"
			emitter.Redis.Publish(context.Background(), roomChannelName, packedMsg)
		}
	}else{
		options["rooms"] = emitter.Rooms
		packedMsg, _ := msgpack.Marshal([]interface{}{emitter.Uid, packet, options})
		emitter.Redis.Publish(context.Background(), channleName, packedMsg)
	}
}
