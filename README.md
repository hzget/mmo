# mmo

A mmo game backend that makes use of the [tcp server][tcpserver] framework.
It's based on the [zinx mmo demo][zinx mmo].
It's client and protocol can be found [here][mmo client].

The animation presentation:

![play](./pic/play.gif)

## functions

- [x] online, offline
- [x] chat
- [x] move
- [ ] jump
- [ ] beat
- [ ] weapon

## how to use the tcpserver framework

```golang
func main() {

	s := tcpserver.NewServer()
	// register handlers for different msg from the client
	s.AddRouter(core.MSG_C_Talk, &core.ChatRouter{})
	s.AddRouter(core.MSG_C_Move, &core.MoveRouter{})
	// hook func when the player is online
	s.SetOnConnStart(OnConnStart)
	// hook func when the player is offline
	s.SetOnConnStop(OnConnStop)

	s.Serve()
}

```

[tcpserver]: https://github.com/hzget/tcpserver
[mmo client]: https://github.com/aceld/zinx/tree/master/zinx_app_demo/mmo_game/game_client
[zinx mmo]: https://github.com/aceld/zinx/tree/master/zinx_app_demo/mmo_game
