package serializer

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

type TestStructs []TestStruct

type TestStruct struct {
	ReqSize    int            `compress:"1"`
	QuestionId string         `compress:"2"`
	Sizes      []int          `compress:"3"`
	P          Person         `compress:"4"`
	Map        map[string]int `compress:"5"`
}

type Person struct {
	Name string `compress:"1"`
	Age  int    `compress:"2"`
}

func hashStr() string {
	now := time.Now().UnixNano()
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%d", now)))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func TestNewSerializer2(t *testing.T) {
	var testStructs TestStructs
	for i := 0; i < 10; i++ {
		testStruct := TestStruct{
			ReqSize:    i,
			QuestionId: fmt.Sprintf("%s-%d", hashStr(), i),
			Sizes:      []int{1, 2, 3, 4, 5},
			P: Person{
				Name: "John",
				Age:  30,
			},
			Map: map[string]int{
				"key1": 1,
				"key2": 2,
				"key3": 3,
			},
		}
		testStructs = append(testStructs, testStruct)
	}

	startTime := time.Now()
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

	fmt.Println("Time taken: ", time.Since(startTime))
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
		Map: map[string]int{
			"key1": 1,
			"key2": 2,
			"key3": 3,
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
