package reversebit

import "testing"

func TestReverseBits(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{name: "用例1", input: 43261596, expected: 964176192},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reverseBits(tt.input)
			if result != tt.expected {
				t.Errorf("reverseBits(%d) = %d, expected %d\n二进制: input=%032b, result=%032b, expected=%032b",
					tt.input, result, tt.expected,
					uint32(tt.input), uint32(result), uint32(tt.expected))
			}
		})
	}
}
