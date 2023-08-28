package patch

type ValidSentinel interface {
	IsValid() bool
	SetValid(bool)
}
