package types

type Position int8

const (
	PositionLeft Position = 1 << iota
	PositionRight
	PositionTop
	PositionBot
)

const (
	PositionTopLeft  Position = PositionLeft | PositionTop
	PositionBotLeft  Position = PositionLeft | PositionBot
	PositionTopRight Position = PositionRight | PositionTop
	PositionBotRight Position = PositionRight | PositionBot
	PositionSentinel Position = -1
)

func (p Position) IsTop() bool {
	return p&PositionTop > 0
}

func (p Position) IsBot() bool {
	return p&PositionBot > 0
}

func (p Position) IsLeft() bool {
	return p&PositionLeft > 0
}

func (p Position) IsRight() bool {
	return p&PositionRight > 0
}

func (p Position) SetLeft() Position {
	return (p &^ PositionRight) | PositionLeft
}

func (p Position) SetRight() Position {
	return (p &^ PositionLeft) | PositionRight
}

func (p Position) SetTop() Position {
	return (p &^ PositionBot) | PositionTop
}

func (p Position) SetBot() Position {
	return (p &^ PositionTop) | PositionBot
}
