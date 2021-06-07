CREATE TABLE IF NOT EXISTS companies (id SERIAL PRIMARY KEY, name varchar(40));

ALTER TABLE
  companies
ADD
  COLUMN deleted_at timestamptz;

CREATE TABLE IF NOT EXISTS profits (
  id SERIAL,
  year char(20),
  ying_shou bigint,
  ying_ye_cheng_ben bigint,
  fei_ying_shou bigint,
  li_run bigint,
  jing_li_run bigint,
  company_id bigint,
  CONSTRAINT fk_profits_companies FOREIGN KEY (company_id) REFERENCES companies (id)
);

ALTER TABLE
  profits
ADD
  COLUMN IF NOT EXISTS deleted_at timestamptz;

ALTER table
  profits
alter column
  year type char(4);

ALTER table
  companies
add
  column IF NOT EXISTS created_at timestamptz;

ALTER table
  companies
add
  column IF NOT EXISTS updated_at timestamptz;

ALTER table
  profits
add
  column IF NOT EXISTS created_at timestamptz;

ALTER table
  profits
add
  column IF NOT EXISTS updated_at timestamptz;

CREATE INDEX IF NOT EXISTS companies_deleted_at_index ON companies USING btree(deleted_at);

CREATE INDEX IF NOT EXISTS Profits_deleted_at_index ON profits USING btree(deleted_at);

/*
* comment for chart  
 */
CREATE TABLE IF NOT EXISTS comments (
  id SERIAL,
  chart varchar(20),
  content text,
  company_id bigint,
  user_id bigint,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);
CREATE INDEX IF NOT EXISTS Comments_deleted_at_index ON comments USING btree(deleted_at);
CREATE INDEX IF NOT EXISTS Comments_company_id_user_id_index ON comments (user_id, company_id);


ALTER table
  companies
add
  column IF NOT EXISTS code char(8);

CREATE INDEX IF NOT EXISTS Companies_code_index ON companies (code);


