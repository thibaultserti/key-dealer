package test

import "testing"

func TestSum(t *testing.T) {

	t.Run("Test", func(t *testing.T) {

		got := 1
		want := 1

		if got != want {
			t.Errorf("got %d want %d given", got, want)
		}
	})

}
