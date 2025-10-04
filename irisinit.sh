#!/bin/bash

if [ ! -f /tmp/password.done ]; then
  # close access to system users
  tr -dc '[:graph:]' </dev/urandom | head -c 16 > /tmp/password

  /usr/irissys/dev/Container/changePassword.sh /tmp/password

  rm -rf /tmp/password

  iris start $ISC_PACKAGE_INSTANCENAME
  /usr/irissys/dev/Cloud/ICM/waitISC.sh
fi

if [ ! -z $IRIS_NAMESPACE ] && ! iris session $ISC_PACKAGE_INSTANCENAME -U $IRIS_NAMESPACE '^%%' > /dev/null; then
echo "Creating Namespace: $IRIS_NAMESPACE" > /proc/1/fd/1
iris sql $ISC_PACKAGE_INSTANCENAME -U%SYS <<-EOSESS > /dev/null
CREATE DATABASE $IRIS_NAMESPACE
quit
EOSESS
fi

if [ ! -z $IRIS_USERNAME ] && [ ! -z $IRIS_PASSWORD ]; then
echo "Creating User: $IRIS_USERNAME"  > /proc/1/fd/1
iris session $ISC_PACKAGE_INSTANCENAME -U%SYS <<-EOSESS > /dev/null
check(sc)	if 'sc { do ##class(%SYSTEM.OBJ).DisplayError(sc) do ##class(%SYSTEM.Process).Terminate(, 1) }
set exists = ##class(Security.Users).Exists("$IRIS_USERNAME", .user)
if 'exists { set sc = ##class(Security.Users).Create("$IRIS_USERNAME", "%All", "$IRIS_PASSWORD", , "$IRIS_NAMESPACE") }
if exists,\$isobject(user) { set user.PasswordExternal = "$IRIS_PASSWORD", sc = user.%Save() }
do check(sc)
halt
EOSESS
fi

# start golang app
/myapp 1> /proc/1/fd/1 2> /proc/1/fd/2 &

touch $IRISSYS/iris-started
