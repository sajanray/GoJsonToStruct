# GoJsonToStruct
json串格式化进结构体，此工具会把json里的数据类型自动匹配转换成结构体内部对应的数据类型。

#安装
使用go mod，直接import 
```gotemplate
import (
    "github.com/sajanray/GoJsonToStruct"
)
```
使用go get，拷贝JTStools.go放进你的项目里即可
```gotemplate
go get github.com/sajanray/GoJsonToStruct
```

#应用
定义json串对用的结构体
```gotemplate
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
```

调用test函数进行转换测试
```gotemplate
func test() {
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

    //先把json串Unmarshal
    var tmp interface{}
    err := json.Unmarshal([]byte(str), &tmp)
    if err != nil {
        fmt.Println("json.Unmarshal error:",err.Error())
        return
    }
    
    //开始转换
    m := JTStools.NewMapToStruct()
    m.Transform(&stu2 , tmp)
    if m.Success {
        fmt.Println("转换成功")
        fmt.Println(stu2)
    } else {
        fmt.Println("转换失败")
    }
}

//转换成功
//{admin 20 13813141567 178 北京 {某某大学 1970-10-01 70 [{语文 90 true} {美术 50 false} {数学 90 true}] [0xc0000046c0 0xc000004700 0xc000004740] map[1:{语文 90 true} 2:{美术 50 false} 3:{数学 90 true}]}}
```
done
