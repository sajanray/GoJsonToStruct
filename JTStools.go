/*
 @project:JsonToStruct
 @author:sajanray@163.com
 @date:2020/6/30
 @note:
*/
package JTStools

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type MapToStruct struct {
	Debug bool        //调试模式
	Success bool      //是否转换成功
	Tagkey string     //结构体标签名

	structTypeOf  reflect.Type
	structTofElem reflect.Type
	structValueOf reflect.Value
	structVofElem reflect.Value

	sourceMapData interface{}
}

func NewMapToStruct() *MapToStruct{
	m := &MapToStruct{}
	m.Tagkey = "json"
	m.Debug = false
	return m
}

func (m *MapToStruct) cloneMapToStruct() *MapToStruct{
	n := &MapToStruct{}
	n.Tagkey = m.Tagkey
	n.Debug = m.Debug
	return n
}

func (m *MapToStruct) getMapValue(i int) (mapVal interface{},ok bool){
	//取tag名
	tagName := m.structTofElem.Field(i).Tag.Get(m.Tagkey)
	if tagName != "" {
		mapVal, ok = m.sourceMapData.(map[string]interface{})[tagName]  //取map对应结构体tagName的值
	} else {
		//todo 结构体字段tag名称这里没有对字段名做任何转换和结构体字段名保持一致（可以对字段名进行多次转换后在尝试取map中的值）
		tagName = m.structTofElem.Field(i).Name
		mapVal, ok = m.sourceMapData.(map[string]interface{})[tagName]
	}
	return mapVal,ok
}

//MapToStruct  把map映射到结构体
func (m *MapToStruct) Transform(destStructData interface{} , sourceMap interface{}) {
	if sourceMap == nil {
		if m.Debug {
			log.Println("sourceMap is nil")
		}
		return
	}
	m.sourceMapData = sourceMap

	if m.Debug {
		log.Println("递归调用，","目标:",destStructData ,"数据源:", m.sourceMapData)
	}

	/*
	for i,v := range mapData.(map[string]interface{}) {
		fmt.Printf("key=%s,val=%v,type=%T\n",i,v,v)
	}
	*/

	m.structTypeOf  = reflect.TypeOf(destStructData)
	m.structTofElem = m.structTypeOf.Elem()
	m.structValueOf = reflect.ValueOf(destStructData)
	m.structVofElem = m.structValueOf.Elem()

	//循环映射每个结构体字段
	numField := m.structVofElem.NumField() //结构体字段个数
	for i := 0; i < numField; i++ {
		//检测是否能被设置值
		if !m.structVofElem.Field(i).CanSet() {
			continue
		}

		//结构体字段类型
		structFieldType := m.structTofElem.Field(i).Type.Kind()

		//获取map对应的value
		mapVal,ok := m.getMapValue(i)
		if !ok {
			if m.Debug {
				log.Printf("结构体第%d个字段，获取对应map的值失败",i)
			}
			continue
		}

		//map对应值的类型
		mapValueType := reflect.TypeOf(mapVal).Kind()
		if m.Debug {
			log.Println("map(",mapValueType,") -> ","struct(",structFieldType,")")
		}
		//类型相同的直接set
		if structFieldType == mapValueType {
			switch structFieldType {
			case reflect.Slice://如果都是切片
				m.setSlice(i,mapVal)
			case reflect.Map://如果都是map
				m.setMap(i,mapVal)
			default://其他基本类型直接set
				m.structVofElem.Field(i).Set(reflect.ValueOf(mapVal))
			}
		} else {
			//结构体字段类型和map key对应值的类型不一致
			switch structFieldType {
			//结构体值类型为int 一类
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				m.transformInt(i,mapVal,mapValueType)
			//结构体值类型为uint 一类
			case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				m.transformUint(i,mapVal,mapValueType)
			//结构体值类型为 float 一类
			case reflect.Float32, reflect.Float64:
				m.transformFloat(i,mapVal,mapValueType)
			//结构体值类型为 bool
			case reflect.Bool:
				m.transformBool(i,mapVal,mapValueType)
			//结构体值类型为 string
			case reflect.String:
				m.transformString(i,mapVal,mapValueType)
			//结构体值类型为 struct
			case reflect.Struct:
				m.transformStruct(i,mapVal,mapValueType)
			//结构体值类型为 struct
			case reflect.Slice:
				m.transformSlice(i,mapVal,mapValueType)
			//结构体值类型为 Map
			case reflect.Map:
				m.transformMap(i,mapVal,mapValueType)
			//引用类型
			case reflect.Ptr:
				m.transformPtr(i,mapVal,mapValueType)
			default:
			}
		}
	}
	m.Success = true
}

