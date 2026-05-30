package entity

type Question struct {
	ID              uint
	Text            string
	PossibleAnswers []PossibleAnswer
	CorrectAnswerID uint
	Difficulity     QuestionDifficulity
	CategoryID      uint
}

type QuestionDifficulity uint8

const (
	QuestionDifficulityEasy QuestionDifficulity = iota + 1
	QuestionDifficulityMedium
	QuestionDifficulityHard
)

func (q QuestionDifficulity) IsValid() bool {
	if q >= QuestionDifficulityEasy && q <= QuestionDifficulityHard {
		return true
	}

	return false

}

type PossibleAnswer struct {
	ID     uint
	Text   string
	Chioce PossibleAnswerChoice
}

type PossibleAnswerChoice uint8

func (p PossibleAnswerChoice) IsValid() bool {
	if p >= PossibleAnswerChoiceA && p <= PossibleAnswerChoiceD {
		return true
	}

	return false
}

const (
	PossibleAnswerChoiceA PossibleAnswerChoice = iota + 1
	PossibleAnswerChoiceB
	PossibleAnswerChoiceC
	PossibleAnswerChoiceD
)
