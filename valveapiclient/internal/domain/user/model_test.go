package user_test

import (
	"testing"

	"github.com/Cludch/csgo-microservices/valveapiclient/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

const TestId = uint64(1)
const TestApiKey = "key"
const TestAuthCode = "auth"
const TestShareCode = "CSGO-Y4DVh-amkvh-OyBrh-SyMHN-2SvPB"

func TestNewUser(t *testing.T) {
	u := user.NewUser(TestId)
	assert.NotNil(t, u)
	assert.Equal(t, TestId, u.ID)
}
