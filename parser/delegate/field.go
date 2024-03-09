package delegate

// Field interface defines the component interface.
type Field interface {
	GetName() string
	GetPattern() string
	GetStrat() parseStrategy
}

// NewSLF returns a pointer to a new instance of singleLineField.
func NewSLF(name string, pattern string, strat parseStrategy) *singleLineField {
	return &singleLineField{
		name:    name,
		pattern: pattern,
		strat:   strat,
	}
}

// singleLineField acts a concrete component.
// It is a concrete implementation of the Field interface.
type singleLineField struct {
	name    string
	pattern string
	strat   parseStrategy
}

func (slf singleLineField) GetName() string {
	return slf.name
}

func (slf singleLineField) GetPattern() string {
	return slf.pattern
}

func (slf singleLineField) GetStrat() parseStrategy {
	return slf.strat
}

// NewMLF returns a pointer to a new instance of multiLineField.
func NewMLF(name string, pattern string, strat parseStrategy, isBeginSeq func(line string) bool, isEndSeq func(line string) bool, clean func([]string) interface{}) *multiLineField {
	return &multiLineField{
		Field:           NewSLF(name, pattern, strat),
		isBeginSequence: isBeginSeq,
		isEndSequence:   isEndSeq,
		cleanSequence:   clean,
	}
}

// multiLineField acts a decorator.
// It embeds the Field interface and adds additional functionality.
type multiLineField struct {
	Field
	isBeginSequence func(line string) bool
	isEndSequence   func(line string) bool
	cleanSequence   func(sequence []string) interface{}
}

func (mlf multiLineField) GetName() string {
	return mlf.Field.GetName()
}

func (mlf multiLineField) GetPattern() string {
	return mlf.Field.GetPattern()
}

func (mlf multiLineField) GetStrat() parseStrategy {
	return mlf.Field.GetStrat()
}

// NewSLMF returns a pointer to a new instance of sameLineMultiField.
func NewSLMF(name string, pattern string, strat parseStrategy, additionalPatterns []string) *sameLineMultiField {
	return &sameLineMultiField{
		Field:              NewSLF(name, pattern, strat),
		additionalPatterns: additionalPatterns,
	}
}

// sameLineMultiField acts as a decorator.
// It embeds the Field interface and adds additional functionality.
type sameLineMultiField struct {
	Field
	additionalPatterns []string
}

func (slmf sameLineMultiField) GetName() string {
	return slmf.Field.GetName()
}

func (slmf sameLineMultiField) GetPattern() string {
	return slmf.Field.GetPattern()
}

func (slmf sameLineMultiField) GetStrat() parseStrategy {
	return slmf.Field.GetStrat()
}
