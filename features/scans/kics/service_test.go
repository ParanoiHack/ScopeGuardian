package kics

import (
	"fmt"
	"scope-guardian/loader"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("Should run kics", func(t *testing.T) {
		service := newKicsService(loader.Kics{"./WebGoat", "Dockerfile"})
		ok, err := service.Start()

		fmt.Println(ok)
		fmt.Println(err)
	})
}
