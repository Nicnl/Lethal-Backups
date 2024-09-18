package save_decoder

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRead(t *testing.T) {
	saveInfo, err := Read(v64_slot1)
	assert.NoError(t, err)

	fmt.Printf("%+v\n", saveInfo)

	assert.Equal(t, 207, saveInfo.GroupCredits.Value)
	assert.Equal(t, 27, saveInfo.Stats_DaysSpent.Value)
	assert.Equal(t, 0, saveInfo.QuotaFulfilled.Value)
	assert.Equal(t, 1350, saveInfo.ProfitQuota.Value)
	assert.Equal(t, 5, saveInfo.CurrentPlanetID.Value)
}
