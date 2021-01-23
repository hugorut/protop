package os_test

import (
	"io"
	"os"
	"sync"
	"syscall"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/hugorut/protop/portgw/internal"
	. "github.com/hugorut/protop/portgw/internal/processor/os"
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

type stubStore struct {
	mu *sync.Mutex
	wg *sync.WaitGroup

	ports []internal.Port
}

func (s *stubStore) Store(port internal.Port) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wg.Done()
	s.ports = append(s.ports, port)

	return nil
}

func TestProcessor(t *testing.T) {
	t.Run("Process", func(t *testing.T) {
		logger, _ := test.NewNullLogger()

		wg := &sync.WaitGroup{}
		store := &stubStore{
			mu: &sync.Mutex{},
			wg: wg,
		}
		wg.Add(2)

		f := os.NewFile(uintptr(syscall.Stdout), "test")

		port1 := internal.Port{Name: "a"}
		port2 := internal.Port{Name: "b"}

		loc := "test/loc"

		p := Processor{
			Processors: 2,
			Open: func(name string) (file *os.File, err error) {
				assert.Equal(t, loc, name)
				return f, nil
			},
			Decode: func(input io.Reader) (entries chan DecodedEntry, err error) {
				assert.Equal(t, f, input)

				ch := make(chan DecodedEntry)

				go func() {
					ch <- DecodedEntry{Port: port1}
					ch <- DecodedEntry{Port: port2}
					close(ch)
				}()

				return ch, nil
			},
			Store:  store,
			Logger: logger,
		}

		_, err := p.Process(loc)
		if !assert.Nil(t, err) {
			return
		}

		wg.Wait()
		assert.Len(t, store.ports, 2)
		assert.Contains(t, store.ports, port1)
		assert.Contains(t, store.ports, port2)
	})
}
