package router

import "testing"

func TestValidatePath_Success(t *testing.T) {
	testTable := []struct {
		Path string
	}{
		{"/p1/p1/p1/"},
		{"/p1/p1/p1/123/3213/"},
		{"/p1/!/@/#/$/%/^/&/*/"},
		{"/:asd/:dasd?a/asd?/"},
	}
	for testCase, test := range testTable {
		err := validatePath(test.Path)
		if err != nil {
			t.Errorf("Test %d: Expected no error but go %v", testCase, err)
			continue
		}
	}
}
func TestValidatePath_Failure(t *testing.T) {
	testTable := []struct {
		Path string
	}{
		{"/p1/p1/p1"},
		{"p1/p1/p1/123/3213/"},
		{"/p1/ !/@/#/$/%/^/&/*/"},
		{":asd/:dasd?a/asd?//"},
	}
	for testCase, test := range testTable {
		err := validatePath(test.Path)
		if err == nil {
			t.Errorf("Test %d: Expected error but got nil", testCase)
			continue
		}
	}
}
