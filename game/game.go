package game

import (
	"2018_2_YetAnotherGame/player"
	"2018_2_YetAnotherGame/room"
	"fmt"
	"log"

	"github.com/google/uuid"

	"golang.org/x/net/websocket"
)

type Game struct {
	Rooms    map[string]*room.Room
	MaxRooms int
	Register chan *websocket.Conn
}

func New() *Game {
	return &Game{
		Rooms:    make(map[string]*room.Room),
		MaxRooms: 2,
		Register: make(chan *websocket.Conn),
	}
}

func (g *Game) Run() {
	fmt.Println("fff")
	for {
		conn := <-g.Register
		g.ProcessConn(conn)
	}
}

func (g *Game) FindRoom() *Room {
	for _, r := range g.Rooms {
		if len(r.Players) < r.MaxPlayers {
			return r
		}
	}
	if len(g.Rooms) >= g.MaxRooms {
		return nil
	}
	r := New()
	go r.ListenToPlayers()
	g.Rooms[r.ID] = r
	log.Println("room %s created", r.ID)
	return r
}

func (g *Game) ProcessConn(conn *websocket.Conn) {
	id := uuid.New().String()

	p := &player.Player{
		Conn: conn,
		ID:   id,
	}
	r := g.FindRoom()
	if r == nil {
		return
	}

	r.Players[p.ID] = p
	p.Room = r
	log.Println("player %s joined room %s", p.ID, r.ID)
	go p.Listen()
	fmt.Println(r.Players)
	if len(r.Players) == r.MaxPlayers {
		go r.Run()
	}

}