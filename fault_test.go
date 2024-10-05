package errs

import "testing"

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

func (m *MockType) As(target any) bool {
	return false
}

func TestAsSuccess(t *testing.T) {
	const (
		ExpectedStr string = "(mock) some error"
	)

	err := &MockType{}

	var target *MockType

	ok := As(err, target)
	if !ok {
		t.Fatalf("As() = %v, want %v", ok, true)
	}

	// if target == nil {
	// 	t.Fatalf("As() = %v, want %v", target, true)
	// }

	if target.Error() != ExpectedStr {
		t.Errorf("As() = %v, want %v", target, ExpectedStr)
	}
}

func TestAs(t *testing.T) {
	const (
		ExpectedStr string = "(mock) some error"
	)

	err := New(MockCode1, "some error")

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
