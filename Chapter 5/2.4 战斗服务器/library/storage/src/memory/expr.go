package memory

import (
	"ProjectX/library/storage/src/tools"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
)

type elemType uint8
type operatorType uint8

const (
	// Value 数值对象，通过 ? 代替
	Value elemType = iota
	// Field 结构体字段对象
	Field
	// Operator 操作符
	Operator
)

const (
	Greater operatorType = iota
	EqualOrGreater
	Less
	EqualOrLess
	Equal
	NotEqual
	And
	Or
	Not
)

type nilTest struct {
}

func (t *nilTest) Logf(_ string, _ ...interface{}) {
}

func (t *nilTest) Errorf(_ string, _ ...interface{}) {
}

func (t *nilTest) FailNow() {
}

var (
	priorityMap = map[operatorType]uint8{
		Greater:        11,
		EqualOrGreater: 11,
		Less:           11,
		EqualOrLess:    11,
		Equal:          11,
		NotEqual:       11,
		And:            2,
		Or:             1,
		Not:            3,
	}

	markMap = map[string]operatorType{
		">":   Greater,
		">=":  EqualOrGreater,
		"<":   Less,
		"<=":  EqualOrLess,
		"=":   Equal,
		"<>":  NotEqual,
		"AND": And,
		"OR":  Or,
		"NOT": Not,
	}

	exprMap = map[string]*exprNode{}

	t = &nilTest{}
)

type exprNode struct {
	typ    elemType
	father *exprNode
	lChild *exprNode
	rChild *exprNode
	value  interface{}
}

// prior before 是否比 after 优先级低， true下，false上
func prior(before, after operatorType) bool {
	if priorityMap[after] > priorityMap[before] {
		return true
	}
	return false
}

// getExprRoot 获取表达式的根节点
// TODO: 完善逻辑，在表达式复杂时，下面可能会导致bug
func getExprRoot(expr string) *exprNode {
	if nod, ok := exprMap[expr]; ok {
		return nod
	} else {
		var root, current *exprNode
		st := newStack()
		fieldContent := ""
		valueCount := 0
		exprLength := len(expr)
		for i := 0; i < exprLength; i++ {
			if expr[i] == ' ' {
				if fieldContent != "" {
					newNode := setNode(Field, fieldContent, current, true)
					if current == nil {
						current = newNode
					}
					fieldContent = ""
				}
				continue
			}

			switch expr[i] {
			case '?':
				if fieldContent != "" {
					return nil
				}
				newNode := setNode(Value, valueCount, current, true)
				if current == nil {
					current = newNode
				}
				valueCount++
				break
			case '(':
				if fieldContent != "" {
					return nil
				}
				st.Push(current)
				current = nil
				break
			case ')':
				if fieldContent != "" {
					newNode := setNode(Field, fieldContent, current, true)
					if current == nil {
						current = newNode
					}
					fieldContent = ""
				}

				beforeCurrent := st.Pop()
				if beforeCurrent != nil {
					setLevel(beforeCurrent, current)
					current = beforeCurrent
				}
				break
			default:
				fieldContent += string(expr[i])
				if opt, ok := markMap[fieldContent]; ok {
					if (fieldContent == ">" || fieldContent == "<") && i < exprLength-1 {
						_, nextOK := markMap[fieldContent+string(expr[i+1])]
						if nextOK {
							continue
						}
					}
					current = setNode(Operator, opt, current, false)
					fieldContent = ""
					if current.father == nil && (root == nil || root.father != nil) {
						root = current
					}
				}

				break
			}
		}

		if fieldContent != "" {
			setNode(Field, fieldContent, current, true)
		}
		if st.Size() != 0 {
			return nil
		}
		exprMap[expr] = root

		return root
	}
}

func setLevel(father *exprNode, child *exprNode) {
	if child != nil && child.father != nil {
		setLevel(father, child.father)
		return
	}

	if father.lChild == nil {
		father.lChild = child
	} else {
		father.rChild = child
	}
	if child != nil {
		child.father = father
	}
}