func (m *MapToStruct) transformInt(i int , mapVal interface{} , mapValueType reflect.Kind) {
	switch mapValueType {
	case reflect.Float64:
		m.structVofElem.Field(i).SetInt(int64(mapVal.(float64)))
	case reflect.String:
		i64, err := strconv.ParseInt(mapVal.(string), 10, 64)
		if err == nil {
			m.structVofElem.Field(i).SetInt(i64)
		} else {
			log.Printf("%q转换成%v失败：%s", mapVal, "int族", err.Error())
		}
	default:
	}
}

func (m *MapToStruct) transformUint(i int , mapVal interface{} , mapValueType reflect.Kind) {
	switch mapValueType {
	case reflect.Float64:
		m.structVofElem.Field(i).SetUint(uint64(mapVal.(float64)))
	case reflect.String:
		ui64, err := strconv.ParseUint(mapVal.(string), 10, 64)
		if err == nil {
			m.structVofElem.Field(i).SetUint(ui64)
		} else {
			log.Printf("%q转换成%v失败：%s", mapVal, "uint族", err.Error())
		}
	default:
	}
}

func (m *MapToStruct) transformFloat(i int , mapVal interface{} , mapValueType reflect.Kind) {
	switch mapValueType {
	case reflect.Float64:
		m.structVofElem.Field(i).SetFloat(mapVal.(float64))
	case reflect.String:
		f64, err := strconv.ParseFloat(mapVal.(string), 64)
		if err == nil {
			m.structVofElem.Field(i).SetFloat(f64)
		} else {
			log.Printf("%q转换成%v失败：%s", mapVal, "float族", err.Error())
		}
	default:
	}
}

func (m *MapToStruct) transformBool(i int , mapVal interface{} , mapValueType reflect.Kind) {
	switch mapValueType {
	case reflect.String:
		mapValStr := strings.ToLower(mapVal.(string))
		if mapValStr == "true" || mapValStr == "1" {
			m.structVofElem.Field(i).SetBool(true)
		} else if mapValStr == "false" || mapValStr == "0" {
			m.structVofElem.Field(i).SetBool(false)
		}
	case reflect.Float64:
		mapValInt := int8(mapVal.(float64))
		if mapValInt == 1 {
			m.structVofElem.Field(i).SetBool(true)
		} else if mapValInt == 0 {
			m.structVofElem.Field(i).SetBool(false)
		}
	default:
	}
}

func (m *MapToStruct) transformString(i int , mapVal interface{} , mapValueType reflect.Kind) {
	switch mapValueType {
	case reflect.Float64:
		//mapValStr := strconv.FormatFloat(mapVal.(float64),'g',10,64)
		mapValStr := fmt.Sprintf("%.f", mapVal.(float64))
		m.structVofElem.Field(i).SetString(mapValStr)
	default:
	}
}

func (m *MapToStruct) transformStruct(i int , mapVal interface{} , mapValueType reflect.Kind) {
	if mapValueType == reflect.Map {
		m.cloneMapToStruct().Transform(m.structVofElem.Field(i).Addr().Interface(), mapVal)
	}
}

func (m *MapToStruct) transformSlice(i int , mapVal interface{} , mapValueType reflect.Kind) {
	if mapValueType == reflect.Map {
		//创建结构体map里面元素的结构体对象
		valTmp := reflect.Indirect(m.structVofElem.Field(i))
		structVal := reflect.New(valTmp.Type().Elem())

		//循环目标map处理
		for _,v := range mapVal.(map[string]interface{}) {
			//递归处理
			m.cloneMapToStruct().Transform(structVal.Interface() , v)

			//把节点append进上层结构体
			m.structVofElem.Field(i).Set(reflect.Append(m.structVofElem.Field(i), structVal.Elem()))
		}
	}
}

