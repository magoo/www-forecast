package models

type Votes []Vote

type Vote struct {
	Index int
	Votes int
}

func (v Votes) Len() int		{ return len(v) }
func (v Votes) Swap(i, j int)		{ v[i], v[j] = v[j], v[i] }
func (v Votes) Less(i, j int) bool { return v[i].Votes < v[j].Votes}
