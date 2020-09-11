package request

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type RequestInfo struct {
	Query url.Values
	Body url.Values
	Header http.Header
	Cookie []*http.Cookie
}

type ResponseInfo struct {
	Body string
	Header http.Header
	Cookie []*http.Cookie
}

type Curl struct {
	Uri string
	Method string
	Request* RequestInfo
	Response *ResponseInfo
	ReTryNum int
	RespondTime time.Duration
	req *http.Request
	resp *http.Response
	err error
}

var DefaultScheme = "http"
var MethodError = errors.New("The Method supports Get AND Post only now")
var UrlError = errors.New("Url error")
var CommonHeader http.Header
var CommonCookie []*http.Cookie

func init() {
	CommonHeader = make(http.Header)
	CommonHeader.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
}

/*
 * 简易 GET 请求
 */
func Get(uri string) (*Curl, error) {
	if uri == "" {
		return nil, UrlError
	}
	c := &Curl{Uri:uri, Request:&RequestInfo{Header:CommonHeader, Cookie:CommonCookie}, Response:&ResponseInfo{}}
	err := c.parseUrl()
	if err != nil {
		return nil, err
	}

	c.Method = http.MethodGet

	return c.base()
}

/*
 * 简易 Post 请求
 */
func Post(uri string, body url.Values) (*Curl, error) {
	if uri == "" {
		return nil, UrlError
	}
	c := &Curl{Uri:uri, Request:&RequestInfo{Body:body, Header:CommonHeader, Cookie:CommonCookie}, Response:&ResponseInfo{}}
	err := c.parseUrl()
	if err != nil {
		return nil, err
	}
	c.Method = http.MethodPost

	return c.base()
}

/*
 * 初始化请求模型
 */
func Client(uri string, method string, body url.Values) (*Curl, error) {
	if http.MethodPost != method && http.MethodGet != method {
		return nil, MethodError
	}
	c := &Curl{Uri: uri, Method:method, Request: &RequestInfo{Body:body, Header:CommonHeader, Cookie:CommonCookie}, Response:&ResponseInfo{}}
	err := c.parseUrl()
	if err != nil {
		c.err = err
		return nil, err
	}

	return c, nil
}

func (c *Curl) Do() (*Curl, error) {
	return c.base()
}

/* *
 * 请求根方法
 */
func (c *Curl) base() (*Curl, error) {
	if c.Request == nil {
		return nil, errors.New("not exists request")
	}

	client := &http.Client{}
	var err error
	c.req, err = http.NewRequest(c.Method, c.combineUrl(), bytes.NewBufferString(c.Request.Body.Encode()))
	fmt.Println(c.combineUrl())

	if err != nil {
		c.err = err
		return nil, err
	}

	if c.Request.Header != nil && len(c.Request.Header) > 0 {
		c.req.Header = c.Request.Header
	}

	c.req.Cookies()
	c.resp, err = client.Do(c.req)
	if err != nil {
		c.err = err
		return nil, err
	}

	c.Response.Header = c.resp.Header
	c.Response.Cookie = c.resp.Cookies()
	b, _ := ioutil.ReadAll(c.resp.Body)
	c.Response.Body = string(b)

	_ = c.resp.Body.Close()

	return c, nil
}

/*
 * 拼接Get参数到url上
 */
func (c *Curl) combineUrl() string {
	if c.Request.Query == nil || len(c.Request.Query) == 0 {
		return c.Uri
	}
	u, _ := url.Parse(c.Uri)
	u.RawQuery = c.Request.Query.Encode()

	return u.String()
}

/*
 * 解析url上的参数
 */
func (c *Curl) parseUrl() error {
	if c.Uri == "" {
		return nil
	}
	u, err := url.Parse(c.Uri)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		u.Scheme = DefaultScheme
	}
	c.Uri = u.Scheme + "://" + u.Host + u.Path
	c.Request.Query = u.Query()
	return nil
}