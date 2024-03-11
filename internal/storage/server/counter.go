package server

type (
	Counter struct {
		name  string
		value int64
	}
)

func (c Counter) Update(s IStorage) error {
	return s.UpdateCounter(c)

}

func NewCounter(name string, value int64) Counter {
	return Counter{name: name, value: value}
}
