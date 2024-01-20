package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
)

var (
	debugClientHttp bool
)

func init() {
	debugClientHttp = beego.AppConfig.DefaultBool("Debug.ClientHttp", true)
}

func HttpGetByParam(urlparam, timeo string) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHttp)
	timeout, _ := time.ParseDuration(timeo)
	if urlparam[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	resp, body, errs := request.Get(urlparam).
		Timeout(timeout).
		Retry(3, time.Second, http.StatusInternalServerError).
		EndBytes()
	if errs == nil {
		if resp.StatusCode == 200 {
			return body, nil
		} else {
			return body, errors.New(fmt.Sprintln("Response Http: %v ", resp.StatusCode))
		}
	}
	fmt.Printf(" Error : %v \n", errs)
	return body, errors.New("Connection Problem") // , nil
}

func HttpPost(url string, jsondata interface{}, timeo string, retry int) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHttp)
	timeout, _ := time.ParseDuration(timeo)
	err := errors.New("Connection Problem")
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	resp, body, errs := request.Post(url).
		Set("Accept", "*/*").
		Set("Content-Type", "application/json").
		//Set("Content-Type", "application/x-www-form-urlencoded").
		Send(jsondata).
		Timeout(timeout).
		Retry(retry, time.Second, http.StatusInternalServerError).
		End()
	if errs == nil {
		if resp.StatusCode/100 == 2 {
			return []byte(body), nil
		}
		return []byte(body), nil
	}
	return []byte(body), err // , nil
}

func HttpPostByParam(urlparam, timeo string) ([]byte, error) {
	fmt.Println("Sending ")
	request := gorequest.New()
	request.SetDebug(debugClientHttp)
	timeout, _ := time.ParseDuration(timeo)
	if urlparam[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	_, body, errs := request.Post(urlparam).
		Timeout(timeout).
		Retry(1, time.Second, http.StatusInternalServerError).
		EndBytes()
	if errs == nil {
		//if resp.StatusCode == 200 {
		return body, nil
		//}
		//return body, errors.New(fmt.Sprintln("Response Http: %v ", resp.StatusCode))
	}
	logs.Error("HttpPostByParam", errs)
	//fmt.Printf(" Error : %v \n", errs)
	return body, errors.New("Connection Problem") // , nil
}

func HttpPostWithHeader(url string, jsondata interface{}, header http.Header, timeo string, retry int) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHttp)
	timeout, _ := time.ParseDuration(timeo)
	//_ := errors.New("Connection Problem")

	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	reqagent := request.Post(url)
	// reqagent.Header.Set("token", Token)
	reqagent.Header = header
	_, body, errs := reqagent.
		//Set("Accept", "*/*").
		//Set("Content-Type", "application/json").
		Send(jsondata).
		Timeout(timeout).
		Retry(retry, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}
