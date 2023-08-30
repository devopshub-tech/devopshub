package filesystem

import (
	"fmt"
	"os"
)

func NewDirectory(name string) error {
	err := os.Mkdir(name, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating directory %s: %w", name, err)
	}
	return nil
}
