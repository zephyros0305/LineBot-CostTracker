package main

import (
	"LineBot-CostTracker/vendor/github.com/line/line-bot-sdk-go/v7/linebot"
	"fmt"
)

func GetListRecordResponse(records []Record) *linebot.BubbleContainer {
	var contents []linebot.FlexComponent

	for _, r := range records {
		recordContent := &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#000000",
					Flex:   linebot.IntPtr(1),
					Text:   r.CreatedAt.Format("2006-01-02"),
				},
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Size:  linebot.FlexTextSizeTypeSm,
					Color: "#111111",
					Flex:  linebot.IntPtr(2),
					Text:  r.Memo,
				},
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Size:  linebot.FlexTextSizeTypeSm,
					Align: linebot.FlexComponentAlignTypeEnd,
					Color: "#000000",
					Flex:  linebot.IntPtr(0),
					Text:  fmt.Sprintf("$%d", r.Cost),
				},
			},
		}
		contents = append(contents, recordContent)
	}

	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeGiga,
		Body: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Contents: contents,
		},
	}

	return container
}
