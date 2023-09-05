package config

import (
	"context"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

//go:generate mockgen -destination=../mocks/mock_database_users.go -package=mocks github.com/mongodb/terraform-provider-mongodbatlas/config DatabaseUserCreator,DatabaseUserDeleter,DatabaseUserUpdater,DatabaseUserDescriber

type DatabaseUserCreator interface {
	CreateDatabaseUser(*atlas.DatabaseUser) (*atlas.DatabaseUser, error)
}

type DatabaseUserDeleter interface {
	DeleteDatabaseUser(string, string, string) error
}

type DatabaseUserUpdater interface {
	UpdateDatabaseUser(*atlas.DatabaseUser) (*atlas.DatabaseUser, error)
}

type DatabaseUserDescriber interface {
	DatabaseUser(string, string, string) (*atlas.DatabaseUser, error)
}

func (c *Config) CreateDatabaseUser(ctx context.Context, user *atlas.DatabaseUser) (*atlas.DatabaseUser, error) {
	dbUser, _, err := c.Client.Atlas.DatabaseUsers.Create(ctx, user.GroupID, user)
	return dbUser, err
}

func (c *Config) UpdateDatabaseUser(ctx context.Context, user *atlas.DatabaseUser) (*atlas.DatabaseUser, error) {
	dbUser, _, err := c.Client.Atlas.DatabaseUsers.Update(ctx, user.GroupID, user.Username, user)
	return dbUser, err
}

func (c *Config) DatabaseUser(ctx context.Context, authDB, groupID, username string) (*atlas.DatabaseUser, *atlas.Response, error) {
	return c.Client.Atlas.DatabaseUsers.Get(ctx, authDB, groupID, username)
}

func (c *Config) DeleteDatabaseUser(ctx context.Context, authDB, groupID, username string) error {
	_, err := c.Client.Atlas.DatabaseUsers.Delete(ctx, authDB, groupID, username)
	return err
}
