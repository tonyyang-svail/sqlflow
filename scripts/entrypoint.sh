# Copyright 2019 The SQLFlow Authors. All rights reserved.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

######################################################################
# Borrowed from start.sh
######################################################################
SQLFLOW_MYSQL_HOST=${SQLFLOW_MYSQL_HOST:-127.0.0.1}
SQLFLOW_MYSQL_PORT=${SQLFLOW_MYSQL_PORT:-3306}
SQLFLOW_NOTEBOOK_DIR=${SQLFLOW_NOTEBOOK_DIR:-/workspace}

function sleep_until_mysql_is_ready() {
  until mysql -u root -proot --host ${SQLFLOW_MYSQL_HOST} --port ${SQLFLOW_MYSQL_PORT} -e ";" ; do
    sleep 1
    read -p "Can't connect, retrying..."
  done
}

function setup_mysql() {
  # Start mysqld
  sed -i "s/.*bind-address.*/bind-address = ${SQLFLOW_MYSQL_HOST}/" /etc/mysql/mysql.conf.d/mysqld.cnf
  service mysql start
  sleep 1
  populate_example_dataset
  # Grant all privileges to all the remote hosts so that the sqlflow server can
  # be scaled to more than one replicas.
  # NOTE: should notice this authorization on the production environment, it's not safe.
  mysql -uroot -proot -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'' IDENTIFIED BY 'root' WITH GRANT OPTION;"
}

function populate_example_dataset() {
  sleep_until_mysql_is_ready
  # FIXME(typhoonzero): should let docker-entrypoint.sh do this work
  for f in /docker-entrypoint-initdb.d/*; do
    cat $f | mysql -uroot -proot --host ${SQLFLOW_MYSQL_HOST} --port ${SQLFLOW_MYSQL_PORT}
  done
}

function setup_sqlflow_server() {
  sleep_until_mysql_is_ready

  DS="mysql://root:root@tcp(${SQLFLOW_MYSQL_HOST}:${SQLFLOW_MYSQL_PORT})/?maxAllowedPacket=0"
  echo "Connect to the datasource ${DS}"
  # Start sqlflowserver
  sqlflowserver --datasource=${DS}
}

######################################################################
# End of borrowing from start.sh
######################################################################

echo "setup all-in-one"
setup_mysql
setup_sqlflow_server &
cp -r ${SQLFLOW_NOTEBOOK_DIR} ${HOME}

# Save command line to a temporary file to avoid arguemnt confliction in su --login
echo "export IPYTHON_STARTUP=${IPYTHON_STARTUP}
export SQLFLOW_MYSQL_HOST=${SQLFLOW_MYSQL_HOST}
export SQLFLOW_MYSQL_PORT=${SQLFLOW_MYSQL_PORT}
export SQLFLOW_NOTEBOOK_DIR=${SQLFLOW_NOTEBOOK_DIR}
export SQLFLOW_SERVER=localhost:50051" > /tmp/start.sh

echo $@ >> /tmp/start.sh

chmod 777 /tmp/start.sh
# Switch to user jovyan and execute the command
echo "running command $@"
su --login jovyan -c "source /miniconda/bin/activate sqlflow-dev && /tmp/start.sh"
# su jovyan -c '$@'

