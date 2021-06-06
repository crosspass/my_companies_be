package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/my-companies-be/models"
	"github.com/my-companies-be/utils"
)

var dsn = "host=localhost user=wu password=gorm dbname=my_companies port=5432 sslmode=disable TimeZone=Asia/Shanghai"
var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

type CompanyReq struct {
	Symbol string
	Name   string
}

type CompanyReqData struct {
	Count uint
	List  []CompanyReq
}

type CompanyReqstruct struct {
	Data             CompanyReqData
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

func getCompanies() {
	client := &http.Client{}
	const path string = "https://xueqiu.com/service/v5/stock/screener/quote/list"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	types := [3]string{"sha", "sza", "cyb"}

	for _, atype := range types {
		page := 0
		for {
			q := req.URL.Query()
			q.Set("type", atype)
			q.Set("order_by", "symbol")
			q.Set("order", "asc")
			q.Set("size", "200")
			page++
			q.Set("page", strconv.Itoa(page))
			req.URL.RawQuery = q.Encode()

			fmt.Println(req.URL.String())
			resp, err := client.Do(req)
			fmt.Println(resp.Status)

			if err != nil {
				fmt.Println("Errored when sending request to the server")
				return
			}
			var companyReq CompanyReqstruct

			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&companyReq)
			if len(companyReq.Data.List) == 0 {
				break
			}
			fmt.Printf("error_code: %d\n", companyReq.ErrorCode)
			for _, companyReq := range companyReq.Data.List {
				var company models.Company
				// fmt.Println(companyReq.Symbol)
				db.Where("code = ?", companyReq.Symbol).First(&company)
				if company.Code != "" {
					if company.Name != companyReq.Name {
						fmt.Println("update compnay name")
						db.Model(&company).Update("name", companyReq.Name)
					}
				} else {
					fmt.Println("create company")
					fmt.Println(companyReq)
					db.Create(&models.Company{Name: companyReq.Name, Code: companyReq.Symbol})
				}
			}
			time.Sleep(time.Duration(2) * time.Second)
		}
	}
}

type ReportSummary struct {
	gorm.Model
	Category                            string
	CompanyCode                         string
	ReportName                          string
	ReportDate                          uint64
	AvgRoe                              float32
	AvgRoeIncrease                      float32
	NpPerShare                          float32
	NpPerShareIncrease                  float32
	OperateCashFlowPs                   float32
	OperateCashFlowPsIncrease           float32
	BasicEps                            float32
	BasicEpsIncrease                    float32
	CapitalReserve                      float32
	CapitalReserveIncrease              float32
	UndistriProfitPs                    float32
	UndistriProfitPsIncrease            float32
	NetInterestOfTotalAssets            float32
	NetInterestOfTotalAssetsIncrease    float32
	NetSellingRate                      float32
	NetSellingRateIncrease              float32
	GrossSellingRate                    float32
	GrossSellingRateIncrease            float32
	TotalRevenue                        float32
	TotalRevenueIncrease                float32
	OperatingIncomeYoy                  float32
	OperatingIncomeYoyIncrease          float32
	NetProfitAtsopc                     float32
	NetProfitAtsopcIncrease             float32
	NetProfitAtsopcYoy                  float32
	NetProfitAtsopcYoyIncrease          float32
	NetProfitAfterNrgalAtsolc           float32
	NetProfitAfterNrgalAtsolcIncrease   float32
	NpAtsopcNrgalYoy                    float32
	NpAtsopcNrgalYoyIncrease            float32
	OreDlt                              float32
	OreDltIncrease                      float32
	Rop                                 float32
	RopIncrease                         float32
	AssetLiabRatio                      float32
	AssetLiabRatioIncrease              float32
	CurrentRatio                        float32
	CurrentRatioIncrease                float32
	QuickRatio                          float32
	QuickRatioIncrease                  float32
	EquityMultiplier                    float32
	EquityMultiplierIncrease            float32
	EquityRatio                         float32
	EquityRatioIncrease                 float32
	HolderEquity                        float32
	HolderEquityIncrease                float32
	NcfFromOaToTotalLiab                float32
	NcfFromOaToTotalLiabIncrease        float32
	InventoryTurnoverDays               float32
	InventoryTurnoverDaysIncrease       float32
	ReceivableTurnoverDays              float32
	ReceivableTurnoverDaysIncrease      float32
	AccountsPayableTurnoverDays         float32
	AccountsPayableTurnoverDaysIncrease float32
	CashCycle                           float32
	CashCycleIncrease                   float32
	OperatingCycle                      float32
	OperatingCycleIncrease              float32
	TotalCapitalTurnover                float32
	TotalCapitalTurnoverIncrease        float32
	InventoryTurnover                   float32
	InventoryTurnoverIncrease           float32
	AccountReceivableTurnover           float32
	AccountReceivableTurnoverIncrease   float32
	AccountsPayableTurnover             float32
	AccountsPayableTurnoverIncrease     float32
	CurrentAssetTurnoverRate            float32
	CurrentAssetTurnoverRateIncrease    float32
	FixedAssetTurnoverRatio             float32
	FixedAssetTurnoverRatioIncrease     float32
}

