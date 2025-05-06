package weather

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/mark3labs/mcp-go/mcp"
)

type WeatherReq struct {
	Location string `json:"location" description:"经纬度坐标，格式是经度,纬度，最多只保留小数点后两位。注意：中国大陆地区应使用GCJ-02坐标系，在其他地区应使用WGS-84坐标系" required:"true"`
}

var TheWeatherCaller *WeatherCaller

type WeatherCaller struct {
	Signer WSigner
	r      *resty.Client
}

func HandleWeatherRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	location := request.Params.Arguments["location"].(string)

	result, err := TheWeatherCaller.GetWeather(location, "zh")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}
