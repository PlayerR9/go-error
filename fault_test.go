package errs

import "testing"

////////////////////////////////////////////////////////////////////

type MockCode int

const (
	MockCode1 MockCode = iota
)

func (m MockCode) String() string {
	return "mock"
}

type MockType struct {
}

func (m MockType) Error() string {
	return "(mock) some error"
}

func (m MockType) Is(target Fault) bool {
	return m.Error() == target.Error()
}

func TestAs(t *testing.T) {
	const (
		ExpectedStr string = "(mock) some error"
	)

	err := &MockType{}

	var target *MockType

	ok := As(err, &target)
	if !ok {
		t.Fatalf("As() = %v, want %v", ok, true)
	}

	if target == nil {
		t.Fatalf("As() = %v, want %v", target, true)
	}

	if target.Error() != ExpectedStr {
		t.Errorf("As() = %v, want %v", target, ExpectedStr)
	}
}

func TestIs(t *testing.T) {
	const (
		ExpectedStr string = "(mock) some error"
	)

	err := &MockType{}

	ok := Is(err, &MockType{})
	if !ok {
		t.Fatalf("Is() = %v, want %v", ok, true)
	}

	if err.Error() != ExpectedStr {
		t.Errorf("Is() = %v, want %v", err, ExpectedStr)
	}
}
