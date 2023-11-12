package juicer

type RayDirection uint8

const (
	N  RayDirection = iota // north
	NE                     // north-east
	E                      // east
	SE                     // south-east
	S                      // south
	SW                     // south-west
	W                      // west
	NW                     // north-west
)
