#!/bin/bash

set -x

bin_file="my-companies-be"
new_bin_file="${bin_file}-new"
work_home="~/my_companies_be/"

GOOS=linux GOARCH=amd64 go build -o ./${new_bin_file}

scp ${new_bin_file} rails:${work_home}

ssh -T rails << 'EOF'
  cd ~/my_companies_be
  git pull
  psql my_companies < ~/my_companies_be/db/schema.sql
  systemctl --user stop vnote.club.service
  mv ~/my_companies_be/my-companies-be-new ~/my_companies_be/my-companies-be
  systemctl --user start vnote.club.service
EOF
