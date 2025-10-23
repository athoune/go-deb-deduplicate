package bytes

type Id string

type ById []Id

func (b ById) Len() int {
	return len(b)
}

func (b ById) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ById) Less(i, j int) bool {
	return b[i] < b[j]
}