func setNode(typ elemType, content interface{}, current *exprNode, isFather bool) *exprNode {
	newNode := &exprNode{
		typ:   typ,
		value: content,
	}

	if current == nil {
		return newNode
	}

	if isFather {
		newNode.father = current
		if current.lChild == nil {
			current.lChild = newNode
		} else if current.rChild != nil {
			newNode.lChild = current.rChild
			current.rChild = newNode
		} else {
			current.rChild = newNode
		}
		return newNode
	}

	if typ == Operator {
		if current.typ != Operator {
			newNode.lChild = current
			current.father = newNode
			return newNode
		} else if prior(current.value.(operatorType), content.(operatorType)) {
			return setNode(typ, content, current, true)
		} else {
			if current.father == nil {
				newNode.lChild = current
				current.father = newNode
				return newNode
			} else {
				return setNode(typ, content, current.father, false)
			}
		}
	}

	return newNode
}

// calculate 计算结果
func (bt *exprNode) calculate(value interface{}, args []interface{}) bool {
	if bt.typ != Operator {
		return getBool(bt.getValue(value, args))
	}

	opt := bt.value.(operatorType)
	switch opt {
	case And:
		if bt.lChild == nil {
			return false
		} else if bt.rChild == nil {
			return bt.lChild.calculate(value, args)
		} else {
			return bt.lChild.calculate(value, args) && bt.rChild.calculate(value, args)
		}
	case Or:
		if bt.lChild == nil {
			return false
		} else if bt.rChild == nil {
			return bt.lChild.calculate(value, args)
		} else {
			return bt.lChild.calculate(value, args) || bt.rChild.calculate(value, args)
		}
	case Not:
		if bt.lChild == nil {
			return false
		} else {
			return !bt.lChild.calculate(value, args)
		}
	case Greater:
		if bt.lChild == nil {
			return false
		} else if bt.rChild == nil {
			return false
		} else {
			lValue := bt.lChild.getValue(value, args)
			rValue := bt.rChild.getValue(value, args)
			if lValue == nil || rValue == nil {
				return false
			}
			return assert.Greater(t, lValue, rValue)
		}
	case EqualOrGreater:
		if bt.lChild == nil {
			return false
		} else if bt.rChild == nil {
			return false
		} else {
			lValue := bt.lChild.getValue(value, args)
			rValue := bt.rChild.getValue(value, args)
			if lValue == nil || rValue == nil {
				return false
			}
			return assert.GreaterOrEqual(t, lValue, rValue)
		}
	case Less:
		if bt.lChild == nil {
			return false
		} else if bt.rChild == nil {
			return false
		} else {
			lValue := bt.lChild.getValue(value, args)
			rValue := bt.rChild.getValue(value, args)
			if lValue == nil || rValue == nil {
				return false
			}
			return assert.Less(t, lValue, rValue)
		}
	case EqualOrLess:
		if bt.lChild == nil {
			return false
		} else if bt.rChild == nil {
			return false
		} else {
			lValue := bt.lChild.getValue(value, args)
			rValue := bt.rChild.getValue(value, args)
			if lValue == nil || rValue == nil {
				return false
			}
			return assert.LessOrEqual(t, lValue, rValue)
		}
	case Equal:
		if bt.lChild == nil {
			return false
		} else if bt.rChild == nil {
			return false
		} else {
			lValue := bt.lChild.getValue(value, args)
			rValue := bt.rChild.getValue(value, args)
			if lValue == nil || rValue == nil {
				return false
			}
			return assert.Equal(t, lValue, rValue)
		}
	case NotEqual:
		if bt.lChild == nil {
			return false
		} else if bt.rChild == nil {
			return false
		} else {
			lValue := bt.lChild.getValue(value, args)
			rValue := bt.rChild.getValue(value, args)
			if lValue == nil || rValue == nil {
				return false
			}
			return assert.NotEqual(t, lValue, rValue)
		}
	}

	return false
}

func (bt *exprNode) getValue(value interface{}, args []interface{}) interface{} {
	if bt.typ == Value {
		index, ok := bt.value.(int)
		if !ok {
			return nil
		}
		if len(args) > index {
			return args[index]
		}
		return nil
	}
	if bt.typ == Field {
		name, ok := bt.value.(string)
		if !ok {
			return nil
		}
		if strings.HasPrefix(name, "'") {
			name = name[1 : len(name)-1]
		}
		return tools.GetFieldValueByRealName(value, name)
	}

	return bt.value
}

func getBool(value interface{}) bool {
	if flag, ok := value.(bool); ok {
		return flag
	}

	str := fmt.Sprintf("%v", value)
	if str != "" && str != "0" {
		return true
	} else {
		return false
	}
}
