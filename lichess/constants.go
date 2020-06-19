package lichess

type ChatRoom string

const (
	ChatRoomPlayer    = "player"
	ChatRoomSpectator = "spectator"
)

type Color string

const (
	ColorRandom Color = "random"
	ColorWhite  Color = "white"
	ColorBlack  Color = "black"
)

const (
	VariantKeyStandard      = "standard"
	VariantKeyChess960      = "chess960"
	VariantKeyCrazyhouse    = "crazyhouse"
	VariantKeyAntichess     = "antichess"
	VariantKeyAtomic        = "atomic"
	VariantKeyHorde         = "horde"
	VariantKeyKingOfTheHill = "kingOfTheHill"
	VariantKeyRacingKings   = "racingKings"
	VariantKeyThreeCheck    = "threeCheck"
)

const (
	GameStateStarted = "started"
	GameStateResign  = "resign"
)