func (m *MapToStruct) transformMap(i int , mapVal interface{} , mapValueType reflect.Kind) {
	//需要把map 对应的slice放进strut的map结构中
	if mapValueType == reflect.Slice {
		//对应切片结构体的值
		valTmp := reflect.Indirect(m.structVofElem.Field(i))
		//切面map的list集合
		mapValSli := mapVal.([]interface{})
		//需要make上层map
		m.structVofElem.Field(i).Set(reflect.MakeMap(m.structVofElem.Field(i).Type()))

		for k , v := range mapValSli {
			//new一个切片结构体里面的元素
			structVal := reflect.New(valTmp.Type().Elem())

			//把map映射进结构体
			m.cloneMapToStruct().Transform(structVal.Interface(), v)

			//把map塞进目标map
			key := strconv.Itoa(k)
			m.structVofElem.Field(i).SetMapIndex(reflect.ValueOf(key) , structVal.Elem())
		}
	}
}

func (m *MapToStruct) transformPtr(i int , mapVal interface{} , mapValueType reflect.Kind) {
	//真实的类型
	structFieldTypeReal := m.structTofElem.Field(i).Type.Elem().Kind()
	//log.Println("真实类型：",structFieldTypeReal)
	switch structFieldTypeReal {
	case reflect.String:
		{
			switch mapValueType {
			case reflect.String:
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(mapVal))
				m.structVofElem.Field(i).Set(n)
			case reflect.Float64:
				mapValStr := fmt.Sprintf("%.f", mapVal.(float64))
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(mapValStr))
				m.structVofElem.Field(i).Set(n)
			default:
			}
		}
	case reflect.Struct:
		{
			//初始化struct
			m.structVofElem.Field(i).Set(reflect.New(m.structTofElem.Field(i).Type.Elem()))
			m.cloneMapToStruct().Transform(m.structVofElem.Field(i).Interface(), mapVal)
		}
	case reflect.Int:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(int(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				i64, err := strconv.ParseInt(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structTofElem.Field(i).Type.Elem())
					n.Elem().Set(reflect.ValueOf(int(i64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "int族", err.Error())
				}
			default:
			}
		}
	case reflect.Int8:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(int8(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				i64, err := strconv.ParseInt(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structTofElem.Field(i).Type.Elem())
					n.Elem().Set(reflect.ValueOf(int8(i64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "int族", err.Error())
				}
			default:
			}
		}
	case reflect.Int16:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(int16(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				i64, err := strconv.ParseInt(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structTofElem.Field(i).Type.Elem())
					n.Elem().Set(reflect.ValueOf(int16(i64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "int族", err.Error())
				}
			default:
			}
		}
	case reflect.Int32:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(int32(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				i64, err := strconv.ParseInt(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structTofElem.Field(i).Type.Elem())
					n.Elem().Set(reflect.ValueOf(int32(i64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "int族", err.Error())
				}
			default:
			}
		}
	case reflect.Int64:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(int64(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				i64, err := strconv.ParseInt(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structTofElem.Field(i).Type.Elem())
					n.Elem().Set(reflect.ValueOf(i64))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "int族", err.Error())
				}
			default:
			}
		}
	case reflect.Uint:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structVofElem.Field(i).Type().Elem())
				n.Elem().Set(reflect.ValueOf(uint(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				ui64, err := strconv.ParseUint(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structVofElem.Field(i).Type().Elem())
					n.Elem().Set(reflect.ValueOf(uint(ui64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "uint族", err.Error())
				}
			default:
			}
		}
	case reflect.Uint8:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structVofElem.Field(i).Type().Elem())
				n.Elem().Set(reflect.ValueOf(uint8(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				ui64, err := strconv.ParseUint(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structVofElem.Field(i).Type().Elem())
					n.Elem().Set(reflect.ValueOf(uint8(ui64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "uint族", err.Error())
				}
			default:
			}
		}
	case reflect.Uint16:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structVofElem.Field(i).Type().Elem())
				n.Elem().Set(reflect.ValueOf(uint16(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				ui64, err := strconv.ParseUint(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structVofElem.Field(i).Type().Elem())
					n.Elem().Set(reflect.ValueOf(uint16(ui64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "uint族", err.Error())
				}
			default:
			}
		}
	case reflect.Uint32:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structVofElem.Field(i).Type().Elem())
				n.Elem().Set(reflect.ValueOf(uint32(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				ui64, err := strconv.ParseUint(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structVofElem.Field(i).Type().Elem())
					n.Elem().Set(reflect.ValueOf(uint32(ui64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "uint族", err.Error())
				}
			default:
			}
		}
	case reflect.Uint64:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structVofElem.Field(i).Type().Elem())
				n.Elem().Set(reflect.ValueOf(uint64(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				ui64, err := strconv.ParseUint(mapVal.(string), 10, 64)
				if err == nil {
					n := reflect.New(m.structVofElem.Field(i).Type().Elem())
					n.Elem().Set(reflect.ValueOf(uint64(ui64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "uint族", err.Error())
				}
			default:
			}
		}
	case reflect.Float32:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(float32(mapVal.(float64))))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				f64, err := strconv.ParseFloat(mapVal.(string), 64)
				if err == nil {
					n := reflect.New(m.structTofElem.Field(i).Type.Elem())
					n.Elem().Set(reflect.ValueOf(float32(f64)))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "float族", err.Error())
				}
			default:
			}
		}
	case reflect.Float64:
		{
			switch mapValueType {
			case reflect.Float64:
				n := reflect.New(m.structTofElem.Field(i).Type.Elem())
				n.Elem().Set(reflect.ValueOf(mapVal))
				m.structVofElem.Field(i).Set(n)
			case reflect.String:
				f64, err := strconv.ParseFloat(mapVal.(string), 64)
				if err == nil {
					n := reflect.New(m.structTofElem.Field(i).Type.Elem())
					n.Elem().Set(reflect.ValueOf(f64))
					m.structVofElem.Field(i).Set(n)
				} else {
					log.Printf("%q转换成%v失败：%s", mapVal, "float族", err.Error())
				}
			default:
			}
		}
	default:
	}
}

func (m *MapToStruct) setMap(i int , mapVal interface{}) {
	//创建结构体map里面元素的结构体对象
	valTmp := reflect.Indirect(m.structVofElem.Field(i))
	valTmpTpy := valTmp.Type().Elem().Kind()

	var structVal reflect.Value
	if valTmpTpy == reflect.Ptr {
		structVal = reflect.New(valTmp.Type().Elem().Elem())
		//todo 二级指针处理
		//structVal = reflect.New(valTmp.Type().Elem().Elem().Elem())
	} else {
		structVal = reflect.New(valTmp.Type().Elem())
	}

	//需要make上层map
	m.structVofElem.Field(i).Set(reflect.MakeMap(m.structVofElem.Field(i).Type()))
	//上层结构体map key的类型
	structVofElemKeyType := m.structVofElem.Field(i).Type().Key()

	//循环目标map处理
	var mapkey reflect.Value
	for mk, mv := range mapVal.(map[string]interface{}) {
		//递归处理
		m.cloneMapToStruct().Transform(structVal.Interface() , mv)

		//对结构体map key转换处理
		var i64 int64
		var ui64 uint64
		var f64 float64
		var err error
		switch structVofElemKeyType.Kind() {
		case reflect.String:
			{
				mapkey = reflect.ValueOf(mk)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			{
				i64,err= strconv.ParseInt(mk, 10, 64)
				if err == nil {
					switch structVofElemKeyType.Kind() {
					case reflect.Int:
						mapkey = reflect.ValueOf(int(i64))
					case reflect.Int8:
						mapkey = reflect.ValueOf(int8(i64))
					case reflect.Int16:
						mapkey = reflect.ValueOf(int16(i64))
					case reflect.Int32:
						mapkey = reflect.ValueOf(int32(i64))
					case reflect.Int64:
						mapkey = reflect.ValueOf(i64)
					}
				}
			}
		case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			{
				ui64,err = strconv.ParseUint(mk,10,64)
				if err == nil {
					switch structVofElemKeyType.Kind() {
					case reflect.Uint:
						mapkey = reflect.ValueOf(uint(ui64))
					case reflect.Uintptr:
						mapkey = reflect.ValueOf(uintptr(ui64))
					case reflect.Uint8:
						mapkey = reflect.ValueOf(uint8(ui64))
					case reflect.Uint16:
						mapkey = reflect.ValueOf(uint16(ui64))
					case reflect.Uint32:
						mapkey = reflect.ValueOf(uint32(ui64))
					case reflect.Uint64:
						mapkey = reflect.ValueOf(ui64)
					}
				}
			}
		case reflect.Float32, reflect.Float64:
			{
				f64,err = strconv.ParseFloat(mk,64)
				if err == nil {
					switch structVofElemKeyType.Kind() {
					case reflect.Float32:
						mapkey = reflect.ValueOf(float32(f64))
					case reflect.Float64:
						mapkey = reflect.ValueOf(f64)
					}
				}
			}
		}
		//如果key有效则set值
		if mapkey.IsValid() {
			//把map塞进目标map
			if valTmpTpy == reflect.Ptr {
				m.structVofElem.Field(i).SetMapIndex(mapkey , structVal)

				//todo 二级指针处理
				//structVofElem.Field(i).SetMapIndex(mapkey , reflect.New(structVal.Type()))
			} else {
				m.structVofElem.Field(i).SetMapIndex(mapkey , structVal.Elem())
			}
		}

		//转换失败
		if err != nil{
			if m.Debug {
				log.Println("map key(string) transform fail:" , err.Error())
			}
		}
	}
}

func (m *MapToStruct) setSlice(i int , mapVal interface{}) {
	//实现方式一 [开始]
	//对应切片结构体的值
	valTmp := reflect.Indirect(m.structVofElem.Field(i))
	valTmpTpy := valTmp.Type().Elem().Kind()

	//切面map的list集合
	mapValSli := mapVal.([]interface{})
	var structVal reflect.Value
	for _, v := range mapValSli {
		//structVal 实现方式一
		//structVal = reflect.Indirect(reflect.New(valTmp.Type().Elem())).Addr()

		//structVal 实现方式二
		if valTmpTpy == reflect.Ptr {
			structVal = reflect.New(valTmp.Type().Elem().Elem())
		} else {
			structVal = reflect.New(valTmp.Type().Elem())
		}

		//把map映射进结构体
		m.cloneMapToStruct().Transform(structVal.Interface(), v)

		//把节点append进上层结构体
		if valTmpTpy == reflect.Ptr {
			m.structVofElem.Field(i).Set(reflect.Append(m.structVofElem.Field(i), structVal))
		} else {
			m.structVofElem.Field(i).Set(reflect.Append(m.structVofElem.Field(i), structVal.Elem()))
		}
	}
	//实现方式一 [结束]

	////实现方式二[开始]
	//valTmp := reflect.Indirect(m.structVofElem.Field(i))
	//valTmpTpy := valTmp.Type().Elem().Kind()
	////节点map切片
	//mapValSli := mapVal.([]interface{})
	//mapValSliLen := len(mapValSli)
	////在结构体make mapValSliLen 个元素
	//if valTmpTpy == reflect.Ptr {
	//	m.structVofElem.Field(i).Set(reflect.MakeSlice(valTmp.Type(),mapValSliLen,mapValSliLen))
	//	for n:=0 ; n<mapValSliLen ; n++ {
	//		m.structVofElem.Field(i).Index(n).Set(reflect.New(valTmp.Type().Elem().Elem()))
	//	}
	//} else {
	//	m.structVofElem.Field(i).Set(reflect.MakeSlice(valTmp.Type(),mapValSliLen,mapValSliLen))
	//}
	//
	////循环往每个节点添加map对象里面的值
	//if valTmpTpy == reflect.Ptr {
	//	for j,v:=range mapValSli {
	//		m.cloneMapToStruct().Transform(valTmp.Index(j).Elem().Addr().Interface(), v)
	//	}
	//} else {
	//	for j,v:=range mapValSli {
	//		m.cloneMapToStruct().Transform(valTmp.Index(j).Addr().Interface(), v)
	//	}
	//}
	////实现方式二[结束]
}