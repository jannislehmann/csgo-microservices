package entity_test

import (
	"testing"

	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	output := entity.NewID()
	assert.NotNil(t, output)
}

func TestStringToID(t *testing.T) {
	input := "3fa4f16f-a5fd-4cae-a943-e47ce13f49b6"
	output, err := entity.StringToID(input)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input, output.String())
}

func TestStringToIDInvalidUUID(t *testing.T) {
	input := "invalid"
	_, err := entity.StringToID(input)
	assert.NotNil(t, err)
}
