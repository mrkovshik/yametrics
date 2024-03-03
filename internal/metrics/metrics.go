package metrics

import "github.com/mrkovshik/yametrics/internal/storage"

type metric interface {
	Update(storage.IStorage) error
}
