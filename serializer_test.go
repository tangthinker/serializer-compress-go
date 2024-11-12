package serializer_compress_go

import (
	"fmt"
	"testing"
)

type TestStructs []TestStruct

type TestStruct struct {
	ReqSize    int    `compress:"1"`
	QuestionId string `compress:"2"`
	Sizes      []int  `compress:"3"`
	P          Person `compress:"4"`
}

type Person struct {
	Name string `compress:"1"`
	Age  int    `compress:"2"`
}

func TestNewSerializer2(t *testing.T) {
	var testStructs TestStructs

	for i := 0; i < 100; i++ {
		testStruct := TestStruct{
			ReqSize:    i,
			QuestionId: fmt.Sprintf("1sfdkio1klmdskfmc1o%d", i),
			Sizes:      []int{1, 2, 3, 4, 5},
		}
		testStructs = append(testStructs, testStruct)
	}

	s := NewSerializer()

	data, err := s.Encode(testStructs)
	if err != nil {
		t.Errorf("Error encoding data: %v", err)
		return
	}

	fmt.Println("encoded length : ", len(data))

	var target TestStructs
	err = s.Decode(data, &target)
	if err != nil {
		t.Errorf("Error decoding data: %v", err)
		return
	}

	fmt.Println("Decoded data: ", target)
}

func TestNewSerializer(t *testing.T) {

	testStruct := TestStruct{
		ReqSize:    10,
		QuestionId: "1sfdkio1klmdskfmc1o",
		Sizes:      []int{1, 2, 3, 4, 5},
		P: Person{
			Name: "John",
			Age:  30,
		},
	}

	s := NewSerializer()
	data, err := s.Encode(testStruct)

	if err != nil {
		t.Errorf("Error encoding data: %v", err)
		return
	}

	fmt.Println("encoded length : ", len(data))

	for _, b := range data {
		fmt.Printf("%08b\t", b)
	}

	fmt.Println()

	var target TestStruct

	err = s.Decode(data, &target)

	if err != nil {
		t.Errorf("Error decoding data: %v", err)
		return
	}

	fmt.Println("Decoded data: ", target)

}
