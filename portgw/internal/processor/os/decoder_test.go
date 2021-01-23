package os

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hugorut/protop/portgw/internal"
)

var (
	testFileContents = []byte(`
{
  "AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "Abu Dhabi",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  }
}
`)
)

func TestJSONDecode(t *testing.T) {
	buf := bytes.NewBuffer(testFileContents)
	ch, err := JSONDecode(buf)
	if !assert.Nil(t, err) {
		return
	}

	var ports []internal.Port
	for msg := range ch {
		assert.Nil(t, msg.Error, "received unexpected decode entry error")
		ports = append(ports, msg.Port)
	}

	assert.Len(t, ports, 2)
	assert.Contains(t, ports, internal.Port{
		Name:        "Ajman",
		Coordinates: []float64{55.5136433, 25.4052165},
		City:        "Ajman",
		Province:    "Ajman",
		Country:     "United Arab Emirates",
		Alias:       []interface{}{},
		Regions:     []interface{}{},
		Timezone:    "Asia/Dubai",
		Unlocs:      []string{"AEAJM"},
		Code:        "52000",
	})

	assert.Contains(t, ports, internal.Port{
		Name:        "Abu Dhabi",
		Coordinates: []float64{54.37, 24.47},
		City:        "Abu Dhabi",
		Province:    "Abu Z¸aby [Abu Dhabi]",
		Country:     "United Arab Emirates",
		Alias:       []interface{}{},
		Regions:     []interface{}{},
		Timezone:    "Asia/Dubai",
		Unlocs:      []string{"AEAUH"},
		Code:        "52001",
	})
}
