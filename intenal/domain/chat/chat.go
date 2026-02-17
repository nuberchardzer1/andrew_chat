package chat

type ChatFlags int
const (
	GroupFlag = 1 << iota
	ProtectedFlag 
)

// represent selectable server
type Chat struct {
	Name string
	Flags ChatFlags
}

//Implements item.Item (bubbles)
func (c Chat) FilterValue() string{
	return c.Name
}

func (c Chat) Title() string {
    return c.Name
}

func (c Chat) Description() string {
    desc := ""
	if c.Flags & GroupFlag > 0{
		desc += "GROUP"
	}else{
		desc += "PRIVATE"
	}

	if c.Flags & ProtectedFlag > 0{
		desc += "| PROTECTED"
	}

	return desc
}