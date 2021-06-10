package Stack

type Stack struct {
	Values []string
}

func (s *Stack) IsEmpty() bool {
	return len((*s).Values) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(str string) {
	(*s).Values = append((*s).Values, str) // Simply append the new value to the end of the stack
}

func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len((*s).Values) - 1 // Get the index of the top most element.
		element := ((*s).Values)[index] // Index into the slice and obtain the element.
		(*s).Values = (*s).Values[:index] // Remove it from the stack by slicing it off.
		return element, true
	}
}