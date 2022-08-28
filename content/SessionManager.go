package content

import (
	"log"
	"net"
	"sync"
)

type GlobalSession struct {
	WorkerNum int
	Jobs chan
}

var instance_gs *GlobalSession
var once_gs sync.Once

func GetGlobalSession() *GlobalSession {
	once_gs.Do(func() {
		instance_gs = &GlobalSession{}
	})
	return instance_gs
}

func (gs *GlobalSession) Init() {
	log.Println("INIT_GlobalSession")
	gs.WorkerNum = 10
	gs.Jobs := make(chan int, numJobs)
}

func (gs *GlobalSession) SendByteByWoker(c net.Conn, data []byte) {
	go gs.worker(c, data)

}
func (gs *GlobalSession) worker(c net.Conn, data []byte) {
	for j := range gs.Jobs {
		gs.SendByte(c, data)
	}
}
func (gs *GlobalSession) SendByte(c net.Conn, data []byte) {
	if c != nil {
		sent, err := c.Write(data)
		if err != nil {
			log.Println("SendPacket ERROR :", err)
		} else {
			if sent != len(data) {
				log.Println("[Sent diffrent size] : SENT =", sent, "BufferSize =", len(data))
			}
			log.Println("SendPacket : ", data)
		}
	}
}

func (gs *GlobalSession) BroadCast(buff []byte) {
	GetContentManager().Players.Range(func(key, value any) bool {
		gs.SendByte(key.(net.Conn), value.([]byte))
		return true
	})
}