type ReportSummaryReq struct {
	ReportDate                  uint64     `json:"report_date"`
	ReportName                  string     `json:"report_name"`
	AvgRoe                      [2]float32 `json:"avg_roe"`
	NpPerShare                  [2]float32 `json:"np_per_share"`
	OperateCashFlowPs           [2]float32 `json:"operate_cash_flow_ps"`
	BasicEps                    [2]float32 `json:"basic_eps"`
	CapitalReserve              [2]float32 `json:"capital_reserve"`
	UndistriProfitPs            [2]float32 `json:"undistri_profit_ps"`
	NetInterestOfTotalAssets    [2]float32 `json:"net_interest_of_total_assets"`
	NetSellingRate              [2]float32 `json:"net_selling_rate"`
	GrossSellingRate            [2]float32 `json:"gross_selling_rate"`
	TotalRevenue                [2]float32 `json:"total_revenue"`
	OperatingIncomeYoy          [2]float32 `json:"operating_income_yoy"`
	NetProfitAtsopc             [2]float32 `json:"net_profit_atsopc"`
	NetProfitAtsopcYoy          [2]float32 `json:"net_profit_atsopc_yoy"`
	NetProfitAfterNrgalAtsolc   [2]float32 `json:"net_profit_after_nrgal_atsolc"`
	NpAtsopcNrgalYoy            [2]float32 `json:"np_atsopc_nrgal_yoy"`
	OreDlt                      [2]float32 `json:"ore_dlt"`
	Rop                         [2]float32 `json:"rop"`
	AssetLiabRatio              [2]float32 `json:"asset_liab_ratio"`
	CurrentRatio                [2]float32 `json:"current_ratio"`
	QuickRatio                  [2]float32 `json:"quick_ratio"`
	EquityMultiplier            [2]float32 `json:"equity_multiplier"`
	EquityRatio                 [2]float32 `json:"equity_ratio"`
	HolderEquity                [2]float32 `json:"holder_equity"`
	NcfFromOaToTotalLiab        [2]float32 `json:"ncf_from_oa_to_total_liab"`
	InventoryTurnoverDays       [2]float32 `json:"inventory_turnover_days"`
	ReceivableTurnoverDays      [2]float32 `json:"receivable_turnover_days"`
	AccountsPayableTurnoverDays [2]float32 `json:"accounts_payable_turnover_days"`
	CashCycle                   [2]float32 `json:"cash_cycle"`
	OperatingCycle              [2]float32 `json:"operating_cycle"`
	TotalCapitalTurnover        [2]float32 `json:"total_capital_turnover"`
	InventoryTurnover           [2]float32 `json:"inventory_turnover"`
	AccountReceivableTurnover   [2]float32 `json:"account_receivable_turnover"`
	AccountsPayableTurnover     [2]float32 `json:"accounts_payable_turnover"`
	CurrentAssetTurnoverRate    [2]float32 `json:"current_asset_turnover_rate"`
	FixedAssetTurnoverRatio     [2]float32 `json:"fixed_asset_turnover_ratio"`
}
type ReportSummaryDataReq struct {
	QuoteName      string
	LastReportRame string
	List           []ReportSummaryReq
}

type ReportSummaryReqstruct struct {
	Data             ReportSummaryDataReq
	ErrorCode        uint
	ErrorDescription string
}

