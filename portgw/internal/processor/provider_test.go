package processor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/hugorut/protop/portgw/internal/processor"
)

type stubFileProcessor struct{}

func (s stubFileProcessor) Process(string) (int, error) { panic("implement me") }

//go:generate mockgen -destination mocks/provider.go -source ./provider.go
func TestProvider(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		type args struct {
			name string
		}

		tests := []struct {
			name     string
			receiver Provider
			args     args
			want     FileProcessor
			wantErr  bool
		}{
			{
				name:     "test valid processor name returns processor",
				receiver: Provider{"valid": stubFileProcessor{}},
				args:     args{name: "valid"},
				want:     stubFileProcessor{},
				wantErr:  false,
			},
			{
				name:     "test invalid processor name returns error",
				receiver: Provider{"valid": stubFileProcessor{}},
				args:     args{name: "invalid"},
				wantErr:  true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := tt.receiver.Get(tt.args.name)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.wantErr {
					assert.EqualValues(t, ErrorInvalidProvider, err)
				}

				assert.EqualValues(t, tt.want, got, "unexpected processor provided")
			})
		}
	})

}
