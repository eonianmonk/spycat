package spycat

import "fmt"

type ComletionStatus string

const (
	Incomplete ComletionStatus = "incomplete"
	Complete   ComletionStatus = "complete"
)

func (cs ComletionStatus) Validate() error {
	switch cs {
	case Incomplete:
		return nil
	case Complete:
		return nil
	default:
		return fmt.Errorf("invalid completion status: %s", cs)
	}
}
