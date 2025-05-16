package domain

type (
	Gender int

	RacquetSportType int

	BestHand  int
	CourtSide int

	RacquetMatchType  int
	RacquetMatchState int
)

const (
	DefaultRacquetRating = 1000
)

const (
	UnknownGender Gender = iota
	Male
	Female
)

const (
	UnknownRacquetSportType RacquetSportType = iota
	Tennis
	Padel
)

const (
	UnknownBestHand BestHand = iota
	Left
	Right
	BothHands
)

const (
	UnknownCourtSide CourtSide = iota
	Backhand
	Forehand
	BothSides
)

const (
	UnknownRacquetMatchType RacquetMatchType = iota
	Singles
	Doubles
)

const (
	UnknownRacquetMatchState RacquetMatchState = iota
	RacquetMatchCreated
	RacquetMatchScheduled
	RacquetMatchInProgress
	RacquetMatchResultSubmission
	RacquetMatchCompleted
	RacquetMatchNonConfirmed
	RacquetMatchCanceled
)
