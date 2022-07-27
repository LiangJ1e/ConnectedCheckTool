package checkTool

import (
	"errors"  //用于错误信息的传递
	"reflect" //关键方法
	"strconv" //实现字符串和数值的转换
	"strings" //字符串操作
)

var structName []string        //储存结构体中的字段名称
var structValueInt []int       //储存结构体中int类型的数值
var structValueString []string //储存结构体中string类型的数值
var structType []string        //储存结构体中字段的数据类型
var structCheckMsg []string    //储存check中的信息
var length int                 //储存结构体字段的数量（总长度）
var intType int                //用于判断int类型
var stringType string          //用于判断string类型

type GetCheckTool interface {
	CheckStruct(c interface{}) (ok bool, err error, msg string)
}
type Check struct {
}

func initArray(c interface{}) {
	structName = nil
	structValueInt = nil
	structValueString = nil
	structType = nil
	structCheckMsg = nil
	item := reflect.ValueOf(c).Elem() //获取结构体字段的数值
	field := reflect.TypeOf(c).Elem() //获取结构体字段类型信息
	length = item.NumField()
	for i := 0; i < length; i++ {
		fieldItem := item.Field(i)
		fieldType := field.Field(i)
		structName = append(structName, fieldType.Name) //储存字段名称
		if fieldItem.Type() == reflect.TypeOf(intType) {
			structValueInt = append(structValueInt, int(fieldItem.Int()))
			structValueString = append(structValueString, "")
			structType = append(structType, "int")
		} //储存字段数值
		if fieldItem.Type() == reflect.TypeOf(stringType) {
			structValueString = append(structValueString, fieldItem.String())
			structValueInt = append(structValueInt, 0)
			structType = append(structType, "string")
		}
		msg := ""
		tag := fieldType.Tag
		label := tag.Get("check")
		msg = label
		structCheckMsg = append(structCheckMsg, msg) //储存check信息
	}
}

/*func getTarget(c interface{}, index int) (msg string, exist bool, err error) {
	msg = ""
	field := reflect.ValueOf(c).Type().Field(index)
	tag := field.Tag
	label := tag.Get("check")
	msg = label
	if msg != "" {
		return msg, true, err
	}
	return msg, false, err
}*/

func getLimitMsg(limitMsg string) (limits []string, err error) {
	result := strings.Split(limitMsg, "&&")
	for i := 0; i < len(result); i++ {
		result[i] = strings.Split(result[i], ")")[0]
	}
	return result, err
}

func getLimitConditions(limitMsg string) (limitConditions []string, err error) {
	level1 := strings.Split(limitMsg, "(")
	head := strings.Index(limitMsg, "(")
	if head == -1 {
		return nil, errors.New("格式出错！请检查是否有when")
	}
	limitConditions = strings.Split(level1[1], ",")
	return limitConditions, err

}

func getLimitRequest(index int, limitMsg string) (limitRequest []string, err error) {
	level1 := strings.Split(limitMsg, "(")
	request := strings.Split(level1[0], ",")
	for i := 0; i < len(request); i++ {
		limitRequest = append(limitRequest, structName[index]+":"+request[i])
	}
	return limitRequest, err
}

func getCheckName(deJudgeString string) (index int, err error) {
	divideString := strings.Split(deJudgeString, ":")
	name := divideString[0]
	for i := 0; i < length; i++ {
		data := structName[i]
		if data == name {
			index = i
			return index, err
		}
	}
	return -1, errors.New("无法找到匹配字段，请检查是否书写错误")
}

func getCheckType(deJudgeString string) (condition int, value string, err error) {
	divideString := strings.Split(deJudgeString, ":")
	message := divideString[1]
	levelDivide := strings.Split(message, "=")

	conditionString := levelDivide[0]
	for i := 0; i < len(conditionType); i++ {
		if conditionString == conditionType[i] {
			if len(levelDivide) == 2 {
				return i, levelDivide[1], err
			}
			return i, "", err
		}
	}
	return -1, "", errors.New("未声明的条件" + conditionString)
}

