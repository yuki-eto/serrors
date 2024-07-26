package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	err := create("test_error")
	s := fmt.Sprintf("%+v", Wrap(err, "wrap_error"))
	if !strings.Contains(s, "testing.tRunner") {
		t.Errorf("invalid stack trace: %s", s)
	}
}

func TestAsAppError(t *testing.T) {
	err := errors.New("new_error")
	werr := Wrap(err)

	t.Run("ng", func(t *testing.T) {
		ae := AsAppError(err)
		if ae != nil {
			t.Errorf("expected nil, got %v", ae)
		}
	})
	t.Run("ok", func(t *testing.T) {
		ae := AsAppError(werr)
		if ae == nil {
			t.Errorf("expected AppError, got nil")
		}
	})
}
