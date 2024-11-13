#### 初始化配置

```shell
[root@centos7 ~]# hostnamectl set-hostname centos7.9-base

[root@centos7 ~]# vim /etc/sysconfig/network-scripts/ifcfg-ens32
TYPE=Ethernet
BOOTPROTO=static
NAME=ens32
DEVICE=ens32
ONBOOT=yes
IPADDR=10.0.1.10
NETMASK=255.255.255.0
GATEWAY=10.0.1.2
DNS1=114.114.114.114

[root@centos7 ~]# systemctl restart network
[root@centos7 ~]# ping www.qq.com
PING ins-r23tsuuf.ias.tencent-cloud.net (61.241.54.211) 56(84) bytes of data.
64 bytes from 61.241.54.211 (61.241.54.211): icmp_seq=1 ttl=128 time=9.96 ms
64 bytes from 61.241.54.211 (61.241.54.211): icmp_seq=2 ttl=128 time=10.3 ms

[root@centos7 ~]# vi centos7-init.sh
[root@centos7 ~]# bash centos7-init.sh
```

centos7优化脚本`centos7-init.sh`

```shell
#!/bin/bash

# Control switch
[[ "$1" != "" ]] && CentOS_ver="$1" || CentOS_ver='7'
[[ "$2" != "" ]] && iptables_yn="$2" || iptables_yn='n'

# Close SELINUX
setenforce 0
sed -i 's/^SELINUX=.*$/SELINUX=disabled/' /etc/selinux/config

# Custom profile
cat > /etc/profile.d/boge.sh << EOF
HISTSIZE=10000
PS1="\[\e[37;40m\][\[\e[32;40m\]\u\[\e[37;40m\]@\h \[\e[35;40m\]\W\[\e[0m\]]\\\\$ "
HISTTIMEFORMAT="%F %T \$(whoami) "

alias l='ls -AFhlt'
alias lh='l | head'

GREP_OPTIONS="--color=auto"
alias grep='grep --color'
alias egrep='egrep --color'
alias fgrep='fgrep --color'
EOF

[ -z "$(grep ^'PROMPT_COMMAND=' /etc/bashrc)" ] && cat >> /etc/bashrc << EOF
PROMPT_COMMAND='{ msg=\$(history 1 | { read x y; echo \$y; });logger "[euid=\$(whoami)]":\$(who am i):[\`pwd\`]"\$msg"; }'
EOF

# Change apt-get source list
# https://developer.aliyun.com/mirror/
yum install wget -y
cd /etc/yum.repos.d/
mkdir bak
mv *.repo bak
wget -O /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-${CentOS_ver}.repo
wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-${CentOS_ver}.repo
yum clean all
yum makecache

# Install package
yum groupinstall -y "base"
yum groupinstall -y "compatibility libraries"
yum groupinstall -y "debuging tools"
yum groupinstall -y "development tools"
yum install -y deltarpm gcc gcc-c++ make cmake autoconf libjpeg libjpeg-devel libpng libpng-devel freetype freetype-devel libxml2 libxml2-devel zlib zlib-devel glibc glibc-devel glib2 glib2-devel bzip2 bzip2-devel ncurses ncurses-devel libaio readline-devel curl curl-devel e2fsprogs e2fsprogs-devel krb5-devel libidn libidn-devel openssl openssl-devel libxslt-devel libicu-devel libevent-devel libtool libtool-ltdl bison gd-devel vim-enhanced pcre-devel zip unzip ntpdate patch bc expect rsync git lsof vim telnet tree nmap sysstat lrzsz dos2unix iotop iftop nethogs nload net-tools bash-completion sshpass
#yum  update -y
# clean yum installd cache
#find  /var/cache/yum/ -type f -exec rm {} \;
cd -

# /etc/security/limits.conf
[ -e /etc/security/limits.d/*nproc.conf ] && rename nproc.conf nproc.conf_bk /etc/security/limits.d/*nproc.conf
sed -i '/^# End of file/,$d' /etc/security/limits.conf
cat >> /etc/security/limits.conf <<EOF
# End of file
* soft nproc 1000000
* hard nproc 1000000
* soft nofile 1000000
* hard nofile 1000000
EOF

ulimit -SHn 1000000

# /etc/hosts
[ "$(hostname -i | awk '{print $1}')" != "127.0.0.1" ] && sed -i "s@127.0.0.1.*localhost@&\n127.0.0.1 $(hostname)@g" /etc/hosts

# Set timezone
rm -rf /etc/localtime
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# Set DNS
#cat > /etc/resolv.conf << EOF
#nameserver 114.114.114.114
#nameserver 8.8.8.8
#EOF

# ip_conntrack table full dropping packets
[ ! -e "/etc/sysconfig/modules/iptables.modules" ] && { echo -e "modprobe nf_conntrack\nmodprobe nf_conntrack_ipv4" > /etc/sysconfig/modules/iptables.modules; chmod +x /etc/sysconfig/modules/iptables.modules; }
modprobe nf_conntrack
modprobe nf_conntrack_ipv4
echo options nf_conntrack hashsize=131072 > /etc/modprobe.d/nf_conntrack.conf

# /etc/sysctl.conf
[ ! -e "/etc/sysctl.conf_bk" ] && /bin/mv /etc/sysctl.conf{,_bk}
cat > /etc/sysctl.conf << EOF
fs.file-max=1000000
net.ipv4.tcp_max_tw_buckets = 6000
net.ipv4.tcp_sack = 1
net.ipv4.tcp_window_scaling = 1
net.ipv4.tcp_rmem = 4096 87380 4194304
net.ipv4.tcp_wmem = 4096 16384 4194304
net.ipv4.tcp_max_syn_backlog = 16384
net.core.netdev_max_backlog = 32768
net.core.somaxconn = 32768
net.core.wmem_default = 8388608
net.core.rmem_default = 8388608
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_timestamps = 1
net.ipv4.tcp_fin_timeout = 20
net.ipv4.tcp_synack_retries = 2
net.ipv4.tcp_syn_retries = 2
net.ipv4.tcp_syncookies = 1
#net.ipv4.tcp_tw_len = 1
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_mem = 94500000 915000000 927000000
net.ipv4.tcp_max_orphans = 3276800
net.ipv4.ip_local_port_range = 1024 65000
net.nf_conntrack_max = 6553500
net.netfilter.nf_conntrack_max = 6553500
net.netfilter.nf_conntrack_tcp_timeout_close_wait = 60
net.netfilter.nf_conntrack_tcp_timeout_fin_wait = 120
net.netfilter.nf_conntrack_tcp_timeout_time_wait = 120
net.netfilter.nf_conntrack_tcp_timeout_established = 3600
EOF
sysctl -p

if [ "${CentOS_ver}" == '5' ]; then
  sed -i 's@^[3-6]:2345:respawn@#&@g' /etc/inittab
  sed -i 's@^ca::ctrlaltdel@#&@' /etc/inittab
  sed -i 's@LANG=.*$@LANG="en_US.UTF-8"@g' /etc/sysconfig/i18n
elif [ "${CentOS_ver}" == '6' ]; then
  sed -i 's@^ACTIVE_CONSOLES.*@ACTIVE_CONSOLES=/dev/tty[1-2]@' /etc/sysconfig/init
  sed -i 's@^start@#start@' /etc/init/control-alt-delete.conf
  sed -i 's@LANG=.*$@LANG="en_US.UTF-8"@g' /etc/sysconfig/i18n
elif [ "${CentOS_ver}" == '7' ]; then
  sed -i 's@LANG=.*$@LANG="en_US.UTF-8"@g' /etc/locale.conf
fi

# Update time
ntpdate pool.ntp.org
[ ! -e "/var/spool/cron/root" -o -z "$(grep 'ntpdate' /var/spool/cron/root)" ] && { echo "*/20 * * * * $(which ntpdate) pool.ntp.org > /dev/null 2>&1" >> /var/spool/cron/root;chmod 600 /var/spool/cron/root; }

# iptables
if [ "${iptables_yn}" == 'y' ]; then
  if [ -e "/etc/sysconfig/iptables" ] && [ -n "$(grep '^:INPUT DROP' /etc/sysconfig/iptables)" -a -n "$(grep 'NEW -m tcp --dport 22 -j ACCEPT' /etc/sysconfig/iptables)" -a -n "$(grep 'NEW -m tcp --dport 80 -j ACCEPT' /etc/sysconfig/iptables)" ]; then
    IPTABLES_STATUS=yes
  else
    IPTABLES_STATUS=no
  fi

  if [ "$IPTABLES_STATUS" == "no" ]; then
    [ -e "/etc/sysconfig/iptables" ] && /bin/mv /etc/sysconfig/iptables{,_bk}
    cat > /etc/sysconfig/iptables << EOF
# Firewall configuration written by system-config-securitylevel
# Manual customization of this file is not recommended.
*filter
:INPUT DROP [0:0]
:FORWARD ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
:syn-flood - [0:0]
-A INPUT -i lo -j ACCEPT
-A INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT
-A INPUT -p tcp -m state --state NEW -m tcp --dport 22 -j ACCEPT
-A INPUT -p tcp -m state --state NEW -m tcp --dport 80 -j ACCEPT
-A INPUT -p tcp -m state --state NEW -m tcp --dport 443 -j ACCEPT
-A INPUT -p icmp -m icmp --icmp-type 8 -j ACCEPT
COMMIT
EOF
  fi

  FW_PORT_FLAG=$(grep -ow "dport ${ssh_port}" /etc/sysconfig/iptables)
  [ -z "${FW_PORT_FLAG}" -a "${ssh_port}" != "22" ] && sed -i "s@dport 22 -j ACCEPT@&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport ${ssh_port} -j ACCEPT@" /etc/sysconfig/iptables
  /bin/cp /etc/sysconfig/{iptables,ip6tables}
  sed -i 's@icmp@icmpv6@g' /etc/sysconfig/ip6tables
  iptables-restore < /etc/sysconfig/iptables
  ip6tables-restore < /etc/sysconfig/ip6tables
  service iptables save
  service ip6tables save
  chkconfig --level 3 iptables on
  chkconfig --level 3 ip6tables on
fi

# sshd optimization
\cp /etc/ssh/sshd_config /etc/ssh/sshd_config.BAK
sed -i '/^AddressFamily.*/d' /etc/ssh/sshd_config
num=$(grep -n "^#ListenAddress.*" /etc/ssh/sshd_config|awk -F: 'NR==1{print $1}')
[[ -z `grep -n "^UseDNS.*" /etc/ssh/sshd_config` ]]&&sed -ir "${num} iUseDNS no" /etc/ssh/sshd_config||{
sed -ri 's+^#UseDNS.*+UseDNS no+g' /etc/ssh/sshd_config
}
[[ -z `grep -n "^GSSAPIAuthentication.*" /etc/ssh/sshd_config` ]]&&sed -ir "${num} iGSSAPIAuthentication no" /etc/ssh/sshd_config||{
sed -ri 's+^GSSAPIAuthentication .*+GSSAPIAuthentication no+g' /etc/ssh/sshd_config
}
service rsyslog restart
service sshd restart

. /etc/profile
```