func fetchQ(str string) string {
	switch {
	case strings.Contains(str, "一"):
		return "Q1"
	case strings.Contains(str, "二"):
		return "Q2"
	case strings.Contains(str, "三"):
		return "Q3"
	case strings.Contains(str, "年"):
		return "Q4"
	default:
		return "Q"
	}
}

func getClient() *http.Client {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar: jar,
	}
	homereq, err := http.NewRequest("GET", "https://xueqiu.com", nil)
	homereq.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(homereq.URL.String())
	resp, err := client.Do(homereq)
	if err != nil {
		log.Fatal("Errored when sending request to the server")
	}
	fmt.Println(resp.Status)
	return client
}

func getReportSummary() {
	var companies []models.Company
	var reportSummaries []ReportSummary
	var companyCodes []string
	var reportCodes []string
	db.Table("companies").Select("code").Find(&companies)
	db.Distinct("company_code").Select("company_code").Find(&reportSummaries)

	for _, company := range companies {
		companyCodes = append(companyCodes, company.Code)
	}

	for _, report := range reportSummaries {
		reportCodes = append(reportCodes, report.CompanyCode)
	}

	codes := utils.Difference(companyCodes, reportCodes)
	fmt.Println("codes", codes)

	client := getClient()
	for _, code := range codes {
		const path string = "http://stock.xueqiu.com/v5/stock/finance/cn/indicator.json?type=ALL&is_detail=true&count=100"
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			log.Fatal(err)
		}
		q := req.URL.Query()
		q.Set("symbol", code)
		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal("Errored when sending request to the server")
		}

		fmt.Println(resp.Status)
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal(string(body))
		}

		var reportSummaryReqstruct ReportSummaryReqstruct

		json.NewDecoder(resp.Body).Decode(&reportSummaryReqstruct)
		fmt.Println(reportSummaryReqstruct)
		for _, reportSummaryReq := range reportSummaryReqstruct.Data.List {
			var reportSummary ReportSummary
			db.Where("company_code = ? AND report_name = ?", code, reportSummaryReq.ReportName).First(&reportSummary)
			if reportSummary.ReportName == "" {
				reportSummary := ReportSummary{
					Category:                            fetchQ(reportSummaryReq.ReportName),
					CompanyCode:                         code,
					ReportName:                          reportSummaryReq.ReportName,
					ReportDate:                          reportSummaryReq.ReportDate,
					AvgRoe:                              reportSummaryReq.AvgRoe[0],
					AvgRoeIncrease:                      reportSummaryReq.AvgRoe[1],
					NpPerShare:                          reportSummaryReq.NpPerShare[0],
					NpPerShareIncrease:                  reportSummaryReq.NpPerShare[1],
					OperateCashFlowPs:                   reportSummaryReq.OperateCashFlowPs[0],
					OperateCashFlowPsIncrease:           reportSummaryReq.OperateCashFlowPs[1],
					BasicEps:                            reportSummaryReq.BasicEps[0],
					BasicEpsIncrease:                    reportSummaryReq.BasicEps[1],
					CapitalReserve:                      reportSummaryReq.CapitalReserve[0],
					CapitalReserveIncrease:              reportSummaryReq.CapitalReserve[1],
					UndistriProfitPs:                    reportSummaryReq.UndistriProfitPs[0],
					UndistriProfitPsIncrease:            reportSummaryReq.UndistriProfitPs[1],
					NetInterestOfTotalAssets:            reportSummaryReq.NetInterestOfTotalAssets[0],
					NetInterestOfTotalAssetsIncrease:    reportSummaryReq.NetInterestOfTotalAssets[1],
					NetSellingRate:                      reportSummaryReq.NetSellingRate[0],
					NetSellingRateIncrease:              reportSummaryReq.NetSellingRate[1],
					GrossSellingRate:                    reportSummaryReq.GrossSellingRate[0],
					GrossSellingRateIncrease:            reportSummaryReq.GrossSellingRate[1],
					TotalRevenue:                        reportSummaryReq.TotalRevenue[0],
					TotalRevenueIncrease:                reportSummaryReq.TotalRevenue[1],
					OperatingIncomeYoy:                  reportSummaryReq.OperatingIncomeYoy[0],
					OperatingIncomeYoyIncrease:          reportSummaryReq.OperatingIncomeYoy[1],
					NetProfitAtsopc:                     reportSummaryReq.NetProfitAtsopc[0],
					NetProfitAtsopcIncrease:             reportSummaryReq.NetProfitAtsopc[1],
					NetProfitAtsopcYoy:                  reportSummaryReq.NetProfitAtsopcYoy[0],
					NetProfitAtsopcYoyIncrease:          reportSummaryReq.NetProfitAtsopcYoy[1],
					NetProfitAfterNrgalAtsolc:           reportSummaryReq.NetProfitAfterNrgalAtsolc[0],
					NetProfitAfterNrgalAtsolcIncrease:   reportSummaryReq.NetProfitAfterNrgalAtsolc[1],
					NpAtsopcNrgalYoy:                    reportSummaryReq.NpAtsopcNrgalYoy[0],
					NpAtsopcNrgalYoyIncrease:            reportSummaryReq.NpAtsopcNrgalYoy[1],
					OreDlt:                              reportSummaryReq.OreDlt[0],
					OreDltIncrease:                      reportSummaryReq.OreDlt[1],
					Rop:                                 reportSummaryReq.Rop[0],
					RopIncrease:                         reportSummaryReq.Rop[1],
					AssetLiabRatio:                      reportSummaryReq.AssetLiabRatio[0],
					AssetLiabRatioIncrease:              reportSummaryReq.AssetLiabRatio[1],
					CurrentRatio:                        reportSummaryReq.CurrentRatio[0],
					CurrentRatioIncrease:                reportSummaryReq.CurrentRatio[1],
					QuickRatio:                          reportSummaryReq.QuickRatio[0],
					QuickRatioIncrease:                  reportSummaryReq.QuickRatio[1],
					EquityMultiplier:                    reportSummaryReq.EquityMultiplier[0],
					EquityMultiplierIncrease:            reportSummaryReq.EquityMultiplier[1],
					EquityRatio:                         reportSummaryReq.EquityRatio[0],
					EquityRatioIncrease:                 reportSummaryReq.EquityRatio[1],
					HolderEquity:                        reportSummaryReq.HolderEquity[0],
					HolderEquityIncrease:                reportSummaryReq.HolderEquity[1],
					NcfFromOaToTotalLiab:                reportSummaryReq.NcfFromOaToTotalLiab[0],
					NcfFromOaToTotalLiabIncrease:        reportSummaryReq.NcfFromOaToTotalLiab[1],
					InventoryTurnoverDays:               reportSummaryReq.InventoryTurnoverDays[0],
					InventoryTurnoverDaysIncrease:       reportSummaryReq.InventoryTurnoverDays[1],
					ReceivableTurnoverDays:              reportSummaryReq.ReceivableTurnoverDays[0],
					ReceivableTurnoverDaysIncrease:      reportSummaryReq.ReceivableTurnoverDays[1],
					AccountsPayableTurnoverDays:         reportSummaryReq.AccountsPayableTurnoverDays[0],
					AccountsPayableTurnoverDaysIncrease: reportSummaryReq.AccountsPayableTurnoverDays[1],
					CashCycle:                           reportSummaryReq.CashCycle[0],
					CashCycleIncrease:                   reportSummaryReq.CashCycle[1],
					OperatingCycle:                      reportSummaryReq.OperatingCycle[0],
					OperatingCycleIncrease:              reportSummaryReq.OperatingCycle[1],
					TotalCapitalTurnover:                reportSummaryReq.TotalCapitalTurnover[0],
					TotalCapitalTurnoverIncrease:        reportSummaryReq.TotalCapitalTurnover[1],
					InventoryTurnover:                   reportSummaryReq.InventoryTurnover[0],
					InventoryTurnoverIncrease:           reportSummaryReq.InventoryTurnover[1],
					AccountReceivableTurnover:           reportSummaryReq.AccountReceivableTurnover[0],
					AccountReceivableTurnoverIncrease:   reportSummaryReq.AccountReceivableTurnover[1],
					AccountsPayableTurnover:             reportSummaryReq.AccountsPayableTurnover[0],
					AccountsPayableTurnoverIncrease:     reportSummaryReq.AccountsPayableTurnover[1],
					CurrentAssetTurnoverRate:            reportSummaryReq.CurrentAssetTurnoverRate[0],
					CurrentAssetTurnoverRateIncrease:    reportSummaryReq.CurrentAssetTurnoverRate[1],
					FixedAssetTurnoverRatio:             reportSummaryReq.FixedAssetTurnoverRatio[0],
					FixedAssetTurnoverRatioIncrease:     reportSummaryReq.FixedAssetTurnoverRatio[1],
				}
				db.Create(&reportSummary)
			}
		}
		time.Sleep(time.Duration(3) * time.Second)
	}

}

