package agent

import (
	"github.com/khorevaa/go-v8platform/types"
	agent "github.com/v8platform/agent"
)

type ClientPool struct {
	pool         map[string]agent.Agent
	OnConnect    func(ConnectString string)
	OnDisconnect func(ConnectString string)
}

type AgentPool struct {
	pool         map[string]agent.Agent
	OnCreate     func(ConnectString string)
	OnDisconnect func(ConnectString string)
}

type RunningAgent struct {
	connectionString string

	//agent.AgentModeOptions

	// Признак запуска конфигуратора в режиме анета
	Running bool

	// Канал для остановки режима агента
	stop chan struct{}
}

func (s RunningAgent) Stop() {

	s.stop <- struct{}{}

}

func (s RunningAgent) Start() {

	//go func() {
	//
	//}()

}

func RunOnAgent(where types.InfoBase, what types.Command, opts ...interface{}) {

}
