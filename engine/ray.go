package juicer

type RayDirection uint8

const (
	Nort RayDirection = iota // north
	NoEa                     // north-east
	East                     // east
	SoEa                     // south-east
	Sout                     // south
	SoWe                     // south-west
	West                     // west
	NoWe                     // north-west
)
