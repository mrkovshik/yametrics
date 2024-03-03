package storage

import "strconv"

type (
	counter struct {
		name  string
		value int64
	}
)

func (c counter) Update(s IStorage) error {
	return s.UpdateCounter(c)

}

func NewCounter(name, value string) (counter, error) {
	intValue, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return counter{}, err
	}
	return counter{name: name, value: intValue}, err
}
