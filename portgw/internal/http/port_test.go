package http_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	. "github.com/hugorut/protop/portgw/internal/http"
	mock_http "github.com/hugorut/protop/portgw/internal/http/mocks"
	mock_processor "github.com/hugorut/protop/portgw/internal/processor/mocks"
)


//go:generate mockgen -destination mocks/port.go -source ./port.go
func TestHandler(t *testing.T) {
	logger, _ := test.NewNullLogger()
	testFileLocation := "/test/location"
	body := []byte(fmt.Sprintf(`{"location": "%s"}`, testFileLocation) )

	t.Run("ProcessFile", func(t *testing.T) {
		t.Run("With valid provider", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pp := mock_http.NewMockProcessorProvider(ctrl)

			handler := Handler{
				Logger:            logger,
				ProcessorProvider: pp,
			}

			provider := "os"
			w := httptest.NewRecorder()
			buf := bytes.NewBuffer(body)

			r := httptest.NewRequest(http.MethodPost, "/ports/file/os/upload/", buf)
			r = mux.SetURLVars(r, map[string]string{"provider": provider})

			fp := mock_processor.NewMockFileProcessor(ctrl)
			pp.EXPECT().Get(provider).Return(fp, nil)

			id := 2
			fp.EXPECT().Process(testFileLocation).Return(id, nil)
			handler.ProcessFile(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, fmt.Sprintf(`{"id": %d}`, id), w.Body.String())
		})

		t.Run("With invalid provider", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pp := mock_http.NewMockProcessorProvider(ctrl)

			handler := Handler{
				Logger:            logger,
				ProcessorProvider: pp,
			}

			provider := "unsup"
			w := httptest.NewRecorder()
			buf := bytes.NewBuffer(body)

			r := httptest.NewRequest(http.MethodPost, "/ports/file/os/upload/", buf)
			r = mux.SetURLVars(r, map[string]string{"provider": provider})

			pp.EXPECT().Get(provider).Return(nil, errors.New("unsupported provider"))
			handler.ProcessFile(w, r)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, fmt.Sprintf(`{"error": "invalid provider %s given"}`, provider), w.Body.String())
		})
	})
}
