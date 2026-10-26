# komari-agent

StarHanami fork of [komari-monitor/komari-agent](https://github.com/komari-monitor/komari-agent).
本分支的安装脚本和二进制自动更新均使用 `StarHanami/komari-agent` 的 Releases，容器镜像使用 `ghcr.io/starhanami/komari-agent`。

## 配置方式

agent 参数可以通过命令行参数、环境变量或 JSON 配置文件传入。

最小启动示例：

```bash
./komari-agent --endpoint "https://example.com" --token "your-token"
```

使用环境变量：

```bash
export AGENT_ENDPOINT="https://example.com"
export AGENT_TOKEN="your-token"
./komari-agent
```

使用 JSON 配置文件：

```bash
./komari-agent --config ./config.json
```

`config.json` 示例：

```json
{
  "endpoint": "https://example.com",
  "token": "your-token",
  "interval": 3,
  "disable_auto_update": false,
  "disable_web_ssh": false,
  "disable_motd": false,
  "ignore_unsafe_cert": false
}
```

配置优先级从低到高为：默认值、命令行参数、环境变量、JSON 配置文件。

常用配置项：

表中支持版本表示该参数本身首次在发布 tag 中出现；环境变量和 JSON 配置文件方式从 `1.1.33` 起支持，早于最早 tag 的参数记为 `0.0.9`。

| JSON 字段 | 环境变量 | 命令行参数 | 说明 | 支持版本 |
| --- | --- | --- | --- | --- |
| `endpoint` | `AGENT_ENDPOINT` | `--endpoint`, `-e` | 面板地址 | `0.0.9` |
| `token` | `AGENT_TOKEN` | `--token`, `-t` | agent token | `0.0.9` |
| `interval` | `AGENT_INTERVAL` | `--interval`, `-i` | 数据采集间隔，单位秒 | `0.0.9` |
| `disable_auto_update` | `AGENT_DISABLE_AUTO_UPDATE` | `--disable-auto-update` | 禁用自动更新 | `0.0.9` |
| `disable_motd` | `AGENT_DISABLE_MOTD` | `--disable-motd` | 禁止修改 Linux MOTD，默认 `false`（仍修改） | StarHanami fork |
| `disable_web_ssh` | `AGENT_DISABLE_WEB_SSH` | `--disable-web-ssh` | 禁用远程控制 | `0.0.9` |
| `ignore_unsafe_cert` | `AGENT_IGNORE_UNSAFE_CERT` | `--ignore-unsafe-cert`, `-u` | 忽略不安全证书 | `0.0.9` |
| `include_nics` | `AGENT_INCLUDE_NICS` | `--include-nics` | 仅统计指定网卡，逗号分隔 | `0.0.22` |
| `exclude_nics` | `AGENT_EXCLUDE_NICS` | `--exclude-nics` | 排除指定网卡，逗号分隔 | `0.0.22` |
| `include_mountpoints` | `AGENT_INCLUDE_MOUNTPOINTS` | `--include-mountpoint` | 仅统计指定挂载点，分号分隔 | `0.1.0` |
| `month_rotate` | `AGENT_MONTH_ROTATE` | `--month-rotate` | 流量统计每月重置日期，`0` 为禁用 | `0.1.0` |
| `auto_discovery_key` | `AGENT_AUTO_DISCOVERY_KEY` | `--auto-discovery` | 自动发现密钥 | `1.0.40` |
| `custom_dns` | `AGENT_CUSTOM_DNS` | `--custom-dns` | 自定义 DNS 服务器 | `1.0.80` |
| `enable_gpu` | `AGENT_ENABLE_GPU` | `--gpu` | 启用详细 GPU 监控 | `1.0.80` |
| `disable_compression` | `AGENT_DISABLE_COMPRESSION` | `--disable-compression` | 禁用 v2 传输压缩 | `1.2.10` |
| `prefer_ip_version` | `AGENT_PREFER_IP_VERSION` | `--prefer-ip-version` | 优先使用 IP 版本，可选 `4` 或 `6` | 未发布 |

完整参数可运行：

```bash
./komari-agent --help
```

详见 `cmd/flags/flags.go` 及 `cmd/root.go`

## 保留自己的 MOTD（Linux）

默认行为与上游一致：启用远程控制时追加 MOTD 警告，正常退出时恢复。要完全跳过 MOTD 修改，在已有启动参数后加上：

```bash
./komari-agent --endpoint "https://example.com" --token "your-token" --disable-motd
```

也可使用 `AGENT_DISABLE_MOTD=true`（或 `1`），或在 JSON 配置中设置 `"disable_motd": true`。此选项不关闭 Web SSH / 远程命令，也不影响 Windows 安全通知。

启用后，Agent 不会创建、写入、删除或恢复 `/etc/motd`，也不会清理 `/etc/update-motd.d/99-komari-agent-warning`。旧实例正常退出时仍可能执行其原有的恢复逻辑，因此升级时请先正常停止旧服务，再启用本选项并启动新版本。之前崩溃留下的警告或旧版 hook 需要自行检查处理；本选项不会自动清理它们。

已有 systemd 服务可通过 `sudo systemctl edit komari-agent` 添加：

```ini
[Service]
Environment=AGENT_DISABLE_MOTD=true
```

然后运行 `sudo systemctl daemon-reload` 和 `sudo systemctl restart komari-agent`。如果已有 JSON 配置显式设置 `disable_motd`，请同时修改它，因为 JSON 优先级最高。用户级服务改用 `systemctl --user`。

从本分支安装开发快照（安装脚本会将 `--disable-motd` 写入服务启动参数）：

```bash
curl -fsSL https://raw.githubusercontent.com/StarHanami/komari-agent/main/install.sh -o install.sh
bash install.sh --install-version snapshot --endpoint "https://example.com" --token "your-token" --disable-motd
```

快照通过本仓库的 `Snapshot-*` 预发布更新。未指定 `--install-version` 时，安装器只选择本仓库的稳定 Release；若尚未发布稳定版，请明确使用 `snapshot`。Windows 安装器也使用本分支的 Releases。

Go module/import 路径保留上游名称以便持续合并上游；它不决定更新来源。默认更新仓库在 `update.Repo` 中设置为 `StarHanami/komari-agent`，所有本地和 CI 构建均使用这一默认值。现有 CI 的 Release 发布目标、快照发布目标和容器镜像名均通过 `github.repository` 解析，因此在本 fork 中会发布至 StarHanami。

本地验证：`go test -short -race ./...`、`go vet ./...`、`sh -n install.sh`。`-short` 只跳过三个既有的外网 ping 集成测试；这些测试需要外部 DNS、IPv6 和 ICMP 权限，可在合适的主机上通过 `go test ./server -run 'Test(ICMP|TCP|HTTP)Ping'` 单独执行。
