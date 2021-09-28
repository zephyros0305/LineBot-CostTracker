package main

import (
	"regexp"
	"strconv"
	"strings"
)

const getRegStr = `^(近|最近)`

func DetermineOperation(text string) (OperationType, *OperationData) {
	var operationInfo OperationData
	fields := strings.Fields(text)

	switch fieldLen := len(fields); fieldLen {
	case 0:
		return Error, nil
	case 1:
		if cost, err := strconv.ParseUint(fields[0], 10, 64); err == nil {
			operationInfo.Number = cost
			return KeepRecord, &operationInfo
		} else {
			operaction := fields[0]
			getReg := regexp.MustCompile(getRegStr)
			switch {
			case getReg.FindStringIndex(operaction) != nil:
				numStr := getReg.ReplaceAllString(operaction, "")
				if num, err := strconv.ParseUint(numStr, 10, 64); err == nil {
					operationInfo.Number = num
					return GetRecord, &operationInfo
				} else {
					return Error, nil
				}
			default:
				return Error, nil
			}
		}
	case 2:
		if cost, err := strconv.ParseUint(fields[0], 10, 64); err == nil {
			operationInfo.Number = cost
			operationInfo.Memo = fields[1]
			return KeepRecord, &operationInfo
		}

		if cost, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			operationInfo.CostType = fields[0]
			operationInfo.Number = cost
			return KeepRecord, &operationInfo
		}
	default:
		if cost, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			operationInfo.CostType = fields[0]
			operationInfo.Number = cost
			operationInfo.Memo = fields[2]
			return KeepRecord, &operationInfo
		}
	}

	return Error, nil
}
