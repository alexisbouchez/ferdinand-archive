package config

import (
	"ferdinand/app/listeners"

	"github.com/caesar-rocks/events"
)

func RegisterEventsEmitter(
	usersListener *listeners.UsersListener,
) *events.EventsEmitter {
	emitter := events.NewEventsEmitter()
	emitter.On("users.created", usersListener.OnCreated)

	return emitter
}
