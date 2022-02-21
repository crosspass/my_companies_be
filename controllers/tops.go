package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// select company_code, avg(avg_roe) as sum_avg_roe, count(avg_roe) as scount from report_summaries where category = 'Q4' and avg_roe_increase > 0 group by company_code having count(avg_roe) > 10 order by sum_avg_roe desc limit 100;
// select company_code, avg(avg_roe) as avg_roe, count(avg_roe) as scount from report_summaries where category = 'Q4' and report_date >= 1514649600000 group by company_code having min(avg_roe) > 0 order by avg_roe desc limit 100;
// select company_code, avg(net_interest_of_total_assets), count(*) as count from report_summaries where company_code in (select company_code from report_summaries where category = 'Q4' and report_date >= 1514649600000 group by company_code having min(net_interest_of_total_assets) > 0 order by avg(net_interest_of_total_assets) desc limit 100) group by company_code order by count desc;
// select code, name from companies where code in (select company_code from report_summaries where category = 'Q4' and report_date >= 1514649600000 group by company_code having min(net_interest_of_total_assets) > 0 order by avg(net_interest_of_total_assets) desc limit 100);

// select code, name from companies where code in (select company_code from report_summaries where report_date >= 1514649600000 group by company_code having min(net_interest_of_total_assets) > 0 order by avg(net_interest_of_total_assets) desc limit 100);

type roaResponse struct {
	Code                        string
	Name                        string
	SumNetInterestOfTotalAssets float64
}

// Roa for roa top 100
func Roa(ctx *gin.Context) {
	sql := `
		select companies.code, companies.name, summary.sum_net_interest_of_total_assets
		from companies, (
		 select company_code, sum(net_interest_of_total_assets) as sum_net_interest_of_total_assets
		 from report_summaries where category = 'Q4' and report_date >= ?
		 group by company_code order by sum(net_interest_of_total_assets) desc limit 100) as summary
		where summary.company_code = companies.code order by summary.sum_net_interest_of_total_assets desc;
	`
	var roas = make([]roaResponse, 0)
	years, _ := strconv.Atoi(ctx.DefaultQuery("years", "2"))
	now := time.Now()
	currentYear, _, _ := now.Date()
	currentLocation := now.Location()
	year := time.Date(currentYear-years, 12, 31, 0, 0, 0, 0, currentLocation)
	db.Raw(sql, year.Unix()*1000).Scan(&roas)
	ctx.JSON(http.StatusOK, gin.H{
		"roas":    roas,
		"message": "ok",
	})
}

// roaIncreaseRespone
type roaIncreaseResponse struct {
	Code                                string
	Name                                string
	SumNetInterestOfTotalAssetsIncrease float64
}

// RoaIncrease for roaincrease top 100
func RoaIncrease(ctx *gin.Context) {
	sql := `
		select companies.code, companies.name, summary.sum_net_interest_of_total_assets_increase
		from companies, (
		 select company_code, sum(net_interest_of_total_assets_increase) as sum_net_interest_of_total_assets_increase
		 from report_summaries where category = 'Q4' and report_date >= ?
		 group by company_code order by sum(net_interest_of_total_assets_increase) desc limit 100) as summary
		where summary.company_code = companies.code order by summary.sum_net_interest_of_total_assets_increase desc;
	`
	var roaIncreases = make([]roaIncreaseResponse, 0)
	years, _ := strconv.Atoi(ctx.DefaultQuery("years", "2"))
	now := time.Now()
	currentYear, _, _ := now.Date()
	currentLocation := now.Location()

	year := time.Date(currentYear-years, 12, 31, 0, 0, 0, 0, currentLocation)
	db.Raw(sql, year.Unix()*1000).Scan(&roaIncreases)
	ctx.JSON(http.StatusOK, gin.H{
		"roaIncreases": roaIncreases,
		"message":      "ok",
	})
}

type roeResponse struct {
	Code      string
	Name      string
	SumAvgRoe float64
}

// Roe for roe top 100
func Roe(ctx *gin.Context) {
	sql := `
		select companies.code, companies.name, summary.sum_avg_roe
		from companies, (
		 select company_code, sum(avg_roe) as sum_avg_roe
		 from report_summaries where category = 'Q4' and report_date >= ?
		 group by company_code order by sum(avg_roe) desc limit 100) as summary
		where summary.company_code = companies.code order by summary.sum_avg_roe desc;
	`
	var roes = make([]roeResponse, 0)
	years, _ := strconv.Atoi(ctx.DefaultQuery("years", "2"))
	now := time.Now()
	currentYear, _, _ := now.Date()
	currentLocation := now.Location()

	year := time.Date(currentYear-years, 1, 1, 0, 0, 0, 0, currentLocation)
	db.Raw(sql, year.Unix()*1000).Scan(&roes)
	ctx.JSON(http.StatusOK, gin.H{
		"roes":    roes,
		"message": "ok",
	})
}

type roeIncreaseResponse struct {
	Code              string
	Name              string
	SumAvgRoeIncrease float64
}

// RoeIncrease for RoeIncrease top 100
func RoeIncrease(ctx *gin.Context) {
	sql := `
		select companies.code, companies.name, summary.sum_avg_roe_increase
		from companies, (
		 select company_code, sum(avg_roe_increase) as sum_avg_roe_increase
		 from report_summaries where category = 'Q4' and report_date >= ?
		 group by company_code order by sum(avg_roe_increase) desc limit 100) as summary
		where summary.company_code = companies.code order by summary.sum_avg_roe_increase desc;
	`
	var roeIncreases = make([]roeIncreaseResponse, 0)
	years, _ := strconv.Atoi(ctx.DefaultQuery("years", "2"))
	now := time.Now()
	currentYear, _, _ := now.Date()
	currentLocation := now.Location()

	year := time.Date(currentYear-years, 1, 1, 0, 0, 0, 0, currentLocation)
	db.Raw(sql, year.Unix()*1000).Scan(&roeIncreases)
	ctx.JSON(http.StatusOK, gin.H{
		"roeIncreases": roeIncreases,
		"message":      "ok",
	})
}
