package livekit

import (
	"time"

	"github.com/livekit/protocol/auth"
)

// TokenGenerator is an interface for generating tokens
type TokenGenerator interface {
	// GenerateToken generates a token for the given identity
	GenerateToken(identity string, roomName string) (string, error)
}

// DefaultTokenGenerator is the default token generator
type DefaultTokenGenerator struct {
	apiKey    string
	apiSecret string
}

// NewDefaultTokenGenerator creates a new DefaultTokenGenerator
func NewDefaultTokenGenerator(apiKey string, apiSecret string) *DefaultTokenGenerator {
	return &DefaultTokenGenerator{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// GenerateToken generates a token for the given identity
func (g *DefaultTokenGenerator) GenerateToken(identity string, roomName string) (string, error) {
	at := auth.NewAccessToken(g.apiKey, g.apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin:   true,
		RoomCreate: true,
		Room:       roomName,
	}
	grant.SetCanPublishData(true)
	grant.SetCanPublish(true)
	grant.SetCanSubscribe(true)
	at.SetVideoGrant(grant).SetIdentity(identity).SetValidFor(time.Hour * 24)
	return at.ToJWT()
}
