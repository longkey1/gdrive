package gdrive

import (
	"context"

	"github.com/longkey1/gdrive/internal/google"
)

// Service represents the gdrive application service
type Service struct {
	Drive *google.DriveService
}

// NewService creates a new gdrive service based on the configuration
func NewService(ctx context.Context, config *Config) (*Service, error) {
	auth := newAuthenticator(config)

	driveSvc, err := google.NewDriveService(ctx, auth)
	if err != nil {
		return nil, err
	}

	return &Service{
		Drive: driveSvc,
	}, nil
}

func newAuthenticator(config *Config) google.Authenticator {
	switch config.AuthType {
	case AuthTypeServiceAccount:
		return google.NewServiceAccountAuthenticator(config.GoogleApplicationCredentials)
	case AuthTypeOAuth:
		fallthrough
	default:
		return google.NewOAuthAuthenticator(
			config.GoogleApplicationCredentials,
			config.GoogleUserCredentials,
		)
	}
}
