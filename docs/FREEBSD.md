# FreeBSD 安装与兼容性

目标系统：FreeBSD 15.1 / amd64。FreeBSD 安装器使用系统自带 `/bin/sh`、`fetch`、`rc.d`、`sysrc` 和 `daemon`；不需要先安装 Bash、wget 或 sudo。需要以 root 执行安装。

## 安装 StarHanami Snapshot

将 `PLACEHOLDER` 替换成此设备的 Komari Agent Token。以下示例保留停用远程控制和每月 1 日流量重置；安装器只读取本 fork 的 Releases。

```sh
fetch -o /tmp/komari-agent-install.sh https://raw.githubusercontent.com/StarHanami/komari-agent/main/install.sh
sh /tmp/komari-agent-install.sh --install-version snapshot \
  -e https://monitor.irenelyra.edu.kg \
  -t 'PLACEHOLDER' --disable-web-ssh --month-rotate 1
```

`--install-version snapshot` 选择本 fork 最新 `Snapshot-*` 预发布。也可指定一个准确的 `Snapshot-yymmddhhMM` tag。安装完成后，Agent 会继续跟踪本 fork 的 Snapshot 更新；`daemon` 会在更新后的进程退出时重新启动新二进制。

标准 FreeBSD 基础系统使用自带的 `fetch` 完成 API 查询和 HTTPS 下载，无需安装额外软件包。仅在裁剪系统缺少 `fetch` 和 `curl` 时，安装器才尝试通过现有 pkg 仓库安装 `curl` 和 `ca_root_nss`；不会切换 latest/quarterly 仓库或升级 FreeBSD 基础系统。新二进制下载并通过 `--help` 运行验证后，才停止旧服务并原子替换二进制。下载或运行验证失败时保留旧服务。

## 服务、目录与日志

| 项目 | 默认值 |
| --- | --- |
| rc.d 服务名 | `komari_agent` |
| 开机启动 | `komari_agent_enable=YES` |
| 服务脚本 | `/usr/local/etc/rc.d/komari_agent` |
| Agent 二进制 | `/usr/local/libexec/komari-agent/agent` |
| 启动参数 | `/usr/local/libexec/komari-agent/run-agent`，权限 `0700` |
| 工作目录与流量数据 | `/var/db/komari_agent`，目录权限 `0700` |
| 流量历史 | `/var/db/komari_agent/net_static.json` |
| 监督进程 PID | `/var/run/komari_agent.pid` |
| Agent PID | `/var/run/komari_agent.child.pid` |
| 日志 | syslog，tag `komari_agent`；默认查看 `/var/log/messages` |

```sh
service komari_agent status
service komari_agent restart
service komari_agent stop
service komari_agent start
sysrc komari_agent_enable
# 查询日志；不要把包含真实 token 的完整日志公开。
grep komari_agent /var/log/messages | tail -n 50
```

监督进程使用 `daemon -R 5`，Agent 意外退出或自动更新后会在 5 秒后重启。`service ... stop` 停止监督进程，并把 SIGTERM 转给 Agent，使其保存流量数据；不会把显式停止误当成需要重启的崩溃。重新运行安装器会保留数据目录。

自定义服务名可使用 `--install-service-name`，其中的连字符在 FreeBSD rc 服务名中转换成下划线。自定义二进制目录可使用 `--install-dir`，必须是绝对路径。FreeBSD 当前只支持 root 系统服务安装。

## Agent 行为与限制

- CPU、内存、Swap、负载、磁盘、网卡吞吐和运行时间使用 gopsutil 的 FreeBSD 实现；不要求挂载 Linux `/proc`。
- TCP/UDP 计数使用系统自带的 `/usr/bin/sockstat`，无需安装 `lsof`。
- `--month-rotate 1` 从本 Agent 开始记录的数据中计算每月累计流量；无法还原安装前未采集的历史。
- `--disable-web-ssh` 保持启用：远程命令、远程文件操作与远程版本切换均受现有开关约束；后台自动更新是独立机制。
- Linux 的 `--memory-exclude-bcf` 模式在 FreeBSD 回退到原生内存统计，不会因缺少 `/proc/meminfo` 上报零内存。
- FreeBSD 不会执行 Linux MOTD 修改逻辑；无需额外传入 `--disable-motd`。
- GPU 名称识别和详细 GPU 指标受平台及设备支持限制，不作为基础监控可用性的前提。
- ZFS 默认按 pool 去重。若希望只上报某个挂载点，可按需添加 `--include-mountpoint /`；不要未经确认就覆盖现有统计范围。

## 验证方式

本地通用检查：

```sh
sh -n install.sh
bash -n install.sh
go test -short -race ./...
go vet ./...
GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build ./...
```

FreeBSD 原生测试（`-short` 跳过既有外网 ping 集成测试）：

```sh
go test -short ./...
```

`monitoring/unit/net_freebsd_test.go` 在 FreeBSD 上检查实际 CPU、内存、磁盘、网卡、运行时间，以及打开本地 TCP/UDP socket 后的连接计数；同时验证 sockstat 输出解析。交叉编译通过不等同于原生运行通过。
