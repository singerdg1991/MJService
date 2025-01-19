package router

import "testing"

func TestRouterInit(t *testing.T) {
	t.Run("When router is nil", func(t *testing.T) {
		r := Init()
		if r == nil {
			t.Error("Router can not be nil")
		}
	})

	t.Run("Check result value with Default variable", func(t *testing.T) {
		r := Init()
		if r != Default {
			t.Error("Default is not equal to result")
		}
	})
}
