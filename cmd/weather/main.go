package main

import (
	"flag"
	"fmt"

	"github.com/Jdchjq/mcpServer/cmd/weather/config"
	"github.com/Jdchjq/mcpServer/internal/tools/weather"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		configDir string
		transport string
	)

	flag.StringVar(&configDir, "configDir", "", "配置文件所在路径")
	flag.StringVar(&configDir, "c", "", "配置文件所在路径")
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.Parse()
	if configDir == "" {
		log.Panic().Msg("config Dir empty")
	}

	config.SetConfig(configDir)

	s := server.NewMCPServer(
		"Get weather data",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)
	InitTools()
	RegisterTools(s)

	if transport == "sse" {
		sseServer := server.NewSSEServer(s, server.WithBaseURL("http://localhost:8080"))
		log.Printf("SSE server listening on :8080")
		if err := sseServer.Start(":8080"); err != nil {
			log.Panic().Err(err).Msg("Server error")
		}
	} else if transport == "stdio" {
		log.Info().Msg("mcp server start on stdio")
		if err := server.ServeStdio(s); err != nil {
			fmt.Printf("Server error: %v\n", err)
			return
		}
	}

	// guangzhouLocation := "113.37,23.12"
	// result, err := weather.TheWeatherCaller.GetWeather(guangzhouLocation, "zh")
	// if err != nil {
	// 	log.Err(err).Msg("")
	// } else {
	// 	log.Info().Str("result", result).Msg("")
	// }
}

func InitTools() {
	weather.TheWeatherCaller = weather.NewWeather(
		config.Config.Weather.BaseUrl,
		config.Config.Weather.Alg,
		config.Config.Weather.Kid,
		config.Config.Weather.Sub,
		config.Config.Weather.Key,
	)
}

func RegisterTools(s *server.MCPServer) {
	nowWeatherTool := mcp.NewTool(
		"weather_now",
		mcp.WithDescription("get the now-time weather on input location"),
		mcp.WithString("location", mcp.Required(), mcp.Description("经纬度坐标，格式是经度,纬度，最多只保留小数点后两位。注意：中国大陆地区应使用GCJ-02坐标系，在其他地区应使用WGS-84坐标系")),
	)

	s.AddTool(nowWeatherTool, weather.HandleWeatherRequest)
}
