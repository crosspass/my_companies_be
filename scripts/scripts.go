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
	"gorm.io/gorm"

	"github.com/my-companies-be/connect"
	"github.com/my-companies-be/models"
	"github.com/my-companies-be/utils"
)

// var dsn = "host=localhost user=wu password=gorm dbname=my_companies port=5432 sslmode=disable TimeZone=Asia/Shanghai"
var db *gorm.DB

func init() {
	db = connect.Db
}

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

	types := [1]string{"sh_sz"}

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
			time.Sleep(time.Duration(1) * time.Second)
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
	// filter updated records
	t := time.Now().AddDate(0, -1, 0)
	db.Where("created_at > ?", t).Distinct("company_code").Select("company_code").Find(&reportSummaries)

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
		const path string = "https://stock.xueqiu.com/v5/stock/finance/cn/indicator.json?type=ALL&is_detail=true&count=100"
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			log.Fatal(err)
		}
		q := req.URL.Query()
		q.Set("symbol", code)
		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Errored when sending request to the server")
		}
		defer resp.Body.Close()

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
		time.Sleep(time.Duration(1) * time.Second)
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
	t := time.Now().AddDate(0, -1, 0)
	db.Where("created_at > ?", t).Distinct("company_code").Select("company_code").Find(&incomes)

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
		const path string = "https://stock.xueqiu.com/v5/stock/finance/cn/income.json?type=ALL&is_detail=true&count=100"
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			log.Fatal(err)
		}
		q := req.URL.Query()
		q.Set("symbol", code)
		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Errored when sending request to the server")
		}
		defer resp.Body.Close()

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

type CashFlowReq struct {
	ReportDate                 uint       `json:"report_date"`
	ReportName                 string     `json:"report_name"`
	NcfFromOa                  [2]float32 `json:"ncf_from_oa"`
	NcfFromIa                  [2]float32 `json:"ncf_from_ia"`
	NcfFromFa                  [2]float32 `json:"ncf_from_fa"`
	CashReceivedOfOthrOa       [2]float32 `json:"cash_received_of_othr_oa"`
	SubTotalOfCiFromOa         [2]float32 `json:"sub_total_of_ci_from_oa"`
	CashPaidToEmployeeEtc      [2]float32 `json:"cash_paid_to_employee_etc"`
	PaymentsOfAllTaxes         [2]float32 `json:"payments_of_all_taxes"`
	OthrcashPaidRelatingToOa   [2]float32 `json:"othrcash_paid_relating_to_oa"`
	SubTotalOfCosFromOa        [2]float32 `json:"sub_total_of_cos_from_oa"`
	CashReceivedOfDspslInvest  [2]float32 `json:"cash_received_of_dspsl_invest"`
	InvestIncomeCashReceived   [2]float32 `json:"invest_income_cash_received"`
	NetCashOfDisposalAssets    [2]float32 `json:"net_cash_of_disposal_assets"`
	NetCashOfDisposalBranch    [2]float32 `json:"net_cash_of_disposal_branch"`
	CashReceivedOfOthrIa       [2]float32 `json:"cash_received_of_othr_ia"`
	SubTotalOfCiFromIa         [2]float32 `json:"sub_total_of_ci_from_ia"`
	InvestPaidCash             [2]float32 `json:"invest_paid_cash"`
	CashPaidForAssets          [2]float32 `json:"cash_paid_for_assets"`
	OthrcashPaidRelatingToIa   [2]float32 `json:"othrcash_paid_relating_to_ia"`
	SubTotalOfCosFromIa        [2]float32 `json:"sub_total_of_cos_from_ia"`
	CashReceivedOfAbsorbInvest [2]float32 `json:"cash_received_of_absorb_invest"`
	CashReceivedFromInvestor   [2]float32 `json:"cash_received_from_investor"`
	CashReceivedFromBondIssue  [2]float32 `json:"cash_received_from_bond_issue"`
	CashReceivedOfBorrowing    [2]float32 `json:"cash_received_of_borrowing"`
	CashReceivedOfOthrFa       [2]float32 `json:"cash_received_of_othr_fa"`
	SubTotalOfCiFromFa         [2]float32 `json:"sub_total_of_ci_from_fa"`
	CashPayForDebt             [2]float32 `json:"cash_pay_for_debt"`
	CashPaidOfDistribution     [2]float32 `json:"cash_paid_of_distribution"`
	BranchPaidToMinorityHolder [2]float32 `json:"branch_paid_to_minority_holder"`
	OthrcashPaidRelatingToFa   [2]float32 `json:"othrcash_paid_relating_to_fa"`
	SubTotalOfCosFromFa        [2]float32 `json:"sub_total_of_cos_from_fa"`
	EffectOfExchangeChgOnCce   [2]float32 `json:"effect_of_exchange_chg_on_cce"`
	NetIncreaseInCce           [2]float32 `json:"net_increase_in_cce"`
	InitialBalanceOfCce        [2]float32 `json:"initial_balance_of_cce"`
	FinalBalanceOfCce          [2]float32 `json:"final_balance_of_cce"`
	CashReceivedOfSalesService [2]float32 `json:"cash_received_of_sales_service"`
	RefundOfTaxAndLevies       [2]float32 `json:"refund_of_tax_and_levies"`
	GoodsBuyAndServiceCashPay  [2]float32 `json:"goods_buy_and_service_cash_pay"`
	NetCashAmtFromBranch       [2]float32 `json:"net_cash_amt_from_branch"`
}
type CashFlowDataReq struct {
	QuoteName      string
	LastReportRame string
	List           []CashFlowReq
}

