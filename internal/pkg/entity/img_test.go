package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_test(t *testing.T) {
	testTable := []struct {
		name           string
		fileName       string
		expectedResult string
	}{
		{
			"Simple test",
			"test.txt",
			"txt",
		},
		{
			"Multiple dots test",
			"data.text.test.ds",
			"ds",
		},
		{
			"Without ext test",
			"binaryFile",
			"",
		},
		{
			"Invalid file name test",
			"InvalidFile.",
			"",
		},
		{
			"Hidden file with ext name test",
			".HiddenFile.config",
			"config",
		},
		{
			"Hidden file without ext name test",
			".HiddenFile",
			"",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			result := resolveExt(testCase.fileName)
			assert.Equal(t, testCase.expectedResult, result)
		})

	}
}
