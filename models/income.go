package models

import "gorm.io/gorm"

// Income for company
type Income struct {
	gorm.Model
	Category                            string
	CompanyCode                         string
	ReportName                          string
	ReportDate                          uint
	NetProfit                           float32
	NetProfitIncrease                   float32
	NetProfitAtsopc                     float32
	NetProfitAtsopcIncrease             float32
	TotalRevenue                        float32
	TotalRevenueIncrease                float32
	Op                                  float32
	OpIncrease                          float32
	IncomeFromChgInFv                   float32
	IncomeFromChgInFvIncrease           float32
	InvestIncomesFromRr                 float32
	InvestIncomesFromRrIncrease         float32
	InvestIncome                        float32
	InvestIncomeIncrease                float32
	ExchgGain                           float32
	ExchgGainIncrease                   float32
	OperatingTaxesAndSurcharge          float32
	OperatingTaxesAndSurchargeIncrease  float32
	AssetImpairmentLoss                 float32
	AssetImpairmentLossIncrease         float32
	NonOperatingIncome                  float32
	NonOperatingIncomeIncrease          float32
	NonOperatingPayout                  float32
	NonOperatingPayoutIncrease          float32
	ProfitTotalAmt                      float32
	ProfitTotalAmtIncrease              float32
	MinorityGal                         float32
	MinorityGalIncrease                 float32
	BasicEps                            float32
	BasicEpsIncrease                    float32
	DltEarningsPerShare                 float32
	DltEarningsPerShareIncrease         float32
	OthrCompreIncomeAtoopc              float32
	OthrCompreIncomeAtoopcIncrease      float32
	OthrCompreIncomeAtms                float32
	OthrCompreIncomeAtmsIncrease        float32
	TotalCompreIncome                   float32
	TotalCompreIncomeIncrease           float32
	TotalCompreIncomeAtsopc             float32
	TotalCompreIncomeAtsopcIncrease     float32
	TotalCompreIncomeAtms               float32
	TotalCompreIncomeAtmsIncrease       float32
	OthrCompreIncome                    float32
	OthrCompreIncomeIncrease            float32
	NetProfitAfterNrgalAtsolc           float32
	NetProfitAfterNrgalAtsolcIncrease   float32
	IncomeTaxExpenses                   float32
	IncomeTaxExpensesIncrease           float32
	CreditImpairmentLoss                float32
	CreditImpairmentLossIncrease        float32
	Revenue                             float32
	RevenueIncrease                     float32
	OperatingCosts                      float32
	OperatingCostsIncrease              float32
	OperatingCost                       float32
	OperatingCostIncrease               float32
	SalesFee                            float32
	SalesFeeIncrease                    float32
	ManageFee                           float32
	ManageFeeIncrease                   float32
	FinancingExpenses                   float32
	FinancingExpensesIncrease           float32
	RadCost                             float32
	RadCostIncrease                     float32
	FinanceCostInterestFee              float32
	FinanceCostInterestFeeIncrease      float32
	FinanceCostInterestIncome           float32
	FinanceCostInterestIncomeIncrease   float32
	AssetDisposalIncome                 float32
	AssetDisposalIncomeIncrease         float32
	OtherIncome                         float32
	OtherIncomeIncrease                 float32
	NoncurrentAssetDisposalGain         float32
	NoncurrentAssetDisposalGainIncrease float32
	NoncurrentAssetDisposalLoss         float32
	NoncurrentAssetDisposalLossIncrease float32
	NetProfitBi                         float32
	NetProfitBiIncrease                 float32
	ContinousOperatingNp                float32
	ContinousOperatingNpIncrease        float32
}