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
