// Copyright 2015 - António Meireles  <antonio.meireles@reformi.st>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"time"

	"github.com/spf13/viper"
)

type (
	sessionInfo struct {
		configDir, imageDir, runDir string
		pwd, uid, gid, username     string
		hasPowers, debug, json      bool
	}
	// VMInfo - per VM settings
	VMInfo struct {
		Name, Channel, Version                 string
		Cpus, Memory                           int
		UUID, Xhyve                            string
		CloudConfig                            string `json:",omitempty"`
		CClocation                             string `json:",omitempty"`
		SSHkey                                 string `json:",omitempty"`
		Extra                                  string `json:",omitempty"`
		Root                                   int
		Ethernet                               []NetworkInterface
		Storage                                storageAssets
		InternalSSHauthKey, InternalSSHprivKey string
		Detached                               bool
		Pid                                    int
		PublicIP                               string
		CreatedAt                              time.Time
	}
	// NetworkInterface ...
	NetworkInterface struct {
		Type int
		// if tap
		Path string
	}
	// StorageDevice ...
	StorageDevice struct {
		Slot       int
		Type, Path string
	}
	storageAssets struct {
		CDDrives   map[string]StorageDevice `json:",omitempty"`
		HardDrives map[string]StorageDevice `json:",omitempty"`
	}
)

var (
	// CoreOS public release streams
	DefaultChannels = []string{"alpha", "beta", "stable"}
	vipre           *viper.Viper
	// SessionContext ...
	SessionContext sessionInfo
)

const (
	_ = iota
	Raw
	Tap
	HDD          = "HDD"
	CDROM        = "CDROM"
	Local        = "localfs"
	Remote       = "URL"
	Attached     = true
	Detached     = false
	HelpTemplate = `{{ $cmd := . }}
Usage: {{if .Runnable}}
  {{.UseLine}}{{if .HasFlags}} [flags]{{end}}{{end}}{{if .HasSubCommands}}
  {{ .CommandPath}} [command]
  {{end}}
{{if gt .Aliases 0}}
Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}
Examples:
{{ .Example }}
{{end}}{{ if .HasAvailableSubCommands}}Available Commands: {{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{end}}{{ if .HasLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages}}{{end}}{{ if .HasInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics: {{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.
{{end}}
All flags can also be configured via upper-case environment variables prefixed with "COREOS_"
For example, "--debug" => "COREOS_DEBUG"
`
	// GPGLongId
	GPGLongID = "50E0885593D2DCB4"
	// GPGKey
	GPGKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v2
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
=1egC
-----END PGP PUBLIC KEY BLOCK-----
`
	//
	CoreOEMsetupEnv = `#!/bin/bash
[[ $(</proc/cmdline) =~ uuid=([^\ ]+) ]]; UUID=${BASH_REMATCH[1]}
[[ $(</proc/cmdline) =~ localuser=([^\ ]+) ]]; CALLERID=${BASH_REMATCH[1]}
STATUSDIR=/Users/${CALLERID}/.coreos/running/${UUID}
# wait for eth0 to get up...
while [ 1 ]; do
  COREOS_PUBLIC_IPV4=$(/bin/ifconfig eth0 | awk '/inet /{print $2}')
  if [ -n "${COREOS_PUBLIC_IPV4}" ]; then
     break
  fi
  sleep 1; done
COREOS_PRIVATE_IPV4=${COREOS_PUBLIC_IPV4}
( echo UUID=${UUID};
  echo CALLERID=${CALLERID};
  echo STATUSDIR=${STATUSDIR};
  echo COREOS_PUBLIC_IPV4=${COREOS_PUBLIC_IPV4};
  echo COREOS_PRIVATE_IPV4=${COREOS_PRIVATE_IPV4};
) > /etc/environment
[[ $(</proc/cmdline) =~ sshkey_internal=\"([^\"]+)\" ]]
echo "${BASH_REMATCH[1]}" | update-ssh-keys -a proc-cmdline-ssh_internal

`
	//
	CoreOEMsetup = `#cloud-config

write-files:
  - path: /etc/conf.d/nfs
    permissions: '0644'
    content: |
      OPTS_RPC_MOUNTD=""
coreos:
  units:
    - name: rpc-statd.service
      command: start
      enable: true
    - name: Users.mount
      command: start
      content: |
        [Mount]
          What=192.168.64.1:/Users
          Where=/Users
          Options=rw,async,nolock,noatime,rsize=32768,wsize=32768
          Type=nfs
    - name: local-cloud-config.service
      command: start
      enable: true
      content: |
        [Unit]
          Description=Load cloud-config from file
          Requires=xhyve.service
          After=xhyve.service Users.mount
          ConditionPathExists=/etc/environment
        [Service]
          Type=oneshot
          RemainAfterExit=yes
          EnvironmentFile=/etc/environment
          ExecStart=/bin/bash -c "[[ -f $STATUSDIR/cloud-config.local ]] && \
                        /usr/bin/coreos-cloudinit \
                            -from-file $STATUSDIR/cloud-config.local || true"
    - name: xhyve.service
      command: start
      content: |
        [Unit]
          Description=updates xhyve context
          Requires=coreos-setup-environment.service Users.mount
          After=coreos-setup-environment.service Users.mount
          ConditionPathExists=/etc/environment
        [Service]
          Type=oneshot
          RemainAfterExit=true
          EnvironmentFile=/etc/environment
          ExecStart=/bin/bash -c "echo ${COREOS_PUBLIC_IPV4} > ${STATUSDIR}/ip"
`
)