var conditionType = [...]string{
	"max",
	"min",
	"val",
	"startwith",
	"required",
	"gte",
	"lte",
}

func judgeOK(index int, value string, judegType int) (err error) {
	targetType := structType[index]
	targetName := structName[index]
	if targetType == "int" {
		targetValue := structValueInt[index]
		switch judegType {
		case 0:
			return errors.New(targetName + "字段max不能用于int类型")
		case 1:
			return errors.New(targetName + "min不能用于int类型")
		case 2:
			val, errs := strconv.Atoi(value)
			if errs != nil {
				return errors.New(targetName + "数值无法和字符串比较！")
			}
			if val == targetValue {
				return err
			}
			return errors.New(targetName + "字段数值不符合要求")
		case 3:
			return errors.New(targetName + "字段startwith不能用于Int类型")
		case 4:
			if targetValue == 0 {
				return errors.New(targetName + "缺少必须数值！")
			}
			return err
		case 5:
			val, errs := strconv.Atoi(value)
			if errs != nil {
				return errors.New(targetName + "数值无法和字符串比较！")
			}
			if targetValue > val {
				return err
			}
			return errors.New(targetName + "数值小于限制")
		case 6:
			val, errs := strconv.Atoi(value)
			if errs != nil {
				return errors.New(targetName + "数值无法和字符串比较！")
			}
			if targetValue < val {
				return err
			}
			return errors.New(targetName + "数值大于限制")
		}
	}
	if targetType == "string" {
		targetString := structValueString[index]
		switch judegType {
		case 0:
			val, errs := strconv.Atoi(value)
			if errs != nil {
				return errors.New(targetName + "数值无法和字符串比较！")
			}
			if len(targetString) < val {
				return err
			}
			return errors.New(targetName + "字符串长度大于限制")
		case 1:
			val, errs := strconv.Atoi(value)
			if errs != nil {
				return errors.New(targetName + "数值无法和字符串比较！")
			}
			if len(targetString) < val {
				return err
			}
			return errors.New(targetName + "字符串长度小于限制")
		case 2:
			if value == targetString {
				return err
			}
			return errors.New(targetName + "字符串不满足确定文字")
		case 3:
			head := strings.Index(targetString, value)
			if head == 0 {
				return err
			}
			return errors.New(targetName + "字符串不以限制字段" + value + "开头")
		case 4:
			if targetString == "" {
				return errors.New(targetName + "缺少必须字符串！")
			}
			return err
		case 5:
			return errors.New(targetName + "标签gte不能用于string类型")
		case 6:
			return errors.New(targetName + "标签lte不能用于string类型")
		}
	}
	return err
}

func (C *Check) CheckStruct(c interface{}) (err error) {
	initArray(c)
	for i := 0; i < length; i++ {
		if structCheckMsg[i] != "" {
			var limitList []string
			limitList, err = getLimitMsg(structCheckMsg[i])
			for j := 0; j < len(limitList); j++ {
				var limitConditions []string
				limitConditions, err = getLimitConditions(limitList[j])
				if err != nil {
					return err
				}
				for a := 0; a < len(limitConditions); a++ {
					var checkName int
					var checkType int
					var value string

					checkName, err = getCheckName(limitConditions[a])
					if err != nil {
						return err
					}
					checkType, value, err = getCheckType(limitConditions[a])
					if err != nil {
						return err
					}
					err = judgeOK(checkName, value, checkType)

					if err != nil {
						return nil
					}
				}
				var limitRequests []string
				limitRequests, err = getLimitRequest(i, limitList[j])
				if err != nil {
					return err
				}
				for a := 0; a < len(limitRequests); a++ {
					var checkName int
					var checkType int
					var value string
					checkName, err = getCheckName(limitRequests[a])
					if err != nil {
						return err
					}
					checkType, value, err = getCheckType(limitRequests[a])
					if err != nil {
						return err
					}
					err = judgeOK(checkName, value, checkType)
					if err != nil {
						return err
					}

				}
			}

		}

	}
	return err

}
