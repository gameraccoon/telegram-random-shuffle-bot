package dialog

type Variant struct {
	Id   string
	Data string // additional id or some info that will be appended to a button/link
	Text string // text visible to the user
}

type Dialog struct {
	Id       string
	Text     string // dialog header visible to the user
	Variants []Variant
}
