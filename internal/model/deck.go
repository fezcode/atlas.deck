package model

type Pad struct {
	Key     string `piml:"key"`
	Label   string `piml:"label"`
	Command string `piml:"cmd"`
	Color   string `piml:"color"` // Optional: cyan, gold, red, etc.
}

type Deck struct {
	Name    string `piml:"name"`
	Version string `piml:"version"`
	Pads    []Pad  `piml:"pads"`
}