type CashFlowReqstruct struct {
	Data             CashFlowDataReq
	ErrorCode        uint
	ErrorDescription string
}

func getCashFlow() {
	var companies []models.Company
	var cashFlow []models.CashFlow
	var companyCodes []string
	var reportCodes []string
	db.Table("companies").Select("code").Find(&companies)
	t := time.Now().AddDate(0, -1, 0)
	db.Where("created_at > ?", t).Distinct("company_code").Select("company_code").Find(&cashFlow)

	for _, company := range companies {
		companyCodes = append(companyCodes, company.Code)
	}

	for _, report := range cashFlow {
		reportCodes = append(reportCodes, report.CompanyCode)
	}

	codes := utils.Difference(companyCodes, reportCodes)
	fmt.Println("codes", codes)

	client := getClient()
	for _, code := range codes {
		const path string = "https://stock.xueqiu.com/v5/stock/finance/cn/cash_flow.json?type=ALL&is_detail=true&count=100"
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			log.Fatal(err)
		}
		q := req.URL.Query()
		q.Set("symbol", code)
		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Errored when sending request to the server")
		}
		defer resp.Body.Close()

		fmt.Println(resp.Status)
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal(string(body))
		}

		var cashFlowReqstruct CashFlowReqstruct

		json.NewDecoder(resp.Body).Decode(&cashFlowReqstruct)
		for _, cashFlowReq := range cashFlowReqstruct.Data.List {
			var cashFlow models.CashFlow
			db.Where("company_code = ? AND report_name = ?", code, cashFlowReq.ReportName).First(&cashFlow)
			if cashFlow.ReportName == "" {
				cashFlow := models.CashFlow{
					Category:                           fetchQ(cashFlowReq.ReportName),
					CompanyCode:                        code,
					ReportName:                         cashFlowReq.ReportName,
					ReportDate:                         cashFlowReq.ReportDate,
					NcfFromOa:                          cashFlowReq.NcfFromOa[0],
					NcfFromOaIncrease:                  cashFlowReq.NcfFromOa[1],
					NcfFromIa:                          cashFlowReq.NcfFromIa[0],
					NcfFromIaIncrease:                  cashFlowReq.NcfFromIa[1],
					NcfFromFa:                          cashFlowReq.NcfFromFa[0],
					NcfFromFaIncrease:                  cashFlowReq.NcfFromFa[1],
					CashReceivedOfOthrOa:               cashFlowReq.CashReceivedOfOthrOa[0],
					CashReceivedOfOthrOaIncrease:       cashFlowReq.CashReceivedOfOthrOa[1],
					SubTotalOfCiFromOa:                 cashFlowReq.SubTotalOfCiFromOa[0],
					SubTotalOfCiFromOaIncrease:         cashFlowReq.SubTotalOfCiFromOa[1],
					CashPaidToEmployeeEtc:              cashFlowReq.CashPaidToEmployeeEtc[0],
					CashPaidToEmployeeEtcIncrease:      cashFlowReq.CashPaidToEmployeeEtc[1],
					PaymentsOfAllTaxes:                 cashFlowReq.PaymentsOfAllTaxes[0],
					PaymentsOfAllTaxesIncrease:         cashFlowReq.PaymentsOfAllTaxes[1],
					OthrcashPaidRelatingToOa:           cashFlowReq.OthrcashPaidRelatingToOa[0],
					OthrcashPaidRelatingToOaIncrease:   cashFlowReq.OthrcashPaidRelatingToOa[1],
					SubTotalOfCosFromOa:                cashFlowReq.SubTotalOfCosFromOa[0],
					SubTotalOfCosFromOaIncrease:        cashFlowReq.SubTotalOfCosFromOa[1],
					CashReceivedOfDspslInvest:          cashFlowReq.CashReceivedOfDspslInvest[0],
					CashReceivedOfDspslInvestIncrease:  cashFlowReq.CashReceivedOfDspslInvest[1],
					InvestIncomeCashReceived:           cashFlowReq.InvestIncomeCashReceived[0],
					InvestIncomeCashReceivedIncrease:   cashFlowReq.InvestIncomeCashReceived[1],
					NetCashOfDisposalAssets:            cashFlowReq.NetCashOfDisposalAssets[0],
					NetCashOfDisposalAssetsIncrease:    cashFlowReq.NetCashOfDisposalAssets[1],
					NetCashOfDisposalBranch:            cashFlowReq.NetCashOfDisposalBranch[0],
					NetCashOfDisposalBranchIncrease:    cashFlowReq.NetCashOfDisposalBranch[1],
					CashReceivedOfOthrIa:               cashFlowReq.CashReceivedOfOthrIa[0],
					CashReceivedOfOthrIaIncrease:       cashFlowReq.CashReceivedOfOthrIa[1],
					SubTotalOfCiFromIa:                 cashFlowReq.SubTotalOfCiFromIa[0],
					SubTotalOfCiFromIaIncrease:         cashFlowReq.SubTotalOfCiFromIa[1],
					InvestPaidCash:                     cashFlowReq.InvestPaidCash[0],
					InvestPaidCashIncrease:             cashFlowReq.InvestPaidCash[1],
					CashPaidForAssets:                  cashFlowReq.CashPaidForAssets[0],
					CashPaidForAssetsIncrease:          cashFlowReq.CashPaidForAssets[1],
					OthrcashPaidRelatingToIa:           cashFlowReq.OthrcashPaidRelatingToIa[0],
					OthrcashPaidRelatingToIaIncrease:   cashFlowReq.OthrcashPaidRelatingToIa[1],
					SubTotalOfCosFromIa:                cashFlowReq.SubTotalOfCosFromIa[0],
					SubTotalOfCosFromIaIncrease:        cashFlowReq.SubTotalOfCosFromIa[1],
					CashReceivedOfAbsorbInvest:         cashFlowReq.CashReceivedOfAbsorbInvest[0],
					CashReceivedOfAbsorbInvestIncrease: cashFlowReq.CashReceivedOfAbsorbInvest[1],
					CashReceivedFromInvestor:           cashFlowReq.CashReceivedFromInvestor[0],
					CashReceivedFromInvestorIncrease:   cashFlowReq.CashReceivedFromInvestor[1],
					CashReceivedFromBondIssue:          cashFlowReq.CashReceivedFromBondIssue[0],
					CashReceivedFromBondIssueIncrease:  cashFlowReq.CashReceivedFromBondIssue[1],
					CashReceivedOfBorrowing:            cashFlowReq.CashReceivedOfBorrowing[0],
					CashReceivedOfBorrowingIncrease:    cashFlowReq.CashReceivedOfBorrowing[1],
					CashReceivedOfOthrFa:               cashFlowReq.CashReceivedOfOthrFa[0],
					CashReceivedOfOthrFaIncrease:       cashFlowReq.CashReceivedOfOthrFa[1],
					SubTotalOfCiFromFa:                 cashFlowReq.SubTotalOfCiFromFa[0],
					SubTotalOfCiFromFaIncrease:         cashFlowReq.SubTotalOfCiFromFa[1],
					CashPayForDebt:                     cashFlowReq.CashPayForDebt[0],
					CashPayForDebtIncrease:             cashFlowReq.CashPayForDebt[1],
					CashPaidOfDistribution:             cashFlowReq.CashPaidOfDistribution[0],
					CashPaidOfDistributionIncrease:     cashFlowReq.CashPaidOfDistribution[1],
					BranchPaidToMinorityHolder:         cashFlowReq.BranchPaidToMinorityHolder[0],
					BranchPaidToMinorityHolderIncrease: cashFlowReq.BranchPaidToMinorityHolder[1],
					OthrcashPaidRelatingToFa:           cashFlowReq.OthrcashPaidRelatingToFa[0],
					OthrcashPaidRelatingToFaIncrease:   cashFlowReq.OthrcashPaidRelatingToFa[1],
					SubTotalOfCosFromFa:                cashFlowReq.SubTotalOfCosFromFa[0],
					SubTotalOfCosFromFaIncrease:        cashFlowReq.SubTotalOfCosFromFa[1],
					EffectOfExchangeChgOnCce:           cashFlowReq.EffectOfExchangeChgOnCce[0],
					EffectOfExchangeChgOnCceIncrease:   cashFlowReq.EffectOfExchangeChgOnCce[1],
					NetIncreaseInCce:                   cashFlowReq.NetIncreaseInCce[0],
					NetIncreaseInCceIncrease:           cashFlowReq.NetIncreaseInCce[1],
					InitialBalanceOfCce:                cashFlowReq.InitialBalanceOfCce[0],
					InitialBalanceOfCceIncrease:        cashFlowReq.InitialBalanceOfCce[1],
					FinalBalanceOfCce:                  cashFlowReq.FinalBalanceOfCce[0],
					FinalBalanceOfCceIncrease:          cashFlowReq.FinalBalanceOfCce[1],
					CashReceivedOfSalesService:         cashFlowReq.CashReceivedOfSalesService[0],
					CashReceivedOfSalesServiceIncrease: cashFlowReq.CashReceivedOfSalesService[1],
					RefundOfTaxAndLevies:               cashFlowReq.RefundOfTaxAndLevies[0],
					RefundOfTaxAndLeviesIncrease:       cashFlowReq.RefundOfTaxAndLevies[1],
					GoodsBuyAndServiceCashPay:          cashFlowReq.GoodsBuyAndServiceCashPay[0],
					GoodsBuyAndServiceCashPayIncrease:  cashFlowReq.GoodsBuyAndServiceCashPay[1],
					NetCashAmtFromBranch:               cashFlowReq.NetCashAmtFromBranch[0],
					NetCashAmtFromBranchIncrease:       cashFlowReq.NetCashAmtFromBranch[1],
				}
				db.Create(&cashFlow)
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

type BalanceReq struct {
	ReportDate                 uint       `json:"report_date"`
	ReportName                 string     `json:"report_name"`
	TotalAssets                [2]float32 `json:"total_assets"`
	TotalLiab                  [2]float32 `json:"total_liab"`
	AssetLiabRatio             [2]float32 `json:"asset_liab_ratio"`
	TotalQuityAtsopc           [2]float32 `json:"total_quity_atsopc"`
	TradableFnnclAssets        [2]float32 `json:"tradable_fnncl_assets"`
	InterestReceivable         [2]float32 `json:"interest_receivable"`
	SaleableFinacialAssets     [2]float32 `json:"saleable_finacial_assets"`
	HeldToMaturityInvest       [2]float32 `json:"held_to_maturity_invest"`
	FixedAsset                 [2]float32 `json:"fixed_asset"`
	IntangibleAssets           [2]float32 `json:"intangible_assets"`
	ConstructionInProcess      [2]float32 `json:"construction_in_process"`
	DtAssets                   [2]float32 `json:"dt_assets"`
	TradableFnnclLiab          [2]float32 `json:"tradable_fnncl_liab"`
	PayrollPayable             [2]float32 `json:"payroll_payable"`
	TaxPayable                 [2]float32 `json:"tax_payable"`
	EstimatedLiab              [2]float32 `json:"estimated_liab"`
	DtLiab                     [2]float32 `json:"dt_liab"`
	BondPayable                [2]float32 `json:"bond_payable"`
	Shares                     [2]float32 `json:"shares"`
	CapitalReserve             [2]float32 `json:"capital_reserve"`
	EarnedSurplus              [2]float32 `json:"earned_surplus"`
	UndstrbtdProfit            [2]float32 `json:"undstrbtd_profit"`
	MinorityEquity             [2]float32 `json:"minority_equity"`
	TotalHoldersEquity         [2]float32 `json:"total_holders_equity"`
	TotalLiabAndHoldersEquity  [2]float32 `json:"total_liab_and_holders_equity"`
	LtEquityInvest             [2]float32 `json:"lt_equity_invest"`
	DerivativeFnnclLiab        [2]float32 `json:"derivative_fnncl_liab"`
	GeneralRiskProvision       [2]float32 `json:"general_risk_provision"`
	FrgnCurrencyConvertDiff    [2]float32 `json:"frgn_currency_convert_diff"`
	Goodwill                   [2]float32 `json:"goodwill"`
	InvestProperty             [2]float32 `json:"invest_property"`
	InterestPayable            [2]float32 `json:"interest_payable"`
	TreasuryStock              [2]float32 `json:"treasury_stock"`
	OthrCompreIncome           [2]float32 `json:"othr_compre_income"`
	OthrEquityInstruments      [2]float32 `json:"othr_equity_instruments"`
	CurrencyFunds              [2]float32 `json:"currency_funds"`
	BillsReceivable            [2]float32 `json:"bills_receivable"`
	AccountReceivable          [2]float32 `json:"account_receivable"`
	PrePayment                 [2]float32 `json:"pre_payment"`
	DividendReceivable         [2]float32 `json:"dividend_receivable"`
	OthrReceivables            [2]float32 `json:"othr_receivables"`
	Inventory                  [2]float32 `json:"inventory"`
	NcaDueWithinOneYear        [2]float32 `json:"nca_due_within_one_year"`
	OthrCurrentAssets          [2]float32 `json:"othr_current_assets"`
	CurrentAssetsSi            [2]float32 `json:"current_assets_si"`
	TotalCurrentAssets         [2]float32 `json:"total_current_assets"`
	LtReceivable               [2]float32 `json:"lt_receivable"`
	DevExpenditure             [2]float32 `json:"dev_expenditure"`
	LtDeferredExpense          [2]float32 `json:"lt_deferred_expense"`
	OthrNoncurrentAssets       [2]float32 `json:"othr_noncurrent_assets"`
	NoncurrentAssetsSi         [2]float32 `json:"noncurrent_assets_si"`
	TotalNoncurrentAssets      [2]float32 `json:"total_noncurrent_assets"`
	StLoan                     [2]float32 `json:"st_loan"`
	BillPayable                [2]float32 `json:"bill_payable"`
	AccountsPayable            [2]float32 `json:"accounts_payable"`
	PreReceivable              [2]float32 `json:"pre_receivable"`
	DividendPayable            [2]float32 `json:"dividend_payable"`
	OthrPayables               [2]float32 `json:"othr_payables"`
	NoncurrentLiabDueIn1y      [2]float32 `json:"noncurrent_liab_due_in1y"`
	CurrentLiabSi              [2]float32 `json:"current_liab_si"`
	TotalCurrentLiab           [2]float32 `json:"total_current_liab"`
	LtLoan                     [2]float32 `json:"lt_loan"`
	LtPayable                  [2]float32 `json:"lt_payable"`
	SpecialPayable             [2]float32 `json:"special_payable"`
	OthrNonCurrentLiab         [2]float32 `json:"othr_non_current_liab"`
	NoncurrentLiabSi           [2]float32 `json:"noncurrent_liab_si"`
	TotalNoncurrentLiab        [2]float32 `json:"total_noncurrent_liab"`
	SalableFinancialAssets     [2]float32 `json:"salable_financial_assets"`
	OthrCurrentLiab            [2]float32 `json:"othr_current_liab"`
	ArAndBr                    [2]float32 `json:"ar_and_br"`
	ContractualAssets          [2]float32 `json:"contractual_assets"`
	BpAndAp                    [2]float32 `json:"bp_and_ap"`
	ContractLiabilities        [2]float32 `json:"contract_liabilities"`
	ToSaleAsset                [2]float32 `json:"to_sale_asset"`
	OtherEqInsInvest           [2]float32 `json:"other_eq_ins_invest"`
	OtherIlliquidFnnclAssets   [2]float32 `json:"other_illiquid_fnncl_assets"`
	FixedAssetSum              [2]float32 `json:"fixed_asset_sum"`
	FixedAssetsDisposal        [2]float32 `json:"fixed_assets_disposal"`
	ConstructionInProcessSum   [2]float32 `json:"construction_in_process_sum"`
	ProjectGoodsAndMaterial    [2]float32 `json:"project_goods_and_material"`
	ProductiveBiologicalAssets [2]float32 `json:"productive_biological_assets"`
	OilAndGasAsset             [2]float32 `json:"oil_and_gas_asset"`
	ToSaleDebt                 [2]float32 `json:"to_sale_debt"`
	LtPayableSum               [2]float32 `json:"lt_payable_sum"`
	NoncurrentLiabDi           [2]float32 `json:"noncurrent_liab_di"`
	PerpetualBond              [2]float32 `json:"perpetual_bond"`
	SpecialReserve             [2]float32 `json:"special_reserve"`
}
type BalanceDataReq struct {
	QuoteName      string
	LastReportRame string
	List           []BalanceReq
}

type BalanceReqstruct struct {
	Data             BalanceDataReq
	ErrorCode        uint
	ErrorDescription string
}

func getBalance() {
	var companies []models.Company
	var balances []models.Balance
	var companyCodes []string
	var reportCodes []string
	db.Table("companies").Select("code").Find(&companies)
	t := time.Now().AddDate(0, -1, 0)
	db.Where("created_at > ?", t).Distinct("company_code").Select("company_code").Find(&balances)

	for _, company := range companies {
		companyCodes = append(companyCodes, company.Code)
	}

	for _, balance := range balances {
		reportCodes = append(reportCodes, balance.CompanyCode)
	}

	codes := utils.Difference(companyCodes, reportCodes)
	fmt.Println("codes", codes)

	client := getClient()
	for _, code := range codes {
		const path string = "https://stock.xueqiu.com/v5/stock/finance/cn/balance.json?type=ALL&is_detail=true&count=100"
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			log.Fatal(err)
		}
		q := req.URL.Query()
		q.Set("symbol", code)
		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Errored when sending request to the server")
		}
		defer resp.Body.Close()

		fmt.Println(resp.Status)
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal(string(body))
		}

		var balanceReqstruct BalanceReqstruct

		json.NewDecoder(resp.Body).Decode(&balanceReqstruct)
		for _, balanceReq := range balanceReqstruct.Data.List {
			var balance models.Balance
			db.Where("company_code = ? AND report_name = ?", code, balanceReq.ReportName).First(&balance)
			if balance.ReportName == "" {
				balance := models.Balance{
					Category:                           fetchQ(balanceReq.ReportName),
					CompanyCode:                        code,
					ReportName:                         balanceReq.ReportName,
					ReportDate:                         balanceReq.ReportDate,
					TotalAssets:                        balanceReq.TotalAssets[0],
					TotalAssetsIncrease:                balanceReq.TotalAssets[1],
					TotalLiab:                          balanceReq.TotalLiab[0],
					TotalLiabIncrease:                  balanceReq.TotalLiab[1],
					AssetLiabRatio:                     balanceReq.AssetLiabRatio[0],
					AssetLiabRatioIncrease:             balanceReq.AssetLiabRatio[1],
					TotalQuityAtsopc:                   balanceReq.TotalQuityAtsopc[0],
					TotalQuityAtsopcIncrease:           balanceReq.TotalQuityAtsopc[1],
					TradableFnnclAssets:                balanceReq.TradableFnnclAssets[0],
					TradableFnnclAssetsIncrease:        balanceReq.TradableFnnclAssets[1],
					InterestReceivable:                 balanceReq.InterestReceivable[0],
					InterestReceivableIncrease:         balanceReq.InterestReceivable[1],
					SaleableFinacialAssets:             balanceReq.SaleableFinacialAssets[0],
					SaleableFinacialAssetsIncrease:     balanceReq.SaleableFinacialAssets[1],
					HeldToMaturityInvest:               balanceReq.HeldToMaturityInvest[0],
					HeldToMaturityInvestIncrease:       balanceReq.HeldToMaturityInvest[1],
					FixedAsset:                         balanceReq.FixedAsset[0],
					FixedAssetIncrease:                 balanceReq.FixedAsset[1],
					IntangibleAssets:                   balanceReq.IntangibleAssets[0],
					IntangibleAssetsIncrease:           balanceReq.IntangibleAssets[1],
					ConstructionInProcess:              balanceReq.ConstructionInProcess[0],
					ConstructionInProcessIncrease:      balanceReq.ConstructionInProcess[1],
					DtAssets:                           balanceReq.DtAssets[0],
					DtAssetsIncrease:                   balanceReq.DtAssets[1],
					TradableFnnclLiab:                  balanceReq.TradableFnnclLiab[0],
					TradableFnnclLiabIncrease:          balanceReq.TradableFnnclLiab[1],
					PayrollPayable:                     balanceReq.PayrollPayable[0],
					PayrollPayableIncrease:             balanceReq.PayrollPayable[1],
					TaxPayable:                         balanceReq.TaxPayable[0],
					TaxPayableIncrease:                 balanceReq.TaxPayable[1],
					EstimatedLiab:                      balanceReq.EstimatedLiab[0],
					EstimatedLiabIncrease:              balanceReq.EstimatedLiab[1],
					DtLiab:                             balanceReq.DtLiab[0],
					DtLiabIncrease:                     balanceReq.DtLiab[1],
					BondPayable:                        balanceReq.BondPayable[0],
					BondPayableIncrease:                balanceReq.BondPayable[1],
					Shares:                             balanceReq.Shares[0],
					SharesIncrease:                     balanceReq.Shares[1],
					CapitalReserve:                     balanceReq.CapitalReserve[0],
					CapitalReserveIncrease:             balanceReq.CapitalReserve[1],
					EarnedSurplus:                      balanceReq.EarnedSurplus[0],
					EarnedSurplusIncrease:              balanceReq.EarnedSurplus[1],
					UndstrbtdProfit:                    balanceReq.UndstrbtdProfit[0],
					UndstrbtdProfitIncrease:            balanceReq.UndstrbtdProfit[1],
					MinorityEquity:                     balanceReq.MinorityEquity[0],
					MinorityEquityIncrease:             balanceReq.MinorityEquity[1],
					TotalHoldersEquity:                 balanceReq.TotalHoldersEquity[0],
					TotalHoldersEquityIncrease:         balanceReq.TotalHoldersEquity[1],
					TotalLiabAndHoldersEquity:          balanceReq.TotalLiabAndHoldersEquity[0],
					TotalLiabAndHoldersEquityIncrease:  balanceReq.TotalLiabAndHoldersEquity[1],
					LtEquityInvest:                     balanceReq.LtEquityInvest[0],
					LtEquityInvestIncrease:             balanceReq.LtEquityInvest[1],
					DerivativeFnnclLiab:                balanceReq.DerivativeFnnclLiab[0],
					DerivativeFnnclLiabIncrease:        balanceReq.DerivativeFnnclLiab[1],
					GeneralRiskProvision:               balanceReq.GeneralRiskProvision[0],
					GeneralRiskProvisionIncrease:       balanceReq.GeneralRiskProvision[1],
					FrgnCurrencyConvertDiff:            balanceReq.FrgnCurrencyConvertDiff[0],
					FrgnCurrencyConvertDiffIncrease:    balanceReq.FrgnCurrencyConvertDiff[1],
					Goodwill:                           balanceReq.Goodwill[0],
					GoodwillIncrease:                   balanceReq.Goodwill[1],
					InvestProperty:                     balanceReq.InvestProperty[0],
					InvestPropertyIncrease:             balanceReq.InvestProperty[1],
					InterestPayable:                    balanceReq.InterestPayable[0],
					InterestPayableIncrease:            balanceReq.InterestPayable[1],
					TreasuryStock:                      balanceReq.TreasuryStock[0],
					TreasuryStockIncrease:              balanceReq.TreasuryStock[1],
					OthrCompreIncome:                   balanceReq.OthrCompreIncome[0],
					OthrCompreIncomeIncrease:           balanceReq.OthrCompreIncome[1],
					OthrEquityInstruments:              balanceReq.OthrEquityInstruments[0],
					OthrEquityInstrumentsIncrease:      balanceReq.OthrEquityInstruments[1],
					CurrencyFunds:                      balanceReq.CurrencyFunds[0],
					CurrencyFundsIncrease:              balanceReq.CurrencyFunds[1],
					BillsReceivable:                    balanceReq.BillsReceivable[0],
					BillsReceivableIncrease:            balanceReq.BillsReceivable[1],
					AccountReceivable:                  balanceReq.AccountReceivable[0],
					AccountReceivableIncrease:          balanceReq.AccountReceivable[1],
					PrePayment:                         balanceReq.PrePayment[0],
					PrePaymentIncrease:                 balanceReq.PrePayment[1],
					DividendReceivable:                 balanceReq.DividendReceivable[0],
					DividendReceivableIncrease:         balanceReq.DividendReceivable[1],
					OthrReceivables:                    balanceReq.OthrReceivables[0],
					OthrReceivablesIncrease:            balanceReq.OthrReceivables[1],
					Inventory:                          balanceReq.Inventory[0],
					InventoryIncrease:                  balanceReq.Inventory[1],
					NcaDueWithinOneYear:                balanceReq.NcaDueWithinOneYear[0],
					NcaDueWithinOneYearIncrease:        balanceReq.NcaDueWithinOneYear[1],
					OthrCurrentAssets:                  balanceReq.OthrCurrentAssets[0],
					OthrCurrentAssetsIncrease:          balanceReq.OthrCurrentAssets[1],
					CurrentAssetsSi:                    balanceReq.CurrentAssetsSi[0],
					CurrentAssetsSiIncrease:            balanceReq.CurrentAssetsSi[1],
					TotalCurrentAssets:                 balanceReq.TotalCurrentAssets[0],
					TotalCurrentAssetsIncrease:         balanceReq.TotalCurrentAssets[1],
					LtReceivable:                       balanceReq.LtReceivable[0],
					LtReceivableIncrease:               balanceReq.LtReceivable[1],
					DevExpenditure:                     balanceReq.DevExpenditure[0],
					DevExpenditureIncrease:             balanceReq.DevExpenditure[1],
					LtDeferredExpense:                  balanceReq.LtDeferredExpense[0],
					LtDeferredExpenseIncrease:          balanceReq.LtDeferredExpense[1],
					OthrNoncurrentAssets:               balanceReq.OthrNoncurrentAssets[0],
					OthrNoncurrentAssetsIncrease:       balanceReq.OthrNoncurrentAssets[1],
					NoncurrentAssetsSi:                 balanceReq.NoncurrentAssetsSi[0],
					NoncurrentAssetsSiIncrease:         balanceReq.NoncurrentAssetsSi[1],
					TotalNoncurrentAssets:              balanceReq.TotalNoncurrentAssets[0],
					TotalNoncurrentAssetsIncrease:      balanceReq.TotalNoncurrentAssets[1],
					StLoan:                             balanceReq.StLoan[0],
					StLoanIncrease:                     balanceReq.StLoan[1],
					BillPayable:                        balanceReq.BillPayable[0],
					BillPayableIncrease:                balanceReq.BillPayable[1],
					AccountsPayable:                    balanceReq.AccountsPayable[0],
					AccountsPayableIncrease:            balanceReq.AccountsPayable[1],
					PreReceivable:                      balanceReq.PreReceivable[0],
					PreReceivableIncrease:              balanceReq.PreReceivable[1],
					DividendPayable:                    balanceReq.DividendPayable[0],
					DividendPayableIncrease:            balanceReq.DividendPayable[1],
					OthrPayables:                       balanceReq.OthrPayables[0],
					OthrPayablesIncrease:               balanceReq.OthrPayables[1],
					NoncurrentLiabDueIn1y:              balanceReq.NoncurrentLiabDueIn1y[0],
					NoncurrentLiabDueIn1yIncrease:      balanceReq.NoncurrentLiabDueIn1y[1],
					CurrentLiabSi:                      balanceReq.CurrentLiabSi[0],
					CurrentLiabSiIncrease:              balanceReq.CurrentLiabSi[1],
					TotalCurrentLiab:                   balanceReq.TotalCurrentLiab[0],
					TotalCurrentLiabIncrease:           balanceReq.TotalCurrentLiab[1],
					LtLoan:                             balanceReq.LtLoan[0],
					LtLoanIncrease:                     balanceReq.LtLoan[1],
					LtPayable:                          balanceReq.LtPayable[0],
					LtPayableIncrease:                  balanceReq.LtPayable[1],
					SpecialPayable:                     balanceReq.SpecialPayable[0],
					SpecialPayableIncrease:             balanceReq.SpecialPayable[1],
					OthrNonCurrentLiab:                 balanceReq.OthrNonCurrentLiab[0],
					OthrNonCurrentLiabIncrease:         balanceReq.OthrNonCurrentLiab[1],
					NoncurrentLiabSi:                   balanceReq.NoncurrentLiabSi[0],
					NoncurrentLiabSiIncrease:           balanceReq.NoncurrentLiabSi[1],
					TotalNoncurrentLiab:                balanceReq.TotalNoncurrentLiab[0],
					TotalNoncurrentLiabIncrease:        balanceReq.TotalNoncurrentLiab[1],
					SalableFinancialAssets:             balanceReq.SalableFinancialAssets[0],
					SalableFinancialAssetsIncrease:     balanceReq.SalableFinancialAssets[1],
					OthrCurrentLiab:                    balanceReq.OthrCurrentLiab[0],
					OthrCurrentLiabIncrease:            balanceReq.OthrCurrentLiab[1],
					ArAndBr:                            balanceReq.ArAndBr[0],
					ArAndBrIncrease:                    balanceReq.ArAndBr[1],
					ContractualAssets:                  balanceReq.ContractualAssets[0],
					ContractualAssetsIncrease:          balanceReq.ContractualAssets[1],
					BpAndAp:                            balanceReq.BpAndAp[0],
					BpAndApIncrease:                    balanceReq.BpAndAp[1],
					ContractLiabilities:                balanceReq.ContractLiabilities[0],
					ContractLiabilitiesIncrease:        balanceReq.ContractLiabilities[1],
					ToSaleAsset:                        balanceReq.ToSaleAsset[0],
					ToSaleAssetIncrease:                balanceReq.ToSaleAsset[1],
					OtherEqInsInvest:                   balanceReq.OtherEqInsInvest[0],
					OtherEqInsInvestIncrease:           balanceReq.OtherEqInsInvest[1],
					OtherIlliquidFnnclAssets:           balanceReq.OtherIlliquidFnnclAssets[0],
					OtherIlliquidFnnclAssetsIncrease:   balanceReq.OtherIlliquidFnnclAssets[1],
					FixedAssetSum:                      balanceReq.FixedAssetSum[0],
					FixedAssetSumIncrease:              balanceReq.FixedAssetSum[1],
					FixedAssetsDisposal:                balanceReq.FixedAssetsDisposal[0],
					FixedAssetsDisposalIncrease:        balanceReq.FixedAssetsDisposal[1],
					ConstructionInProcessSum:           balanceReq.ConstructionInProcessSum[0],
					ConstructionInProcessSumIncrease:   balanceReq.ConstructionInProcessSum[1],
					ProjectGoodsAndMaterial:            balanceReq.ProjectGoodsAndMaterial[0],
					ProjectGoodsAndMaterialIncrease:    balanceReq.ProjectGoodsAndMaterial[1],
					ProductiveBiologicalAssets:         balanceReq.ProductiveBiologicalAssets[0],
					ProductiveBiologicalAssetsIncrease: balanceReq.ProductiveBiologicalAssets[1],
					OilAndGasAsset:                     balanceReq.OilAndGasAsset[0],
					OilAndGasAssetIncrease:             balanceReq.OilAndGasAsset[1],
					ToSaleDebt:                         balanceReq.ToSaleDebt[0],
					ToSaleDebtIncrease:                 balanceReq.ToSaleDebt[1],
					LtPayableSum:                       balanceReq.LtPayableSum[0],
					LtPayableSumIncrease:               balanceReq.LtPayableSum[1],
					NoncurrentLiabDi:                   balanceReq.NoncurrentLiabDi[0],
					NoncurrentLiabDiIncrease:           balanceReq.NoncurrentLiabDi[1],
					PerpetualBond:                      balanceReq.PerpetualBond[0],
					PerpetualBondIncrease:              balanceReq.PerpetualBond[1],
					SpecialReserve:                     balanceReq.SpecialReserve[0],
					SpecialReserveIncrease:             balanceReq.SpecialReserve[1],
				}
				db.Create(&balance)
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// parse xueqiu company's finace data
func main() {
	getCompanies()
	// os.Setenv("HTTP_PROXY", "http://127.0.0.1:8008")
	getReportSummary()
	getIncome()
	getCashFlow()
	getBalance()
}
