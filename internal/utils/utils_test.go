package utils

import (
	"fmt"
	"testing"
)

type TTest struct {
	Name string
}

func TestGobMarshal(t *testing.T) {
	var tt = TTest{
		"M!!",
	}

	b := GobMarshal(&tt)

	br, err := GobUnmarshal[*TTest](b)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(br.Name)
}
