// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"fmt"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"gopkg.in/ini.v1"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	//# 将accessKeyId改成自己的accessKeyId
	AccessKeyId string
	//# 将accessSecret改成自己的accessSecret
	AccessSecret string
	//# 是否开启ipv4 ddns解析,1为开启，0为关闭
	Ipv4Flag int
	//# 是否开启ipv6 ddns解析,1为开启，0为关闭
	Ipv6Flag int
	//# 你的主域名
	Domain string
	//# 要进行ipv4 ddns解析的子域名
	NameIpv4 string
	//# 要进行ipv6 ddns解析的子域名
	NameIpv6 string
	//# 日志同时输出到文件,1为开启，0为关闭
	LogFileFlag int
	//# 探测本机公共ip的服务器,可以不配置,有内部默认值,默认服务器挂了,可以设置新的服务器
	Ip4Url string //= https://api-ipv4.ip.sb/ip
	Ip6Url string //= https://api6.ipify.org
}

// 全局配置
var config Config
var programDir string

func Log(a interface{}) (n int, err error) {

	var timeStr = time.Now().Format("[2006-01-02 15:04:05] ")
	fmt.Print(timeStr)

	switch v := a.(type) {
	case *string:
		return fmt.Fprintln(os.Stdout, *v)
	case string:
		return fmt.Fprintln(os.Stdout, v)
	}
	return fmt.Fprintln(os.Stdout, a)
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("alidns.cn-shanghai.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

func _getRecords(client *alidns20150109.Client) (
	Record []*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord, _err error) {

	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(config.Domain),
		PageNumber: tea.Int64(1),
		PageSize:   tea.Int64(500),
	}

	runtime := &util.RuntimeOptions{}
	resp, _err := client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, runtime)
	if _err != nil {
		return nil, _err
	}

	var ls = resp.Body.DomainRecords.Record

	//Log(util.ToJSONString(tea.ToMap(resp)))
	return ls, _err
}

/*
*
{"DomainName":"baise.tk","Line":"default","Locked":false,"RR":"tb","RecordId":"768988894054603776","Status":"ENABLE",
"TTL":600,"Type":"AAAA","Value":"2409:8a62:3e2:f2c0:799c:d0cf:6423:3
bc4","Weight":1}
*/
func getRecordIdByPR(Record []*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord,
	PR string, ip string, isV4 bool) (id *string, _err error) {

	if Record == nil {
		//添加
		return nil, nil
	}

	var aType = "AAAA"
	if isV4 {
		aType = "A"
	}

	for _, value := range Record {
		//fmt.Println("Index =", index, "Value =", value)
		if *value.RR == PR && *value.Type == aType {
			if *value.Value == ip {
				return value.RecordId, fmt.Errorf("ip地址没变")
			}
			//修改
			return value.RecordId, nil
		}
	}
	//添加
	return nil, nil
}

func _add(client *alidns20150109.Client, isV4 bool, ip string) (_err error) {

	var name = config.NameIpv6
	var aType = "AAAA"
	if isV4 {
		name = config.NameIpv4
		aType = "A"
	}

	Log("添加解析: " + name + "." + config.Domain)

	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(config.Domain),
		RR:         tea.String(name),
		Type:       tea.String(aType),
		Value:      tea.String(ip),
	}
	runtime := &util.RuntimeOptions{}
	_, _err = client.AddDomainRecordWithOptions(addDomainRecordRequest, runtime)
	if _err != nil {

		return _err
	}

	//Log(util.ToJSONString(tea.ToMap(resp)))

	return nil
}

func _update(client *alidns20150109.Client, isV4 bool, ip string, RecordId string) (_err error) {
	var name = config.NameIpv6
	var aType = "AAAA"
	if isV4 {
		name = config.NameIpv4
		aType = "A"
	}

	Log("更新解析: " + name + "." + config.Domain)

	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(RecordId),
		RR:       tea.String(name),
		Type:     tea.String(aType),
		Value:    tea.String(ip),
	}
	runtime := &util.RuntimeOptions{}
	_, _err = client.UpdateDomainRecordWithOptions(updateDomainRecordRequest, runtime)
	if _err != nil {
		return _err
	}

	//Log(util.ToJSONString(tea.ToMap(resp)))

	return nil
}

