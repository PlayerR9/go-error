package internal

import "testing"

func TestInfo(t *testing.T) {
	info := &Info{
		info: make(map[Key]any),
	}

	key1 := NewKey()
	key2 := NewKey()

	info.info[key1] = 1
	info.info[key2] = 2

	if info.info[key1] != 1 {
		t.Errorf("info[%v] = %v, want %v", key1, info.info[key1], 1)
	}
}