type IncomeReq struct {
	ReportDate                  uint       `json:"report_date"`
	ReportName                  string     `json:"report_name"`
	NetProfit                   [2]float32 `json:"net_profit"`
	NetProfitAtsopc             [2]float32 `json:"net_profit_atsopc"`
	TotalRevenue                [2]float32 `json:"total_revenue"`
	Op                          [2]float32 `json:"op"`
	IncomeFromChgInFv           [2]float32 `json:"income_from_chg_in_fv"`
	InvestIncomesFromRr         [2]float32 `json:"invest_incomes_from_rr"`
	InvestIncome                [2]float32 `json:"invest_income"`
	ExchgGain                   [2]float32 `json:"exchg_gain"`
	OperatingTaxesAndSurcharge  [2]float32 `json:"operating_taxes_and_surcharge"`
	AssetImpairmentLoss         [2]float32 `json:"asset_impairment_loss"`
	NonOperatingIncome          [2]float32 `json:"non_operating_income"`
	NonOperatingPayout          [2]float32 `json:"non_operating_payout"`
	ProfitTotalAmt              [2]float32 `json:"profit_total_amt"`
	MinorityGal                 [2]float32 `json:"minority_gal"`
	BasicEps                    [2]float32 `json:"basic_eps"`
	DltEarningsPerShare         [2]float32 `json:"dlt_earnings_per_share"`
	OthrCompreIncomeAtoopc      [2]float32 `json:"othr_compre_income_atoopc"`
	OthrCompreIncomeAtms        [2]float32 `json:"othr_compre_income_atms"`
	TotalCompreIncome           [2]float32 `json:"total_compre_income"`
	TotalCompreIncomeAtsopc     [2]float32 `json:"total_compre_income_atsopc"`
	TotalCompreIncomeAtms       [2]float32 `json:"total_compre_income_atms"`
	OthrCompreIncome            [2]float32 `json:"othr_compre_income"`
	NetProfitAfterNrgalAtsolc   [2]float32 `json:"net_profit_after_nrgal_atsolc"`
	IncomeTaxExpenses           [2]float32 `json:"income_tax_expenses"`
	CreditImpairmentLoss        [2]float32 `json:"credit_impairment_loss"`
	Revenue                     [2]float32 `json:"revenue"`
	OperatingCosts              [2]float32 `json:"operating_costs"`
	OperatingCost               [2]float32 `json:"operating_cost"`
	SalesFee                    [2]float32 `json:"sales_fee"`
	ManageFee                   [2]float32 `json:"manage_fee"`
	FinancingExpenses           [2]float32 `json:"financing_expenses"`
	RadCost                     [2]float32 `json:"rad_cost"`
	FinanceCostInterestFee      [2]float32 `json:"finance_cost_interest_fee"`
	FinanceCostInterestIncome   [2]float32 `json:"finance_cost_interest_income"`
	AssetDisposalIncome         [2]float32 `json:"asset_disposal_income"`
	OtherIncome                 [2]float32 `json:"other_income"`
	NoncurrentAssetDisposalGain [2]float32 `json:"noncurrent_assets_dispose_gain"`
	NoncurrentAssetDisposalLoss [2]float32 `json:"noncurrent_assets_dispose_loss"`
	NetProfitBi                 [2]float32 `json:"net_profit_bi"`
	ContinousOperatingNp        [2]float32 `json:"continous_operating_np"`
}
type IncomeDataReq struct {
	QuoteName      string
	LastReportRame string
	List           []IncomeReq
}

