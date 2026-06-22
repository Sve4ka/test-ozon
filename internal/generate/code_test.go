package generate

import "testing"

func TestGenerator_Generate_Length(t *testing.T) {
	generator := NewRandomGenerator()

	code, err := generator.Generate()
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	if len(code) != CodeLength {
		t.Fatalf("expected code length %d, got %d", CodeLength, len(code))
	}
}

func TestGenerator_Generate_AllowedCharacters(t *testing.T) {
	generator := NewRandomGenerator()

	for i := 0; i < 1000; i++ {
		code, err := generator.Generate()
		if err != nil {
			t.Fatalf("Generate() returned error: %v", err)
		}

		for _, ch := range string(code) {
			if !isAllowedChar(ch) {
				t.Fatalf("generated code contains invalid char %q in %q", ch, code)
			}
		}
	}
}

func isAllowedChar(ch rune) bool {
	return ch >= 'a' && ch <= 'z' ||
		ch >= 'A' && ch <= 'Z' ||
		ch >= '0' && ch <= '9' ||
		ch == '_'
}

func TestIsValidCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "valid code",
			code: "aB12_cdEF3",
			want: true,
		},
		{
			name: "too short",
			code: "abc",
			want: false,
		},
		{
			name: "too long",
			code: "aB12_cdEF30",
			want: false,
		},
		{
			name: "dash is invalid",
			code: "aB12-cdEF3",
			want: false,
		},
		{
			name: "cyrillic is invalid",
			code: "абвгдежзик",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidCode(tt.code)
			if got != tt.want {
				t.Fatalf("IsValidCode(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}
