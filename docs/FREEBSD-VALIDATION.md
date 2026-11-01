# FreeBSD 15.1 验证记录

日期：2026-11-02。

## 原生环境

- 官方 `FreeBSD-15.1-RELEASE-amd64-ufs.qcow2.xz`，下载后核对官方 SHA-256：`25ad16e9d0c037f40a43c6f3c474b5c1ba62eaf9204dbd7f034fdbffe2c3dae5`。
- `uname`：`FreeBSD 15.1-RELEASE releng/15.1-n283562-96841ea08dcf GENERIC amd64`。
- QEMU TCG，2 vCPU、1 GiB RAM、UFS，网卡 `vtnet0`。
- Agent 和测试程序以 `GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0` 构建，在上述 FreeBSD 内执行，非仅交叉编译。

## 已通过

| 项目 | 结果 |
| --- | --- |
| `monitoring/unit` 原生短测试 | PASS，包含 CPU、原生及回退内存、磁盘、网卡、运行时间、实际 TCP/UDP socket |
| `monitoring/netstatic` 原生短测试 | PASS |
| `cmd` 原生短测试 | PASS |
| `update` 原生短测试 | PASS |
| `server` 原生短测试 | PASS，保留停用远程操作的既有保护测试 |
| `/bin/sh` 执行安装器 | PASS，使用 base `fetch`，自动识别 amd64 |
| `sysrc` 开机启用设置 | PASS，`komari_agent_enable=YES` |
| 服务 start/status/stop | PASS |
| Agent SIGTERM 后 supervisor 重启 | PASS，监督进程 PID 不变、Agent PID 改变 |
| 显式停止后的状态 | PASS，等待超过重启延时仍未运行，子进程 PID 文件删除 |
| 流量落盘 | PASS，`vtnet0` 实际记录 TX=147061、RX=8026602 bytes |
| Snapshot API 查询、版本解析、重装 | PASS，重装后保留既有网卡历史 |
| HTTP 404 下载失败 | PASS，旧二进制哈希及运行中 supervisor PID 均未改变 |
| 新二进制 `--help` 非零退出 | PASS，旧二进制哈希及运行中 supervisor PID 均未改变 |
| 参数与目录权限 | PASS，launcher 和数据目录均为 0700 |
| sh / Bash 参数转义 | PASS，包含空值、空格、单引号、命令替换文本及末尾换行 |
| Linux 通用短测试 + race detector | PASS |
| `go vet ./...`、sh/Bash 语法、git diff 格式 | PASS |

服务集成测试的关键输出：

```text
RESPAWN_PASS
STOP_AND_PERSISTENCE_PASS
SNAPSHOT_REINSTALL_PASS
DOWNLOAD_FAILURE_PRESERVES_SERVICE_PASS
INVALID_BINARY_PRESERVES_SERVICE_PASS
ALL_SERVICE_TESTS_PASS
```

## 验证范围

虚拟机的外部 DNS/SRV 查询受执行环境限制。安装与服务集成测试使用本机 HTTP 源提供从当前源码构建的 FreeBSD 二进制和模拟 GitHub Snapshot JSON；真实执行的是完整安装脚本、原生 fetch、rc.d/sysrc/daemon 和 Agent。未把这一结果写成“虚拟机直连 GitHub HTTPS / pkg 已通过”。标准 FreeBSD 安装路径不依赖 pkg。

未使用实际设备 Token 连接用户面板，未修改用户服务器。既有外网 ping 集成测试在 `-short` 模式下跳过。UFS 原生验证不代表完成 ZFS、多架构或 GPU 设备实机验证。

复验方法与安装指令见 [FREEBSD.md](FREEBSD.md)。
