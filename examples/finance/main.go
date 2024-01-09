package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/infinitete/dynamic"
	"github.com/xuri/excelize/v2"
)

type Provider struct {
	Company  string `xlsx:"col:公司名称"`
	Surveyor string `xlsx:"col:叫车人"`
}

type Servant struct {
	Company string `xlsx:"col:公司名称"`
	Driver  string `xlsx:"col:驾驶员"`
}

// WechatPay
// 微信支付
type WechatPay struct {
	Receivable int    `xlsx:"col:应收(元)"`
	Received   int    `xlsx:"col:已收(元)"`
	Datetime   string `xlsx:"col:时间"`
}

// ProviderPay
// 直赔
type ProviderPay struct {
	Receivable  int     `xlsx:"col:应收(元)"`
	Received    int     `xlsx:"col:已收(元)"`
	NotReceived int     `xlsx:"col:未收(元)"`
	Percent     float64 `xlsx:"col:到款率"`
	Datetime    string  `xlsx:"col:时间"`
}

// CashPay
// 现金收款
type CashPay struct {
	Fee      int    `xlsx:"col:现金收款(元)"`
	Datetime string `xlsx:"现金收款时间"`
}

type FeeDetails struct {
	WechatPay       WechatPay   `xlsx:"col:微信支付"`
	ProviderPay     ProviderPay `xlsx:"col:直赔"`
	CashPay         int         `xlsx:"col:现金收款(元)"`
	CashPayDatetime string      `xlsx:"col:现金收款时间"`
}

// Fee
// 应收已收
type Fee struct {
	Receivable  int        `xlsx:"col:应收(元)"`
	Received    int        `xlsx:"col:已收(元)"`
	NotReceived int        `xlsx:"col:未收(元)"`
	Delay       int        `xlsx:"col:超期天数"`
	Details     FeeDetails `xlsx:"col:明细"`
}

// Paid
// 向下付款
type Paid struct {
	Pay      int    `xlsx:"col:应付(元)"`
	Paid     int    `xlsx:"col:已付(元)"`
	NotPaid  int    `xlsx:"col:未付(元)"`
	Datetime string `xlsx:"col:付款时间"`
}

type Sheet struct {
	OutNo          string   `xlsx:"col:报案号"`
	License        string   `xlsx:"col:车牌号"`
	RescueDate     string   `xlsx:"col:救援日期"`
	Area           string   `xlsx:"col:区域"`
	Provider       Provider `xlsx:"col:叫车方"`
	Servant        Servant  `xlsx:"col:服务人"`
	Fee            Fee      `xlsx:"col:应收/已收"`
	Paid           Paid     `xlsx:"col:向下付款"`
	Invoice        string   `xlsx:"col:开票"`
	ServiceRemarks string   `xlsx:"col:服务备注"`
	FinanceRemarks string   `xlsx:"col:财务备注"`
	ReviewDatetime string   `xlsx:"col: w审核日期"`
	RescueMode     string   `xlsx:"col:救援模式"`
}

func main() {
	file, err := excelize.OpenFile("./财务报表.xlsx")
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	reader, err := dynamic.NewReader[Sheet]()
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	data := reader.Read(file, "财务报表")
	var totaReceivable = 0
	var totalWechat = 0
	var totalProvider = 0

	for _, node := range data {
		totaReceivable += node.Fee.Receivable
		totalWechat += node.Fee.Details.WechatPay.Receivable
		totalProvider += node.Fee.Details.ProviderPay.Receivable
	}

	log.Printf("总应收: %d, 微信应收: %d, 直赔应收: %d", totaReceivable, totalWechat, totalProvider)

	for _, node := range data {
		if node.OutNo == "SDAA2024520110S0001298" {
			b, _ := json.MarshalIndent(node, "", "  ")
			fmt.Printf("\n%s\n", b)
		}
	}
}
