package input

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
)

type Item struct {
	values []float32
	index  int
}

func (i Item) Title() string {
	return fmt.Sprint(i.values[i.index])
}

func (i Item) Description() string { return "" }

func (i Item) FilterValue() string {
	return fmt.Sprint(i.index)
}

func NewItems(values []float32) []list.Item {
	items := make([]list.Item, len(values))
	for i := range values {
		items[i] = Item{
			values: values,
			index:  i,
		}
	}

	return items
}
