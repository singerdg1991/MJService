package migrator

import (
	"github.com/hoitek/Maja-Service/cmd/migrator/ports"
)

func New[T ports.Migrator](driver T) (T, error) {
	return driver, nil
}
