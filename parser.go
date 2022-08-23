package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type StructInfos struct {
	StructInfos []StructInfo
}

type StructInfo struct {
	FieldInfos []FieldInfo
	StructName string
	FileName   string
}

type FieldInfo struct {
	FieldName string
	FieldType string
	DbName    string
}

type NestedStruct struct {
	StartLine int64
	EndLine   int64
	Content   string
	FieldName string
	IsArray   bool
}

func (n NestedStruct) Parse() (FieldInfo, []StructInfo) {
	fieldInfo := FieldInfo{
		FieldName: UpperFirst(n.FieldName),
		FieldType: UpperFirst(n.FieldName),
		DbName:    n.FieldName,
	}
	if n.IsArray {
		if strings.HasSuffix(fieldInfo.FieldType, "ies") {
			fieldInfo.FieldType = strings.TrimRight(fieldInfo.FieldType, "ies")
			fieldInfo.FieldType = fmt.Sprintf("%sy", fieldInfo.FieldType)
		}
		if strings.HasSuffix(fieldInfo.FieldType, "es") {
			fieldInfo.FieldType = strings.TrimRight(fieldInfo.FieldType, "es")
		}
		if strings.HasSuffix(fieldInfo.FieldType, "s") {
			fieldInfo.FieldType = strings.TrimRight(fieldInfo.FieldType, "s")
		}
		fieldInfo.FieldType = fmt.Sprintf("[]%s", fieldInfo.FieldType)
	}
	return fieldInfo, Parse(n.Content)
}

var typeMap map[string]string

func init() {
	b, err := ioutil.ReadFile("type_mapping.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &typeMap)
	if err != nil {
		panic(err)
	}
}

func Parse(content string) []StructInfo {
	infos := []StructInfo{}
	lines := strings.Split(content, "\n")
	if len(lines) == 1 {
		lines = getLines(content)
	}
	currentStruct := StructInfo{}
	nestedStructs := []NestedStruct{}
	currentNestedStruct := NestedStruct{}
	for i, line := range lines {
		if strings.HasPrefix(line, "#") {
			structName, fileName := getStructNameAndFileName(line)
			currentStruct.StructName = structName
			currentStruct.FileName = fileName
			continue
		}
		if strings.HasPrefix(line, "```") || line == "{" || line == "}" {
			continue
		}
		if target := trimAllSpaces(line); target != "" {
			if strings.HasPrefix(target, "[") || strings.HasSuffix(target, "]") {
				continue
			}
			if strings.HasSuffix(target, "[") {
				currentNestedStruct.IsArray = true
				currentNestedStruct.FieldName = strings.Split(target, ": ")[0]
			}
			if strings.HasSuffix(target, "{") {
				currentNestedStruct.StartLine = int64(i)
				if currentNestedStruct.FieldName == "" {
					currentNestedStruct.FieldName = strings.Split(target, ": ")[0]
				}
				continue
			}
			if strings.HasPrefix(target, "}") {
				currentNestedStruct.EndLine = int64(i)
				currentNestedStruct.Content = trimAllSpaces(strings.Join(lines[currentNestedStruct.StartLine:currentNestedStruct.EndLine+1], ""))
				nestedStructs = append(nestedStructs, currentNestedStruct)
				currentNestedStruct = NestedStruct{}
				continue
			}

			if !strings.HasSuffix(target, "[") && !strings.HasSuffix(target, "{") {
				if currentNestedStruct.StartLine != 0 {
					continue
				}
				currentStruct.FieldInfos = append(currentStruct.FieldInfos, genFieldInfo(target))
			}
		}
	}
	for _, n := range nestedStructs {
		fieldInfo, structInfos := n.Parse()
		currentStruct.FieldInfos = append(currentStruct.FieldInfos, fieldInfo)
		infos = append(infos, structInfos...)
	}
	infos = append(infos, currentStruct)
	return infos
}

func genFieldInfo(line string) FieldInfo {
	s := strings.Split(strings.Split(line, ",")[0], ": ")
	result := FieldInfo{
		DbName: s[0],
	}
	if s[0] == "_id" {
		result.FieldName = "Id"
	} else {
		result.FieldName = UpperFirst(s[0])
	}
	result.FieldType = typeMap[s[1]]
	return result
}

func getStructNameAndFileName(line string) (string, string) {
	s := strings.Split(line, " ")
	s = strings.Split(s[1], ".")
	return UpperFirst(s[len(s)-1]), ToSnakeCase(s[len(s)-1])
}

func trimAllSpaces(content string) string {
	temp := strings.TrimSpace(content)
	if temp == content {
		return temp
	}
	return trimAllSpaces(temp)
}

func beforeGetLines(content string) string {
	for !strings.HasPrefix(content, "{") {
		content = strings.TrimPrefix(content, content[0:1])
	}
	return content
}

func getLines(content string) []string {
	content = beforeGetLines(content)
	lines := []string{"{"}
	beforeProcess := []string{}
	temp := strings.Split(content, ",")
	for _, v := range temp {
		beforeProcess = append(beforeProcess, strings.Split(v, " ")...)
	}
	for i, v := range beforeProcess {
		if strings.HasSuffix(v, ":") {
			lines = append(lines, fmt.Sprintf("%s %s", v, beforeProcess[i+1]))
		}
	}
	lines = append(lines, "}")
	return lines
}
