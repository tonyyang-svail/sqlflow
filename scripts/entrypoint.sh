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

# Start MySQL service in entrypoint since it requires root
set -e
echo begin entrypoint.sh
service mysql start
echo end entrypoint.sh

# Save command line to a temporary file to avoid arguemnt confliction in su --login
echo $@ > /tmp/start.sh
chmod 777 /tmp/start.sh

# Switch to user jovyan and execute the command
echo "running command $@"
su --login jovyan -c "source /miniconda/bin/activate sqlflow-dev && /tmp/start.sh"
# su jovyan -c '$@'

