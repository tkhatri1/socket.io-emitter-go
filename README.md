
# SocketIO::Emitter

Golang version of [socket.io-emitter](https://github.com/Automattic/socket.io-emitter).

## How to use

```go
import "github.com/go-redis/redis/v8"

rdb := redis.NewClient(&redis.Options{
  Addr:     "localhost:6379",
  Password: "",
  DB:       0,
})

Emitter{Redis: rdb}.Of("/").In("chat_roome_1").In("chat_room_2").Emit("event_name", map[string]interface{}{"sample": "data"})

With namespace

Emitter{Redis: rdb}.Of("/namespace").In("chat_roome_1").In("chat_room_2").Emit("event_name", map[string]interface{}{"sample": "data"})
```