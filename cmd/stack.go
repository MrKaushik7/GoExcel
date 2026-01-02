package cmd

type StringStack struct {
	items []string
}

func (s* StringStack) Top() string {
	return s.items[len(s.items) - 1]
}

func (s* StringStack) Push(item string) {
	s.items = append(s.items, item)
}

func (s* StringStack) Pop() (string, bool) {
	if len(s.items) == 0 {
		return "0", false
	}
	item := s.items[len(s.items) - 1]
	s.items = s.items[:len(s.items) - 1]
	return item, true
}

func (s* StringStack) IsEmpty() bool {
	return (len(s.items) == 0)
}