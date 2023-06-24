#!/bin/bash

source ./config/maker.env
echo $TEMPLATE_PATH
echo $GO_MOD_URL
bash $TEMPLATE_PATH/go-gin-bearerauth-postgres-monolithic-template-crud.sh
make swag_init