// Copyright 2022 Praetorian Security, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scan

// These import statements ensure that the init functions run in each plugin.
// When a new plugin is added, this list should be updated.

import (
	_ "fingerprintx/pkg/plugins/services/dhcp"
	_ "fingerprintx/pkg/plugins/services/dns"
	_ "fingerprintx/pkg/plugins/services/echo"
	_ "fingerprintx/pkg/plugins/services/ftp"
	_ "fingerprintx/pkg/plugins/services/http"
	_ "fingerprintx/pkg/plugins/services/imap"
	_ "fingerprintx/pkg/plugins/services/ipmi"
	_ "fingerprintx/pkg/plugins/services/ipsec"
	_ "fingerprintx/pkg/plugins/services/jdwp"
	_ "fingerprintx/pkg/plugins/services/kafka/kafkaNew"
	_ "fingerprintx/pkg/plugins/services/kafka/kafkaOld"
	_ "fingerprintx/pkg/plugins/services/ldap"
	_ "fingerprintx/pkg/plugins/services/linuxrpc"
	_ "fingerprintx/pkg/plugins/services/modbus"
	_ "fingerprintx/pkg/plugins/services/mongodb"
	_ "fingerprintx/pkg/plugins/services/mqtt/mqtt3"
	_ "fingerprintx/pkg/plugins/services/mqtt/mqtt5"
	_ "fingerprintx/pkg/plugins/services/mssql"
	_ "fingerprintx/pkg/plugins/services/mysql"
	_ "fingerprintx/pkg/plugins/services/netbios"
	_ "fingerprintx/pkg/plugins/services/ntp"
	_ "fingerprintx/pkg/plugins/services/openvpn"
	_ "fingerprintx/pkg/plugins/services/oracledb"
	_ "fingerprintx/pkg/plugins/services/pop3"
	_ "fingerprintx/pkg/plugins/services/postgresql"
	_ "fingerprintx/pkg/plugins/services/rdp"
	_ "fingerprintx/pkg/plugins/services/redis"
	_ "fingerprintx/pkg/plugins/services/rsync"
	_ "fingerprintx/pkg/plugins/services/rtsp"
	_ "fingerprintx/pkg/plugins/services/smb"
	_ "fingerprintx/pkg/plugins/services/smtp"
	_ "fingerprintx/pkg/plugins/services/snmp"
	_ "fingerprintx/pkg/plugins/services/ssh"
	_ "fingerprintx/pkg/plugins/services/stun"
	_ "fingerprintx/pkg/plugins/services/telnet"
	_ "fingerprintx/pkg/plugins/services/vnc"
)
