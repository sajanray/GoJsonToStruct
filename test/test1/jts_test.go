package test1

import (
	"encoding/json"
	"testing"
)

type Student struct {
	Name    string `stm:"name" json:"name"`
	Age     uint   `stm:"age" json:"age"`
	Mobile  string `stm:"mobile" json:"mobile"`
	Height  int    `stm:"height" json:"height"`
	Address string `stm:"address" json:"address"`
}

func TestDemo1(t *testing.T) {
	stu := Student{
		Name:    "admin",
		Age:     20,
		Mobile:  "13813141567",
		Height:  178,
		Address: "北京",
	}
	stu2 := Student{}
	bytes, _ := json.Marshal(stu)
	str := string(bytes)
	//fmt.Println("原json串：",str)

	err := json.Unmarshal([]byte(str), &stu2)
	if err == nil {
		t.Log("json转struct成功，stu2 =", stu2)
	} else {
		t.Fatal("json转struct失败", err)
	}
}

func TestDemo2(t *testing.T) {
	stu2 := Student{}
	str := `{"name":"admin","age":"20","mobile":"13813141567","height":178,"address":"北京"}`
	//fmt.Println("原json串：",str)
	err := json.Unmarshal([]byte(str), &stu2)
	if err == nil {
		t.Log("json转struct成功，stu2 =", stu2)
	} else {
		t.Fatal("json转struct失败，err =", err)
	}
}

func TestDemo3(t *testing.T) {
	stu2 := Student{}
	str := `{"name":"admin","age":"20","mobile":13813141567,"height":"178","address":"北京"}`
	//fmt.Println("原json串：",str)

	var tmp interface{}
	_ = json.Unmarshal([]byte(str), &tmp)

	m := JTStools.NewMapToStruct()
	m.Debug = false
	m.Tagkey = "stm"
	m.Transform(&stu2, tmp)
	t.Log("json转struct成功，stu2 =", stu2)
}