CREATE TABLE IF NOT EXISTS report_summaries (
  id SERIAL,
  category char(2),
  company_code char(8),
  report_name char(10),
  report_date bigint,
  avg_roe float,
  avg_roe_increase float,
  np_per_share float,
  np_per_share_increase float,
  operate_cash_flow_ps float,
  operate_cash_flow_ps_increase float,
  basic_eps float,
  basic_eps_increase float,
  capital_reserve float,
  capital_reserve_increase float,
  undistri_profit_ps float,
  undistri_profit_ps_increase float,
  net_interest_of_total_assets float,
  net_interest_of_total_assets_increase float,
  net_selling_rate float,
  net_selling_rate_increase float,
  gross_selling_rate float,
  gross_selling_rate_increase float,
  total_revenue float,
  total_revenue_increase float,
  operating_income_yoy float,
  operating_income_yoy_increase float,
  net_profit_atsopc float,
  net_profit_atsopc_increase float,
  net_profit_atsopc_yoy float,
  net_profit_atsopc_yoy_increase float,
  net_profit_after_nrgal_atsolc float,
  net_profit_after_nrgal_atsolc_increase float,
  np_atsopc_nrgal_yoy float,
  np_atsopc_nrgal_yoy_increase float,
  ore_dlt float,
  ore_dlt_increase float,
  rop float,
  rop_increase float,
  asset_liab_ratio float,
  asset_liab_ratio_increase float,
  current_ratio float,
  current_ratio_increase float,
  quick_ratio float,
  quick_ratio_increase float,
  equity_multiplier float,
  equity_multiplier_increase float,
  equity_ratio float,
  equity_ratio_increase float,
  holder_equity float,
  holder_equity_increase float,
  ncf_from_oa_to_total_liab float,
  ncf_from_oa_to_total_liab_increase float,
  inventory_turnover_days float,
  inventory_turnover_days_increase float,
  receivable_turnover_days float,
  receivable_turnover_days_increase float,
  accounts_payable_turnover_days float,
  accounts_payable_turnover_days_increase float,
  cash_cycle float,
  cash_cycle_increase float,
  operating_cycle float,
  operating_cycle_increase float,
  total_capital_turnover float,
  total_capital_turnover_increase float,
  inventory_turnover float,
  inventory_turnover_increase float,
  account_receivable_turnover float,
  account_receivable_turnover_increase float,
  accounts_payable_turnover float,
  accounts_payable_turnover_increase float,
  current_asset_turnover_rate float,
  current_asset_turnover_rate_increase float,
  fixed_asset_turnover_ratio float,
  fixed_asset_turnover_ratio_increase float,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS Report_summary_code_index ON  report_summaries (company_code);
CREATE INDEX IF NOT EXISTS Report_summary_category_index ON report_summaries (category);


CREATE TABLE IF NOT EXISTS incomes (
  id SERIAL,
  category char(2),
  company_code char(8),
  report_name char(10),
  report_date bigint,
  net_profit float,
  net_profit_increase float,
  net_profit_atsopc float,
  net_profit_atsopc_increase float,
  total_revenue float,
  total_revenue_increase float,
  op float,
  op_increase float,
  income_from_chg_in_fv float,
  income_from_chg_in_fv_increase float,
  invest_incomes_from_rr float,
  invest_incomes_from_rr_increase float,
  invest_income float,
  invest_income_increase float,
  exchg_gain float,
  exchg_gain_increase float,
  operating_taxes_and_surcharge float,
  operating_taxes_and_surcharge_increase float,
  asset_impairment_loss float,
  asset_impairment_loss_increase float,
  non_operating_income float,
  non_operating_income_increase float,
  non_operating_payout float,
  non_operating_payout_increase float,
  profit_total_amt float,
  profit_total_amt_increase float,
  minority_gal float,
  minority_gal_increase float,
  basic_eps float,
  basic_eps_increase float,
  dlt_earnings_per_share float,
  dlt_earnings_per_share_increase float,
  othr_compre_income_atoopc float,
  othr_compre_income_atoopc_increase float,
  othr_compre_income_atms float,
  othr_compre_income_atms_increase float,
  total_compre_income float,
  total_compre_income_increase float,
  total_compre_income_atsopc float,
  total_compre_income_atsopc_increase float,
  total_compre_income_atms float,
  total_compre_income_atms_increase float,
  othr_compre_income float,
  othr_compre_income_increase float,
  net_profit_after_nrgal_atsolc float,
  net_profit_after_nrgal_atsolc_increase float,
  income_tax_expenses float,
  income_tax_expenses_increase float,
  credit_impairment_loss float,
  credit_impairment_loss_increase float,
  revenue float,
  revenue_increase float,
  operating_costs float,
  operating_costs_increase float,
  operating_cost float,
  operating_cost_increase float,
  sales_fee float,
  sales_fee_increase float,
  manage_fee float,
  manage_fee_increase float,
  financing_expenses float,
  financing_expenses_increase float,
  rad_cost float,
  rad_cost_increase float,
  finance_cost_interest_fee float,
  finance_cost_interest_fee_increase float,
  finance_cost_interest_income float,
  finance_cost_interest_income_increase float,
  asset_disposal_income float,
  asset_disposal_income_increase float,
  other_income float,
  other_income_increase float,
  noncurrent_asset_disposal_gain float,
  noncurrent_asset_disposal_gain_increase float,
  noncurrent_asset_disposal_loss float,
  noncurrent_asset_disposal_loss_increase float,
  net_profit_bi float,
  net_profit_bi_increase float,
  continous_operating_np float,
  continous_operating_np_increase float,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS income_code_index ON  incomes (company_code);
CREATE INDEX IF NOT EXISTS income_category_index ON incomes (category);

CREATE TABLE  IF NOT EXISTS cash_flows (
  id SERIAL,
  category char(2),
  company_code char(8),
  report_name char(10),
  report_date bigint,
	ncf_from_oa float,
	ncf_from_oa_increase float,
	ncf_from_ia float,
  ncf_from_ia_increase float,
	ncf_from_fa float,
	ncf_from_fa_increase float,
	cash_received_of_othr_oa float,
	cash_received_of_othr_oa_increase float,
	sub_total_of_ci_from_oa float,
	sub_total_of_ci_from_oa_increase float,
	cash_paid_to_employee_etc float,
	cash_paid_to_employee_etc_increase float,
	payments_of_all_taxes float,
	payments_of_all_taxes_increase float,
	othrcash_paid_relating_to_oa float,
	othrcash_paid_relating_to_oa_increase float,
	sub_total_of_cos_from_oa float,
	sub_total_of_cos_from_oa_increase float,
	cash_received_of_dspsl_invest float,
	cash_received_of_dspsl_invest_increase float,
	invest_income_cash_received float,
	invest_income_cash_received_increase float,
	net_cash_of_disposal_assets float,
	net_cash_of_disposal_assets_increase float,
	net_cash_of_disposal_branch float,
	net_cash_of_disposal_branch_increase float,
	cash_received_of_othr_ia float,
	cash_received_of_othr_ia_increase float,
	sub_total_of_ci_from_ia float,
	sub_total_of_ci_from_ia_increase float,
	invest_paid_cash float,
	invest_paid_cash_increase float,
	cash_paid_for_assets float,
	cash_paid_for_assets_increase float,
	othrcash_paid_relating_to_ia float,
	othrcash_paid_relating_to_ia_increase float,
	sub_total_of_cos_from_ia float,
	sub_total_of_cos_from_ia_increase float,
	cash_received_of_absorb_invest float,
	cash_received_of_absorb_invest_increase float,
	cash_received_from_investor float,
	cash_received_from_investor_increase float,
	cash_received_from_bond_issue float,
	cash_received_from_bond_issue_increase float,
	cash_received_of_borrowing float,
	cash_received_of_borrowing_increase float,
	cash_received_of_othr_fa float,
	cash_received_of_othr_fa_increase float,
	sub_total_of_ci_from_fa float,
	sub_total_of_ci_from_fa_increase float,
	cash_pay_for_debt float,
	cash_pay_for_debt_increase float,
	cash_paid_of_distribution float,
	cash_paid_of_distribution_increase float,
	branch_paid_to_minority_holder float,
	branch_paid_to_minority_holder_increase float,
	othrcash_paid_relating_to_fa float,
	othrcash_paid_relating_to_fa_increase float,
	sub_total_of_cos_from_fa float,
	sub_total_of_cos_from_fa_increase float,
	effect_of_exchange_chg_on_cce float,
	effect_of_exchange_chg_on_cce_increase float,
	net_increase_in_cce float,
	net_increase_in_cce_increase float,
	initial_balance_of_cce float,
	initial_balance_of_cce_increase float,
	final_balance_of_cce float,
	final_balance_of_cce_increase float,
	cash_received_of_sales_service float,
	cash_received_of_sales_service_increase float,
	refund_of_tax_and_levies float,
	refund_of_tax_and_levies_increase float,
	goods_buy_and_service_cash_pay float,
	goods_buy_and_service_cash_pay_increase float,
	net_cash_amt_from_branch float,
	net_cash_amt_from_branch_increase float,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS cash_flow_code_index ON  cash_flows (company_code);