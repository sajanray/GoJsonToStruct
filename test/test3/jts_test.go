package demo03_test

import (
	"encoding/json"
	JTStools "github.com/sajanray/GoJsonToStruct"
	"testing"
)

type Student struct {
	Name    string `stm:"name" json:"name"`
	Age     uint   `stm:"age" json:"age"`
	Mobile  string `stm:"mobile" json:"mobile"`
	Height  int    `stm:"height" json:"height"`
	Address string `stm:"address" json:"address"`
	School  School `stm:"school" json:"school"`
}

type School struct {
	Name string    `json:"name" stm:"name"`
	Founded string `json:"founded" stm:"founded"`
	Age uint       `json:"age" stm:"age"`
	Subject []Subject `json:"subject" stm:"subject"`
	Subject2 []*Subject `json:"subject2" stm:"subject2"`
	Subject3 map[string]Subject `json:"subject3" stm:"subject3"`
}

type Subject struct {
	Name string
	Score float32
	Pass bool
}

func TestDemo1(t *testing.T) {
	stu := Student{
		Name: "admin",
		Age: 20,
		Mobile: "13813141567",
		Height: 178,
		Address: "北京",
		School: School{
			Name: "某某大学",
			Founded: "1970-10-01",
			Age: 70,
			Subject: []Subject{
				{Name: "语文",Score: 90,Pass: true},
				{Name: "美术",Score: 50,Pass: false},
				{Name: "数学",Score: 90,Pass: true},
			},
			Subject2: []*Subject{
				&Subject{Name: "语文",Score: 90,Pass: true},
				&Subject{Name: "美术",Score: 50,Pass: false},
				{Name: "数学",Score: 90,Pass: true},
			},
			Subject3: map[string]Subject{
				"0":{Name: "语文",Score: 90,Pass: true},
				"2":{Name: "美术",Score: 50,Pass: false},
				"5":{Name: "数学",Score: 90,Pass: true},
			},
		},
	}
	stu2 := Student{}
	bytes,_ := json.Marshal(stu)
	str := string(bytes)
	//fmt.Println("原json串：",str)

	err := json.Unmarshal([]byte(str) , &stu2)
	if err==nil{
		t.Log("json转struct成功，stu2 =",stu2)
	} else {
		t.Fatal("json转struct失败",err)
	}
}

func TestDemo2(t *testing.T) {
	stu2 := Student{}
	str := `{
  "name": "admin",
  "age": 20,
  "mobile": "13813141567",
  "height": 178,
  "address": "北京",
  "school": {
    "name": "某某大学",
    "founded": "1970-10-01",
    "age": 70,
    "subject": [
      {
        "Name": "语文",
        "Score": 90,
        "Pass": true
      },
      {
        "Name": "美术",
        "Score": 50,
        "Pass": false
      },
      {
        "Name": "数学",
        "Score": 90,
        "Pass": true
      }
    ],
    "subject2": [
      {
        "Name": "语文",
        "Score": 90,
        "Pass": true
      },
      {
        "Name": "美术",
        "Score": 50,
        "Pass": false
      },
      {
        "Name": "数学",
        "Score": 90,
        "Pass": true
      }
    ],
    "subject3": {
      "0": {
        "Name": "语文",
        "Score": 90,
        "Pass": true
      },
      "2": {
        "Name": "美术",
        "Score": 50,
        "Pass": false
      },
      "5": {
        "Name": "数学",
        "Score": 90,
        "Pass": true
      }
    }
  }
}`
	//fmt.Println("原json串：",str)
	err := json.Unmarshal([]byte(str) , &stu2)
	if err==nil{
		t.Log("json转struct成功，stu2 =",stu2)
	} else {
		t.Fatal("json转struct失败，err =",err)
	}
}

func TestDemo3(t *testing.T) {
	stu2 := Student{}
	str := `
{
  "name": "admin",
  "age": 20,
  "mobile": "13813141567",
  "height": 178,
  "address": "北京",
  "school": {
    "name": "某某大学",
    "founded": "1970-10-01",
    "age": 70,
    "subject": [
      {
        "Name": "语文",
        "Score": 90,
        "Pass": true
      },
      {
        "Name": "美术",
        "Score": 50,
        "Pass": false
      },
      {
        "Name": "数学",
        "Score": 90,
        "Pass": true
      }
    ],
    "subject2": [
      {
        "Name": "语文",
        "Score": 90,
        "Pass": true
      },
      {
        "Name": "美术",
        "Score": 50,
        "Pass": false
      },
      {
        "Name": "数学",
        "Score": 90,
        "Pass": true
      }
    ],
    "subject3": {
      "1":{
        "Name": "语文",
        "Score": 90,
        "Pass": true
      },
      "2":{
        "Name": "美术",
        "Score": 50,
        "Pass": false
      },
      "3":{
        "Name": "数学",
        "Score": 90,
        "Pass": true
      }
    }
  }
}
`
	//fmt.Println("原json串：",str)

	var tmp interface{}
	_ = json.Unmarshal([]byte(str), &tmp)

	m := JTStools.NewMapToStruct()
	m.Debug  = false
	m.Tagkey = "stm"
	m.Transform(&stu2 , tmp)

	t.Log("json转struct成功，stu2 =",stu2)
}