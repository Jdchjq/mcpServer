package weather

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"reflect"
	"time"

	"github.com/Jdchjq/mcpServer/internal/types"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type NowWeather struct {
	ObsTime   types.StrText `json:"obsTime" description:"天气数据观测时间"`
	Temp      types.StrText `json:"temp" description:"温度，默认单位：摄氏度"`
	FeelsLike types.StrText `json:"feelsLike" description:"体感温度，默认单位：摄氏度"`
	Text      types.StrText `json:"text" description:"天气状况的文字描述，包括阴晴雨雪等天气状态的描述"`
	Wind360   types.StrText `json:"wind360" description:"风向360角度"`
	WindDir   types.StrText `json:"windDir" description:"风向"`
	WindScale types.StrText `json:"windScale" description:"风力等级"`
	WindSpeed types.StrText `json:"windSpeed" description:"风速，公里/小时"`
	Humidity  types.StrText `json:"humidity" description:"相对湿度，百分比数值"`
	Precip    types.StrText `json:"precip" description:"过去1小时降水量，默认单位：毫米"`
	Pressure  types.StrText `json:"pressure" description:"大气压强，默认单位：百帕"`
	Vis       types.StrText `json:"vis" description:"能见度，默认单位：公里"`
	Cloud     types.StrText `json:"cloud" description:"云量，百分比数值。可能为空"`
	Dew       types.StrText `json:"dew" description:"露点温度。可能为空"`
}

func (n *NowWeather) UnmarshalJSON(data []byte) error {
	type Alias NowWeather
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	v := reflect.ValueOf(n).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		feild := t.Field(i)
		if feild.Type == reflect.TypeOf(types.StrText{}) {
			customTag := feild.Tag.Get(n.getCommentTag())

			f := v.Field(i)
			ct := f.Interface().(types.StrText)
			if ct.Description == "" {
				ct.Description = customTag
				f.Set(reflect.ValueOf(ct))
			}

		}

	}
	return nil
}

func (n *NowWeather) getCommentTag() string {
	return "description"
}

type WeatherResponse struct {
	Code       string     `json:"code"`
	UpdateTime string     `json:"updateTime"`
	FxLink     string     `json:"fxLink"`
	Now        NowWeather `json:"now"`
}

func NewWeather(baseurl string, alg string, kid string, sub string, key []byte) *WeatherCaller {
	signer := NewSigner(alg, kid, sub, []byte(key))
	cli := resty.New().SetBaseURL(baseurl)
	return &WeatherCaller{
		Signer: signer,
		r:      cli,
	}
}

// 获取实时天气
func (w *WeatherCaller) GetWeather(location string, lang string) (result string, err error) {
	token := w.Signer.Sign()

	url := "/v7/weather/now"
	res, err := w.r.R().SetAuthToken(token).SetQueryParams(map[string]string{"location": location, "lang": lang}).Get(url)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	if res.StatusCode() != 200 {
		log.Err(fmt.Errorf("get statusCode:%d, body:%s", res.StatusCode(), string(res.Body()))).Msg("")
		return "", err
	}

	return w.DecodeWeather(res.Body())
}

func (w *WeatherCaller) DecodeWeather(rawBody []byte) (result string, err error) {
	var res WeatherResponse
	err = json.Unmarshal(rawBody, &res)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	resultData, err := json.Marshal(res)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	return string(resultData), nil
}

type WSigner struct {
	sub        string // 签发主体，在控制台获取
	alg        string // 签名算法
	kid        string // 凭据ID 在控制台获取
	privateKey []byte
}

type signHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type signPayload struct {
	Sub string `json:"sub"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

func NewSigner(Alg string, Kid string, Sub string, privatekey []byte) WSigner {
	block, _ := pem.Decode(privatekey)
	if block == nil || block.Type != "PRIVATE KEY" {
		panic("无效的PEM文件或类型不正确")
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Panic().Err(err).Msg("decode private key")
	}

	edPrivateKey, ok := private.(ed25519.PrivateKey)
	if !ok {
		log.Panic().Err(err).Msg("not ed25519 key")
	}

	return WSigner{
		alg:        Alg,
		sub:        Sub,
		kid:        Kid,
		privateKey: edPrivateKey,
	}

}

func (w WSigner) Sign() string {
	header := signHeader{
		Alg: w.alg,
		Kid: w.kid,
	}

	headerContent, _ := json.Marshal(header)
	encodeHeader := make([]byte, base64.URLEncoding.EncodedLen(len(headerContent)))
	base64.URLEncoding.Encode(encodeHeader, headerContent)

	payload := signPayload{
		Sub: w.sub,
		Iat: time.Now().Add(-30 * time.Second).Unix(), // 	签发时间建议当前时间30s前
		Exp: time.Now().Add(20 * time.Minute).Unix(),
	}
	payloadContent, _ := json.Marshal(payload)
	encodePayload := make([]byte, base64.URLEncoding.EncodedLen(len(payloadContent)))
	base64.URLEncoding.Encode(encodePayload, payloadContent)

	signatureContent := fmt.Sprintf("%s.%s", encodeHeader, encodePayload)
	signedContent := ed25519.Sign(w.privateKey, []byte(signatureContent))

	signed := make([]byte, base64.URLEncoding.EncodedLen(len(signedContent)))
	base64.URLEncoding.Encode(signed, signedContent)

	token := fmt.Sprintf("%s.%s.%s", encodeHeader, encodePayload, signed)
	return token
}