type IncomeReqstruct struct {
	Data             IncomeDataReq
	ErrorCode        uint
	ErrorDescription string
}

func getIncome() {
	var companies []models.Company
	var incomes []models.Income
	var companyCodes []string
	var reportCodes []string
	db.Table("companies").Select("code").Find(&companies)
	db.Distinct("company_code").Select("company_code").Find(&incomes)

	for _, company := range companies {
		companyCodes = append(companyCodes, company.Code)
	}

	for _, report := range incomes {
		reportCodes = append(reportCodes, report.CompanyCode)
	}

	codes := utils.Difference(companyCodes, reportCodes)
	fmt.Println("codes", codes)

	client := getClient()
	for _, code := range codes {
		const path string = "http://stock.xueqiu.com/v5/stock/finance/cn/income.json?type=ALL&is_detail=true&count=100"
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			log.Fatal(err)
		}
		q := req.URL.Query()
		q.Set("symbol", code)
		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal("Errored when sending request to the server")
		}

		fmt.Println(resp.Status)
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal(string(body))
		}

		var incomeReqstruct IncomeReqstruct

		json.NewDecoder(resp.Body).Decode(&incomeReqstruct)
		for _, incomeReq := range incomeReqstruct.Data.List {
			var income models.Income
			db.Where("company_code = ? AND report_name = ?", code, incomeReq.ReportName).First(&income)
			if income.ReportName == "" {
				income := models.Income{
					Category:                            fetchQ(incomeReq.ReportName),
					CompanyCode:                         code,
					ReportName:                          incomeReq.ReportName,
					ReportDate:                          incomeReq.ReportDate,
					NetProfit:                           incomeReq.NetProfit[0],
					NetProfitIncrease:                   incomeReq.NetProfit[1],
					NetProfitAtsopc:                     incomeReq.NetProfitAtsopc[0],
					NetProfitAtsopcIncrease:             incomeReq.NetProfitAtsopc[1],
					TotalRevenue:                        incomeReq.TotalRevenue[0],
					TotalRevenueIncrease:                incomeReq.TotalRevenue[1],
					Op:                                  incomeReq.Op[0],
					OpIncrease:                          incomeReq.Op[1],
					IncomeFromChgInFv:                   incomeReq.IncomeFromChgInFv[0],
					IncomeFromChgInFvIncrease:           incomeReq.IncomeFromChgInFv[1],
					InvestIncomesFromRr:                 incomeReq.InvestIncomesFromRr[0],
					InvestIncomesFromRrIncrease:         incomeReq.InvestIncomesFromRr[1],
					InvestIncome:                        incomeReq.InvestIncome[0],
					InvestIncomeIncrease:                incomeReq.InvestIncome[1],
					ExchgGain:                           incomeReq.ExchgGain[0],
					ExchgGainIncrease:                   incomeReq.ExchgGain[1],
					OperatingTaxesAndSurcharge:          incomeReq.OperatingTaxesAndSurcharge[0],
					OperatingTaxesAndSurchargeIncrease:  incomeReq.OperatingTaxesAndSurcharge[1],
					AssetImpairmentLoss:                 incomeReq.AssetImpairmentLoss[0],
					AssetImpairmentLossIncrease:         incomeReq.AssetImpairmentLoss[1],
					NonOperatingIncome:                  incomeReq.NonOperatingIncome[0],
					NonOperatingIncomeIncrease:          incomeReq.NonOperatingIncome[1],
					NonOperatingPayout:                  incomeReq.NonOperatingPayout[0],
					NonOperatingPayoutIncrease:          incomeReq.NonOperatingPayout[1],
					ProfitTotalAmt:                      incomeReq.ProfitTotalAmt[0],
					ProfitTotalAmtIncrease:              incomeReq.ProfitTotalAmt[1],
					MinorityGal:                         incomeReq.MinorityGal[0],
					MinorityGalIncrease:                 incomeReq.MinorityGal[1],
					BasicEps:                            incomeReq.BasicEps[0],
					BasicEpsIncrease:                    incomeReq.BasicEps[1],
					DltEarningsPerShare:                 incomeReq.DltEarningsPerShare[0],
					DltEarningsPerShareIncrease:         incomeReq.DltEarningsPerShare[1],
					OthrCompreIncomeAtoopc:              incomeReq.OthrCompreIncomeAtoopc[0],
					OthrCompreIncomeAtoopcIncrease:      incomeReq.OthrCompreIncomeAtoopc[1],
					OthrCompreIncomeAtms:                incomeReq.OthrCompreIncomeAtms[0],
					OthrCompreIncomeAtmsIncrease:        incomeReq.OthrCompreIncomeAtms[1],
					TotalCompreIncome:                   incomeReq.TotalCompreIncome[0],
					TotalCompreIncomeIncrease:           incomeReq.TotalCompreIncome[1],
					TotalCompreIncomeAtsopc:             incomeReq.TotalCompreIncomeAtsopc[0],
					TotalCompreIncomeAtsopcIncrease:     incomeReq.TotalCompreIncomeAtsopc[1],
					TotalCompreIncomeAtms:               incomeReq.TotalCompreIncomeAtms[0],
					TotalCompreIncomeAtmsIncrease:       incomeReq.TotalCompreIncomeAtms[1],
					OthrCompreIncome:                    incomeReq.OthrCompreIncome[0],
					OthrCompreIncomeIncrease:            incomeReq.OthrCompreIncome[1],
					NetProfitAfterNrgalAtsolc:           incomeReq.NetProfitAfterNrgalAtsolc[0],
					NetProfitAfterNrgalAtsolcIncrease:   incomeReq.NetProfitAfterNrgalAtsolc[1],
					IncomeTaxExpenses:                   incomeReq.IncomeTaxExpenses[0],
					IncomeTaxExpensesIncrease:           incomeReq.IncomeTaxExpenses[1],
					CreditImpairmentLoss:                incomeReq.CreditImpairmentLoss[0],
					CreditImpairmentLossIncrease:        incomeReq.CreditImpairmentLoss[1],
					Revenue:                             incomeReq.Revenue[0],
					RevenueIncrease:                     incomeReq.Revenue[1],
					OperatingCosts:                      incomeReq.OperatingCosts[0],
					OperatingCostsIncrease:              incomeReq.OperatingCosts[1],
					OperatingCost:                       incomeReq.OperatingCost[0],
					OperatingCostIncrease:               incomeReq.OperatingCost[1],
					SalesFee:                            incomeReq.SalesFee[0],
					SalesFeeIncrease:                    incomeReq.SalesFee[1],
					ManageFee:                           incomeReq.ManageFee[0],
					ManageFeeIncrease:                   incomeReq.ManageFee[1],
					FinancingExpenses:                   incomeReq.FinancingExpenses[0],
					FinancingExpensesIncrease:           incomeReq.FinancingExpenses[1],
					RadCost:                             incomeReq.RadCost[0],
					RadCostIncrease:                     incomeReq.RadCost[1],
					FinanceCostInterestFee:              incomeReq.FinanceCostInterestFee[0],
					FinanceCostInterestFeeIncrease:      incomeReq.FinanceCostInterestFee[1],
					FinanceCostInterestIncome:           incomeReq.FinanceCostInterestIncome[0],
					FinanceCostInterestIncomeIncrease:   incomeReq.FinanceCostInterestIncome[1],
					AssetDisposalIncome:                 incomeReq.AssetDisposalIncome[0],
					AssetDisposalIncomeIncrease:         incomeReq.AssetDisposalIncome[1],
					OtherIncome:                         incomeReq.OtherIncome[0],
					OtherIncomeIncrease:                 incomeReq.OtherIncome[1],
					NoncurrentAssetDisposalGain:         incomeReq.NoncurrentAssetDisposalGain[0],
					NoncurrentAssetDisposalGainIncrease: incomeReq.NoncurrentAssetDisposalGain[1],
					NoncurrentAssetDisposalLoss:         incomeReq.NoncurrentAssetDisposalLoss[0],
					NoncurrentAssetDisposalLossIncrease: incomeReq.NoncurrentAssetDisposalLoss[1],
					NetProfitBi:                         incomeReq.NetProfitBi[0],
					NetProfitBiIncrease:                 incomeReq.NetProfitBi[1],
					ContinousOperatingNp:                incomeReq.ContinousOperatingNp[0],
					ContinousOperatingNpIncrease:        incomeReq.ContinousOperatingNp[1],
				}
				db.Create(&income)
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

}

// parse xueqiu company's finace data
func main() {
	// getCompanies()
	// os.Setenv("HTTP_PROXY", "http://127.0.0.1:8008")
	// getReportSummary()
	getIncome()
}
