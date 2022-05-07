package baidu

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// FReqPrefix Baidu翻译Request前缀
const FReqPrefix = "kw="

// TranslationUrl Baidu翻译RequestURL
const TranslationUrl = "https://fanyi.baidu.com/sug"

// DictResponse Baidu翻译返回结果
type DictResponse struct {
	Errno int `json:"errno"`
	Data  []struct {
		K string `json:"k"`
		V string `json:"v"`
	} `json:"data"`
}

func Query(word string, rep chan string) {
	client := &http.Client{}
	// 封装请求参数
	requestBuilder := strings.Builder{}
	requestBuilder.WriteString(FReqPrefix)
	requestBuilder.WriteString("hello")
	var data = strings.NewReader(requestBuilder.String())
	req, err := http.NewRequest("POST", TranslationUrl, data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", `BAIDUID_BFESS=41BC552B3AB37D0D8D06C86AE79D7BAA:FG=1; RT="z=1&dm=baidu.com&si=9i8im6r7y6s&ss=l2vkyse0&sl=3&tt=39k&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=44r&ul=j8n2&hd=j8o1"; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1651915928; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1651915943; ab_sr=1.0.1_OTBlNzIyZjViOTBmZGYwNmFmNzgxMDExMzQ3ZDY3MDQ5NzRhZTNhM2ExNTY1NjhmM2I2MDUzYzBkOTg0N2Y0NDA0Njk4ODNmYTcxMGIxY2NlNjhjNmU0NmUwYWEzZDQ5MzZjZmE3NTRmYjM5ODA3ZjY0YmE4NzEzYjJhZmM5NTMyY2Q5YzRlMDk0YjIxODY3N2FjMmRiZjRiMGNmNmZjMw==`)
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Referer", "https://fanyi.baidu.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 处理返回结果
	var baiduDictResponse DictResponse
	err = json.Unmarshal(bodyText, &baiduDictResponse)
	if err != nil {
		log.Fatal(err)
	}

	builder := strings.Builder{}
	for _, item := range baiduDictResponse.Data {
		builder.WriteString(item.V + "\n")
	}
	rep <- builder.String()
}
