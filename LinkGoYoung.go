package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/qifengzhang007/goCurl"
	"github.com/tinyhubs/tinydom"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient = goCurl.CreateHttpClient()

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&MyTheme{})
	myWindow := myApp.NewWindow("LinkGoYoung")
	myWindow.Resize(fyne.Size{Width: 320, Height: 480})

	message := "微信公众号:corehub\n本程序为免费程序，请勿相信付费购买"
	msg := widget.NewMultiLineEntry()
	msg.SetText(message)
	label := widget.NewLabel("welcome!")
	userHard := widget.NewSelect([]string{"!^Adcm0", "!^Iqnd0"}, func(value string) {})

	var user userInfo
	userHard.SetSelectedIndex(0)
	accountEntry := widget.NewEntry()
	accountEntry.SetPlaceHolder("账号")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("密码")
	user.ReadUserInfoJson(label)
	userHard.SetSelected(user.UserHard)
	accountEntry.SetText(user.UserAccount)
	passwordEntry.SetText(user.PassWord)
	label.SetText(CheckServer("baidu.com:443"))

	form := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "账号", Widget: accountEntry}},
		OnSubmit: func() {
			loginURL, userIp, nasIp, userMac := nextUrl(msg, message, "http://www.msftconnecttest.com/redirect")
			fmt.Printf("校园网关IP：%s\n", nasIp)
			if nasIp == "" {
				label.SetText("重定向失败，请检查网络环境")
				return
			}
			message = "当前网络信息："
			message += "\nnasip:" + nasIp
			message += "\nuserip:" + userIp
			message += "\nusermac:" + userMac
			loginURL = Get(loginURL)
			loginURL = parseXML(loginURL, "WISPAccessGatewayParam", "Redirect", "LoginURL")
			fmt.Println("loginURL:" + loginURL)
			dialog.ShowInformation("登陆……", login(msg, message, label, loginURL, time.Now().Format("20060102150405"),
				userHard.Selected, accountEntry.Text, passwordEntry.Text), myWindow)
			if label.Text == "50：认证成功" {
				user.UserHard = userHard.Selected
				user.UserAccount = accountEntry.Text
				user.PassWord = passwordEntry.Text
				user.LastLoginURL = loginURL
				user.SaveUserInfoJson(label)
			}
		},
		OnCancel: func() {
			u := user.LastLoginURL
			dialog.ShowInformation("下线……", logout(label, u), myWindow)
		},
		SubmitText: "登录",
		CancelText: "下线",
	}

	form.Append("密码", passwordEntry)

	content := container.New(layout.NewGridLayoutWithColumns(1), msg, container.NewVBox(
		label,
		userHard,
		form,
	))

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func login(msg *widget.Entry, message string, label *widget.Label, url, dateTime, userHard, userName, passWord string) string {
	token := "UserName=" + userHard + userName + "&Password=" + passWord + "&AidcAuthAttr1=" + dateTime +
		"&AidcAuthAttr3=keuyGQlK&AidcAuthAttr4=zrDgXllCChyJHjwkcRwhygP0&AidcAuthAttr5=kfe1GQhXdGqOFDtee" +
		"go5zwP9IsNoxX7djTWspPrYm1A%3D%3D&AidcAuthAttr6=5Ia4cQhDfXSFbTtUDGY1yx8%3D&AidcAuthAttr7=6ZWiVl" +
		"wdNiHMXCpOagQv2w2MQs0ohTWJnTu8qK5OibhCydTpTxkI88wadKPWby%2F2PKCVaZUxglbBs96%2FtmLE89M8AJ6y28o7" +
		"qolpFep%2FcYFFRLd7H4MAMrDUMRO0F%2B93jh14fiAZYmtk9hdp%2BZ5w%2BjMQUoV4TCtM9VJ07XQwxlMVg%2F0YKrS1" +
		"s3hXAstdQ1fvdSn3nAVGgdxc%2BJQDrQ%3D%3D&AidcAuthAttr8=jPSyBQxVaXWTQWUaakluj06scJ98nyqCyX7y%2FLU" +
		"k1OkXiNjkXhVGvJhyTuLDaCPhK%2FOFJttlxxiVqNKupnDXkp9%2BR9D9j8p2j5h8FOxoatMaGu0oRdk%3D&createAuthorFlag=0"
	resp, err := httpClient.Post(url, goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent":   "CDMA+WLAN(Mios)",
			"Content-Type": "application/x-www-form-urlencoded",
		},
		XML:           token,
		SetResCharset: "utf-8",
	})
	if err != nil {
		message += fmt.Sprintf("\nLogin请求出错：%s", err.Error())
		msg.SetText(message)
		return "Login请求出错"
	}
	body, err := resp.GetContents()
	if err != nil {
		message += fmt.Sprintf("\nLogin请求失败,错误明细：%s", err.Error())
		msg.SetText(message)
		return "Login请求出错"
	}
	body = parseXML(body, "WISPAccessGatewayParam", "AuthenticationReply", "ReplyMessage")
	message += fmt.Sprintf("\n请求结果：%s", body)
	msg.SetText(message)
	label.SetText(body)
	return body
}

func logout(label *widget.Label, url string) string {
	url = "http://58.53.199.144:8001/wispr_logout.jsp?" + strings.Split(url, "?")[1]
	ReplyMessage := Get(url)
	ReplyMessage = parseXML(ReplyMessage, "WISPAccessGatewayParam", "LogoffReply", "ResponseCode")
	if ReplyMessage == "150" {
		label.SetText("150:下线成功")
		return "150:下线成功"
	} else {
		label.SetText("255:下线失败")
		return "255:下线失败"
	}
}

func nextUrl(msg *widget.Entry, message, urlIn string) (urlOut, userIp, nasIp, userMac string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}
	res, err := client.Get(urlIn)
	if err != nil {
		message += fmt.Sprintf("\n重定向请求出错：%s", err.Error())
		msg.SetText(message)
		return
	}
	if res.StatusCode != http.StatusFound {
		message += fmt.Sprintf("\n[Error]StatusCode:%v", res.StatusCode)
		msg.SetText(message)
		return
	}
	u, _ := url.Parse(res.Header.Get("Location"))
	query := u.Query()
	userIp, nasIp, userMac = query.Get("userip"), query.Get("nasip"), query.Get("usermac")
	urlOut = fmt.Sprintf("http://58.53.199.144:8001/?userip=%s&wlanacname=&nasip=%s&usermac=%s&aidcauthtype=0", userIp, nasIp, userMac)
	return urlOut, userIp, nasIp, userMac
}

func Get(url string) string {
	resp, err := httpClient.Get(url, goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent": "CDMA+WLAN(Mios)",
		},
		SetResCharset: "utf-8",
	})
	if err != nil && resp == nil {
		return ""
	} else {
		body, err := resp.GetContents()
		if err != nil {
			return ""
		}
		return body
	}
}

func parseXML(str, str1, str2, str3 string) string {
	doc, err := tinydom.LoadDocument(strings.NewReader(str))
	if err != nil {
		fmt.Printf("XML解析失败，err：%s\n", err.Error())
	}
	elem := doc.FirstChildElement(str1).FirstChildElement(str2).FirstChildElement(str3)
	return elem.Text()
}

func CheckServer(url string) string {
	timeout := 5 * time.Second
	//t1 := time.Now()
	_, err := net.DialTimeout("tcp", url, timeout)
	//massage += "\n网络测试时长 :" + time.Now().Sub(t1).String()

	if err == nil {
		return "已接入互联网，只能进行下线操作"
	} else {
		return "未接入互联网"
	}
}
