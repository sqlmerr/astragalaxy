package model

import (
	"context"

	"github.com/google/uuid"
)

type Event interface {
	EventType() string
}

type EventBus interface {
	Publish(ctx context.Context, events ...Event)
}

type AgentWarpedEvent struct {
	AgentID uuid.UUID `json:"agent_id"`
	SystemX int       `json:"system_x"`
	SystemY int       `json:"system_y"`
}

func (AgentWarpedEvent) EventType() string {
	return "agent.warped"
}
