package models

type StatsChanVariable struct {
	Varname string
	Value   any  // updated value to publish (int , float, etc)
	Inc     bool // Treat the count as an increment as opposite to the final value.
}