func getIniConfig() (_err error) {
	//默认值
	config.Ipv4Flag = 1
	config.Ipv6Flag = 1
	config.LogFileFlag = 1
	config.Ip4Url = "https://api-ipv4.ip.sb/ip"
	config.Ip6Url = "https://api6.ipify.org"

	path := "aliddns.ini"
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			//工作目录下面没有,就拼接路径到程序目录; 还没有就报错
			path = programDir + string(os.PathSeparator) + path
		}
	}

	load, err := ini.Load(path)
	if err != nil {
		Log("配置文件读取失败:" + path)
		return err
	}

	section, err := load.GetSection("DNS")
	if err != nil {
		return err
	}
	{
		key, err := section.GetKey("AccessKeyId")
		if err != nil {
			return err
		}
		config.AccessKeyId = strings.TrimSpace(key.String())
	}
	{
		key, err := section.GetKey("AccessSecret")
		if err != nil {
			return err
		}
		config.AccessSecret = strings.TrimSpace(key.String())
	}
	{
		key, err := section.GetKey("Domain")
		if err != nil {
			return err
		}
		config.Domain = strings.TrimSpace(key.String())
	}
	{
		key, err := section.GetKey("NameIpv4")
		if err != nil {
			return err
		}
		config.NameIpv4 = strings.TrimSpace(key.String())
	}
	{
		key, err := section.GetKey("NameIpv6")
		if err != nil {
			return err
		}
		config.NameIpv6 = strings.TrimSpace(key.String())
	}
	{
		key, err := section.GetKey("Ipv4Flag")
		if err != nil {
			return err
		}
		config.Ipv4Flag, _ = key.Int()
	}
	{
		key, err := section.GetKey("Ipv6Flag")
		if err != nil {
			return err
		}
		config.Ipv6Flag, _ = key.Int()
	}
	{
		key, err := section.GetKey("LogFileFlag")
		if err != nil {
			return err
		}
		config.LogFileFlag, _ = key.Int()
	}
	{
		key, err := section.GetKey("Ip4Url")
		if err != nil {
			return err
		}
		var t = strings.TrimSpace(key.String())
		if len(t) > 0 {
			config.Ip4Url = t
		}
	}
	{
		key, err := section.GetKey("Ip6Url")
		if err != nil {
			return err
		}
		var t = strings.TrimSpace(key.String())
		if len(t) > 0 {
			config.Ip6Url = t
		}

	}
	Log(util.ToJSONString(config))

	return nil
}

func getPublicIp(url string) (ip string, _err error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	//刚才找到的浏览器中User-agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		//console.Log(tea.String(err.Error()))
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		//Log("取得ip地址:" + string(body))
	}

	ip = strings.TrimSpace(string(body))
	Log("取得ip地址:" + ip)
	if len(ip) > 0 {
		return ip, nil
	} else {
		return "", fmt.Errorf("ip获取失败")
	}
}

func _main(args []*string) (_err error) {

	_err = getIniConfig()
	if _err != nil {
		return _err
	}

	if config.Ipv4Flag == 0 && config.Ipv6Flag == 0 {
		Log(tea.String("ipv4,ipv6的更新都没有开启,退出"))
		return nil
	}

	client, _err := CreateClient(tea.String(config.AccessKeyId), tea.String(config.AccessSecret))
	if _err != nil {
		return _err
	}

	records, _err := _getRecords(client)
	if _err != nil {
		return _err
	}

	if config.Ipv4Flag == 1 {
		ip, _err := getPublicIp(config.Ip4Url)
		if _err != nil {
			Log("获取ip4失败:" + _err.Error())
		} else {

			id, _err := getRecordIdByPR(records, config.NameIpv4, ip, true)
			if _err != nil {
				Log(_err.Error())
			} else if id == nil {
				//add
				_err := _add(client, true, ip)
				if _err != nil {
					Log("添加ip4失败:" + _err.Error())
				} else {
					Log("添加ip4成功")
				}
			} else {
				//edit
				_err := _update(client, true, ip, *id)
				if _err != nil {
					Log("更新ip4失败:" + _err.Error())
				} else {
					Log("更新ip4成功")
				}
			}

		}
	}
	if config.Ipv6Flag == 1 {
		ip, _err := getPublicIp(config.Ip6Url)
		if _err != nil {
			Log("获取ip6失败:" + _err.Error())
		} else {

			id, _err := getRecordIdByPR(records, config.NameIpv6, ip, false)
			if _err != nil {
				Log(_err.Error())
			} else if id == nil {
				//add
				_err := _add(client, false, ip)
				if _err != nil {
					Log("添加ip6失败:" + _err.Error())
				} else {
					Log("添加ip6成功")
				}
			} else {
				//edit
				_err := _update(client, false, ip, *id)
				if _err != nil {
					Log("更新ip6失败:" + _err.Error())
				} else {
					Log("添加ip6成功")
				}
			}

		}
	}

	return nil
}

func main() {
	{
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			//log.Fatal(err)
		}
		//fmt.Println(dir)
		programDir = dir

		/*
			u, err := user.Current()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(u)
		*/

		//获取 go 工作目录
		//var home string = os.Getenv("GOROOT")
		//fmt.Printf("GO 工作目录是 %s\n", home)
		//获取 go 项目目录
		//path := os.Getenv("GOPATH")
		//fmt.Printf("GO 项目目录是 %s\n", path)

	}
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
