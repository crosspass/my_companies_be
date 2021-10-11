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

CREATE INDEX IF NOT EXISTS Report_summary_code_index ON report_summaries (company_code);

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

CREATE INDEX IF NOT EXISTS income_code_index ON incomes (company_code);

CREATE INDEX IF NOT EXISTS income_category_index ON incomes (category);

CREATE TABLE IF NOT EXISTS cash_flows (
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

CREATE INDEX IF NOT EXISTS cash_flow_code_index ON cash_flows (company_code);

CREATE TABLE IF NOT EXISTS balances (
  id SERIAL,
  category char(2),
  company_code char(8),
  report_name char(10),
  report_date bigint,
  total_assets float,
  total_assets_increase float,
  total_liab float,
  total_liab_increase float,
  asset_liab_ratio float,
  asset_liab_ratio_increase float,
  total_quity_atsopc float,
  total_quity_atsopc_increase float,
  tradable_fnncl_assets float,
  tradable_fnncl_assets_increase float,
  interest_receivable float,
  interest_receivable_increase float,
  saleable_finacial_assets float,
  saleable_finacial_assets_increase float,
  held_to_maturity_invest float,
  held_to_maturity_invest_increase float,
  fixed_asset float,
  fixed_asset_increase float,
  intangible_assets float,
  intangible_assets_increase float,
  construction_in_process float,
  construction_in_process_increase float,
  dt_assets float,
  dt_assets_increase float,
  tradable_fnncl_liab float,
  tradable_fnncl_liab_increase float,
  payroll_payable float,
  payroll_payable_increase float,
  tax_payable float,
  tax_payable_increase float,
  estimated_liab float,
  estimated_liab_increase float,
  dt_liab float,
  dt_liab_increase float,
  bond_payable float,
  bond_payable_increase float,
  shares float,
  shares_increase float,
  capital_reserve float,
  capital_reserve_increase float,
  earned_surplus float,
  earned_surplus_increase float,
  undstrbtd_profit float,
  undstrbtd_profit_increase float,
  minority_equity float,
  minority_equity_increase float,
  total_holders_equity float,
  total_holders_equity_increase float,
  total_liab_and_holders_equity float,
  total_liab_and_holders_equity_increase float,
  lt_equity_invest float,
  lt_equity_invest_increase float,
  derivative_fnncl_liab float,
  derivative_fnncl_liab_increase float,
  general_risk_provision float,
  general_risk_provision_increase float,
  frgn_currency_convert_diff float,
  frgn_currency_convert_diff_increase float,
  goodwill float,
  goodwill_increase float,
  invest_property float,
  invest_property_increase float,
  interest_payable float,
  interest_payable_increase float,
  treasury_stock float,
  treasury_stock_increase float,
  othr_compre_income float,
  othr_compre_income_increase float,
  othr_equity_instruments float,
  othr_equity_instruments_increase float,
  currency_funds float,
  currency_funds_increase float,
  bills_receivable float,
  bills_receivable_increase float,
  account_receivable float,
  account_receivable_increase float,
  pre_payment float,
  pre_payment_increase float,
  dividend_receivable float,
  dividend_receivable_increase float,
  othr_receivables float,
  othr_receivables_increase float,
  inventory float,
  inventory_increase float,
  nca_due_within_one_year float,
  nca_due_within_one_year_increase float,
  othr_current_assets float,
  othr_current_assets_increase float,
  current_assets_si float,
  current_assets_si_increase float,
  total_current_assets float,
  total_current_assets_increase float,
  lt_receivable float,
  lt_receivable_increase float,
  dev_expenditure float,
  dev_expenditure_increase float,
  lt_deferred_expense float,
  lt_deferred_expense_increase float,
  othr_noncurrent_assets float,
  othr_noncurrent_assets_increase float,
  noncurrent_assets_si float,
  noncurrent_assets_si_increase float,
  total_noncurrent_assets float,
  total_noncurrent_assets_increase float,
  st_loan float,
  st_loan_increase float,
  bill_payable float,
  bill_payable_increase float,
  accounts_payable float,
  accounts_payable_increase float,
  pre_receivable float,
  pre_receivable_increase float,
  dividend_payable float,
  dividend_payable_increase float,
  othr_payables float,
  othr_payables_increase float,
  noncurrent_liab_due_in1y float,
  noncurrent_liab_due_in1y_increase float,
  current_liab_si float,
  current_liab_si_increase float,
  total_current_liab float,
  total_current_liab_increase float,
  lt_loan float,
  lt_loan_increase float,
  lt_payable float,
  lt_payable_increase float,
  special_payable float,
  special_payable_increase float,
  othr_non_current_liab float,
  othr_non_current_liab_increase float,
  noncurrent_liab_si float,
  noncurrent_liab_si_increase float,
  total_noncurrent_liab float,
  total_noncurrent_liab_increase float,
  salable_financial_assets float,
  salable_financial_assets_increase float,
  othr_current_liab float,
  othr_current_liab_increase float,
  ar_and_br float,
  ar_and_br_increase float,
  contractual_assets float,
  contractual_assets_increase float,
  bp_and_ap float,
  bp_and_ap_increase float,
  contract_liabilities float,
  contract_liabilities_increase float,
  to_sale_asset float,
  to_sale_asset_increase float,
  other_eq_ins_invest float,
  other_eq_ins_invest_increase float,
  other_illiquid_fnncl_assets float,
  other_illiquid_fnncl_assets_increase float,
  fixed_asset_sum float,
  fixed_asset_sum_increase float,
  fixed_assets_disposal float,
  fixed_assets_disposal_increase float,
  construction_in_process_sum float,
  construction_in_process_sum_increase float,
  project_goods_and_material float,
  project_goods_and_material_increase float,
  productive_biological_assets float,
  productive_biological_assets_increase float,
  oil_and_gas_asset float,
  oil_and_gas_asset_increase float,
  to_sale_debt float,
  to_sale_debt_increase float,
  lt_payable_sum float,
  lt_payable_sum_increase float,
  noncurrent_liab_di float,
  noncurrent_liab_di_increase float,
  perpetual_bond float,
  perpetual_bond_increase float,
  special_reserve float,
  special_reserve_increase float,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS balance_code_index ON balances (company_code);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL,
  user_name varchar(20),
  full_name varchar(20),
  email varchar(20),
  register_token text,
  active_time timestamptz,
  password_hash text,
  password_salt text,
  is_actived boolean,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS sessions (
  id SERIAL,
  user_id int not null,
  key varchar,
  login_time timestamptz,
  last_seen_time timestamptz,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS sessions_user_idx on sessions (id);
CREATE INDEX IF NOT EXISTS sessoins_key_index ON sessions (key);

/*
 * articles
 */
CREATE TABLE IF NOT EXISTS articles (
  id SERIAL,
  user_id bigint,
  content text,
  raw_content text,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);
CREATE INDEX IF NOT EXISTS articles_user_id_index ON articles (user_id);

/*
 * users_companies
 */
CREATE TABLE IF NOT EXISTS users_companies (
  user_id bigint,
  company_id bigint
);
CREATE INDEX IF NOT EXISTS users_companies_user_id_index ON users_companies (user_id);
CREATE INDEX IF NOT EXISTS users_companies_company_id_index ON users_companies (company_id);
CREATE UNIQUE INDEX IF NOT EXISTS users_companies_user_id_company_id_index ON users_companies (user_id, company_id);

CREATE TABLE IF NOT EXISTS articles_companies (
  article_id bigint,
  company_id bigint
);
CREATE INDEX IF NOT EXISTS articles_companies_article_id_index ON articles_companies (article_id);
CREATE INDEX IF NOT EXISTS articles_companies_company_id_index ON articles_companies (company_id);
CREATE UNIQUE INDEX IF NOT EXISTS articles_companies_article_id_company_id_index ON articles_companies (article_id, company_id);

/*
 * csvs
 */
CREATE TABLE IF NOT EXISTS csvs (
  id SERIAL,
  user_id bigint,
  company_id bigint,
  title varchar,
  chart_type varchar,
  data text,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);
CREATE INDEX IF NOT EXISTS csvs_user_id_index ON csvs (user_id);
CREATE INDEX IF NOT EXISTS csv_company_id_index ON csvs (company_id);
CREATE INDEX IF NOT EXISTS csv_user_company_id_index ON csvs (company_id, user_id);