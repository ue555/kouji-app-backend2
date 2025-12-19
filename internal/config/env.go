package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var loadOnce sync.Once

// LoadEnv loads variables from a .env file exactly once, if the file exists.
func LoadEnv(files ...string) {
	loadOnce.Do(func() {
		paths := files
		if len(paths) == 0 {
			paths = []string{".env"}
		}

		var existing []string
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				existing = append(existing, path)
			}
		}

		if len(existing) == 0 {
			return
		}

		if err := godotenv.Load(existing...); err != nil {
			log.Printf("warning: cannot load env files %v: %v", existing, err)
		}
	})
}
