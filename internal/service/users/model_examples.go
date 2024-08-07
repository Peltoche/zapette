package users

import (
	"time"

	"github.com/Peltoche/zapette/internal/tools/secret"
	"github.com/Peltoche/zapette/internal/tools/uuid"
)

var now = time.Now().UTC()

var ExampleAlice = User{
	id:                uuid.UUID("86bffce3-3f53-4631-baf8-8530773884f3"),
	username:          "Alice",
	isAdmin:           true,
	status:            Active,
	password:          secret.NewText("alice-encrypted-password"),
	passwordChangedAt: now,
	createdAt:         now,
	createdBy:         uuid.UUID("86bffce3-3f53-4631-baf8-8530773884f3"),
}

var ExampleBob = User{
	id:                uuid.UUID("0923c86c-24b6-4b9d-9050-e82b8408edf4"),
	username:          "Bob",
	isAdmin:           false,
	status:            Active,
	password:          secret.NewText("bob-encrypted-password"),
	passwordChangedAt: now,
	createdAt:         now,
	createdBy:         ExampleAlice.id,
}

var ExampleInitializingBob = User{
	id:                uuid.UUID("0923c86c-24b6-4b9d-9050-e82b8408edf4"),
	username:          "Bob",
	isAdmin:           false,
	status:            Initializing,
	password:          secret.NewText("bob-encrypted-password"),
	passwordChangedAt: now,
	createdAt:         now,
	createdBy:         ExampleAlice.id,
}

var ExampleDeletingAlice = User{
	id:                uuid.UUID("86bffce3-3f53-4631-baf8-8530773884f3"),
	username:          "Alice",
	isAdmin:           true,
	status:            Deleting,
	password:          secret.NewText("alice-encrypted-password"),
	passwordChangedAt: now,
	createdAt:         now,
	createdBy:         uuid.UUID("86bffce3-3f53-4631-baf8-8530773884f3"),
}
