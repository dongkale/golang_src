package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/tables"
	"bizt.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// 설정 화일 읽기
// https://forteleaf.tistory.com/entry/json-%EC%84%A4%EC%A0%95%ED%8C%8C%EC%9D%BC-%EC%9D%BD%EC%96%B4%EC%98%A4%EA%B8%B0

// // Config ...
// type Config struct {
// 	LiveNvnApplyMaxCount int `json:"LiveNvnApplyMaxCount"`

// 	Test int `json:"TEST"`
// }

// // ConfigFile ...
// var ConfigFile string = "config.json"

// ReloadConfigController ...
type ReloadConfigController struct {
	beego.Controller
}

// Post ...
func (c *ReloadConfigController) Post() {

	session := c.StartSession()

	memNo := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNo == nil {
		//c.Ctx.Redirect(302, "/login")
		c.Data["json"] = &models.DefaultResult{
			RtnCd:  99,
			RtnMsg: "FAIL",
		}

		c.ServeJSON()
		return
	}

	fmt.Printf(fmt.Sprintf("[ReloadConfigController] Start ----------------------------->"))

	//fmt.Printf(fmt.Sprintf("[ReloadTemplateController] reload File: %v", beego.BConfig.WebConfig.ViewsPath))
	//fmt.Printf(fmt.Sprintf("[ReloadTemplateController] reload File: %v", ConfigFile))

	// var rtnCd int64 = 0
	// var rtnMsg string = ""

	// config, err := LoadConfig()
	// if err == nil {
	// 	// fmt.Printf(fmt.Sprintf("[Error] %v", err))
	// 	// fmt.Printf(fmt.Sprintf("[Config] %v", config))
	// 	// fmt.Println(err)
	// 	// fmt.Println(config)

	// 	ConfigSetting(config)
	// 	beego.Info(fmt.Sprintf("[Config] Load Ok"))

	// 	// rtnCd = 1
	// 	// rtnMsg = "SUCCESS"
	// } else {
	// 	beego.Info(fmt.Sprintf("[Config] Error: %v", err))
	// 	beego.Info(fmt.Sprintf("[Config] Load Fail"))
	// 	// fmt.Println(err)

	// 	// rtnCd = 0
	// 	// rtnMsg = "ERROR"
	// }

	var loadTableConf tables.TableConfig

	loadTableConf = tables.TableConf

	confFile, err := beego.AppConfig.String("setConfFile")
	if err != nil {
		logs.Error(fmt.Sprintf("[LoadConfig] Error getting config file: %v", err))
		c.Data["json"] = &models.DefaultResult{
			RtnCd:  99,
			RtnMsg: "FAIL",
		}
		c.ServeJSON()
		return
	}
	
	errCnt := utils.LoadConfigFile(confFile, &loadTableConf)
	if errCnt == 0 {
		result := loadTableConf.IsCheckValue()
		if result != "" {
			logs.Info(fmt.Sprintf("[LoadConfig] Load Fail!! -> %s InValid Value", result))
		} else {
			tables.TableConf = loadTableConf
			logs.Info(fmt.Sprintf("[LoadConfig] Load Ok!!"))
		}
	} else {
		logs.Info(fmt.Sprintf("[LoadConfig] Load Fail!!"))
	}

	rtnData := models.DefaultResult{
		RtnCd:  1,
		RtnMsg: "SUCCESS",
	}

	c.Data["json"] = &rtnData

	c.ServeJSON()

	fmt.Printf(fmt.Sprintf("[ReloadConfigController] End ----------------------------->"))
}

// // LoadConfig ...
// func LoadConfig() (Config, error) {
// 	var config Config

// 	file, err := os.Open(ConfigFile)
// 	defer file.Close()
// 	if err != nil {
// 		return config, fmt.Errorf("Error %v", err)
// 	}

// 	decoder := json.NewDecoder(file)
// 	//decoder.DisallowUnknownFields() // Force errors

// 	err = decoder.Decode(&config)
// 	if err != nil {
// 		return config, fmt.Errorf("Error %v", err)
// 	}

// 	return config, nil
// }

// func LoadConfig2() (Config, error) {
// 	var config Config

// 	// file, err := ioutil.Open(ConfigFile)
// 	// defer file.Close()
// 	// if err != nil {
// 	// 	return config, fmt.Errorf("Error %v", err)
// 	// }

// 	//decoder := json.NewDecoder(file)
// 	// decoder.DisallowUnknownFields() // Force errors

// 	//m := map[string]interface{}{}
// 	// jsonMap := make(map[string]interface{})
// 	// err := json.Unmarshal([]byte(jsonStr), &jsonMap)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// err = decoder.Decode(m)
// 	// if err != nil {
// 	// 	return config, fmt.Errorf("Error %v", err)
// 	// }

// 	jsonFile, _ := os.Open(ConfigFile)
// 	byteValue, _ := ioutil.ReadAll(jsonFile)

// 	var result map[string]interface{}
// 	json.Unmarshal([]byte(byteValue), &result)

// 	for k := range result {

// 		if k == "Test" {

// 		}

// 	}

// 	// data := make(map[string]interface{})
// 	// err := json.Unmarshal([]byte(file), &data)

// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// for k := range m {

// 	// 	fmt.Println("================ ", keyExists("Test", k))
// 	// }

// 	//fmt.Println("================ ", keyExists("Test", m))

// 	//m := map[string]interface{}{}
// 	// if err := json.Unmarshal(file, &m); err != nil {
// 	// 	panic(err)
// 	// }

// 	// if err := decoder.Decode(&config); err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(a)

// 	// jsonParsed, err := gabs.ParseJSONFile(ConfigFile)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// fmt.Println("================ %v", jsonParsed.Path("Test").Data())
// 	// fmt.Println("================ %v", jsonParsed.Path("LiveNvnApplyMaxCount").Data())

// 	// fmt.Println("================ ", jsonParsed.Exists("Test"))

// 	// fmt.Println("================ ", jsonParsed.Exists("Test"))

// 	// conf, ok := jsonParsed.Data().(Config)

// 	return config, nil
// }

// func keyExists(key string, keys []string) bool {
// 	for _, k := range keys {
// 		if k == key {
// 			return true
// 		}
// 	}
// 	return false
// }

// // ConfigSetting ...
// func ConfigSetting(conf Config) int {

// 	tables.LiveNvnApplyMaxCount = conf.LiveNvnApplyMaxCount
// 	beego.Info(fmt.Sprintf("[Config] tables.LiveNvnApplyMaxCount: %v", tables.LiveNvnApplyMaxCount))

// 	//fmt.Println(conf.Test != nil)

// 	return 0
// }
