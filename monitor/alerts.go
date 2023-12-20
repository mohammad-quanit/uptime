package monitor

import (
	"encore.app/site"
	"encore.dev/pubsub"
)

// TransitionEvent describes a transition of a monitored site
// from up->down or from down->up.
type TransitionEvent struct {
	Site *site.Site `json:"site"` // Site is the monitored site in question.
	Up   bool       `json:"up"`   // Up specifies whether the site is now up or down (the new value).
}

// TransitionTopic is a pubsub topic with transition events for when a monitored site
// transitions from up->down or from down->up.
var TransitionTopic = pubsub.NewTopic[*TransitionEvent]("uptime-transition", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
