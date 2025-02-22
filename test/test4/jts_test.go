package test4

import (
	JTStools "github.com/sajanray/GoJsonToStruct"
	"testing"
)

type Student struct {
	Name   string `stm:"name" json:"name"`
	Age    int    `stm:"age" json:"age"`
	Mobile string `stm:"mobile" json:"mobile"`
}

func TestString2Int(t *testing.T) {
	stu2 := Student{}
	str := `
{
  "name": "admin",
  "age": " \n",
  "mobile": "13813141567"
}
`
	JTStools.NewMapToStruct().Transform(&stu2, str)
	t.Log("json转struct成功，stu2 =", stu2)
}
