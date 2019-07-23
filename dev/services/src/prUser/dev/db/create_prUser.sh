#!/bin/bash

#. ./dbcmd.sh

CMD=`which cockroach`" sql"
PORT="26257"
IP="127.0.0.1"
USER="root"

CMD=`echo $CMD "--insecure" --host=$IP:$PORT`

echo "$CMD"

#1 Create kevlarAdmin if it doesn not already exists
$CMD < PavedRoadAdmin.sql

#2 Create kevlarWeb db
$CMD < pavedroad.sql

#3 Create kevlarAdmin all on kevlarWeb db
$CMD < grantPavedRoadAdminAll.sql

#4 Create prTokenTable 
$CMD < createPrUserTable.sql



