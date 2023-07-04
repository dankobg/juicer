package engine

type Player struct {
	Name  string
	Color Color
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:  name,
		Color: White,
	}
}

func (p Player) String() string {
	return p.Name
}
