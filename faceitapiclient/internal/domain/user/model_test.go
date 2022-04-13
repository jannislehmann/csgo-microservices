package user_test

import (
	"testing"

	"github.com/Cludch/csgo-microservices/faceitapiclient/internal/domain/user"
	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
	"github.com/stretchr/testify/assert"
)

var TestID = entity.NewID()

func TestNewUser(t *testing.T) {
	u := user.NewUser(TestID)
	assert.NotNil(t, u)
	assert.Equal(t, TestID, u.ID)
}
