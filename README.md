<div align="center">
  <br>

  # BFHelper

  Battlefield 游戏辅助机器人 — 玩家战绩查询、服务器管理、联ban查询
</div>

<p align="center">
    <a href="https://goreportcard.com/report/github.com/KomeiDiSanXian/BFHelper">
        <img src="https://goreportcard.com/badge/github.com/KomeiDiSanXian/BFHelper">
    </a>
    <a href="LICENSE">
        <img alt="GitHub" src="https://img.shields.io/github/license/KomeiDiSanXian/BFHelper">
    </a>
</p>

---

## 架构

```
BFHelper/
├── bfhelper/          # 核心库（框架无关）
│   ├── core/          # 类型定义、错误哨兵、存储层
│   ├── public/        # BTR、Bili22、反作弊查询（零私有库依赖）
│   └── battlefield/   # EA 认证 + Sparta 服务器管理（依赖 BattlefieldAPI）
├── remilia/           # Remilia 框架插件适配器
├── main.go            # 独立入口（Remilia）
└── config.yaml        # 全局配置
```

## 构建

```bash
# 需要配置私有库（remilia + GoBattlefieldAPI）
go build -o bfhelper.exe .
```

### 作为外部插件的导入方式

```go
// Remilia 项目
import "github.com/KomeiDiSanXian/BFHelper/remilia"
pm.Register(remilia.New())
```

## 配置

编辑 `config.yaml`：

```yaml
mode: remilia

# OneBot 正向 WebSocket 适配器
remilia:
  ws_server: "ws://127.0.0.1:6700"
  ws_token: ""
  super_users: [123456]

# 插件专属配置（remilia 框架约定: plugins.<插件名>）
plugins:
  bfhelper:
    # 机器人超级管理员 QQ（可执行登录/登出，并绕过群内权限检查）
    super_users: [123456]

    # EA 账号凭据（首次需手动填写，从浏览器 cookie 中提取）
    account:
      sid: "your-ea-sid-cookie"
      remid: "your-ea-remid-cookie"
      game: "bf1"        # bf1 | bfv | bf4 | bf2042 | bf6
    bfeac:
      api_key: ""
    blaze:
      enabled: false
```

### 获取 EA 账号凭据

服务器管理接口依赖 EA 账号，需要从浏览器提取 `sid` 和 `remid` 两个 cookie：

1. 打开浏览器并登录 <https://www.ea.com/>
2. 按 `F12` 打开开发者工具 → **Application** → **Cookies** → `accounts.ea.com`
3. 复制 `sid` 和 `remid` 的值，填入上方 `config.yaml` 的 `plugins.bfhelper.account`

> `sid` 约 1 小时过期，`remid` 可维持数天。BFHelper 登录成功后会持久化凭据并自动续期，日常使用无需反复更新。

## 功能

### 玩家命令

| 命令 | 说明 | 数据源 |
|------|------|--------|
| `/bind <name>` | 绑定玩家名到 QQ | 本地数据库 |
| `/unbind` | 解绑 | 本地数据库 |
| `/stats [name]` | 查询玩家战绩 | BTR (tracker.gg) |
| `/weapons [name]` | 武器数据 | Sparta Gateway |
| `/vehicles [name]` | 载具数据 | Sparta Gateway |
| `/recent [name]` | 最近战绩 | Bili22 API |
| `/exchange` | 本期交换皮肤 | Sparta Gateway |
| `/campaign` | 本期行动包 | Sparta Gateway |
| `/cb [name]` | 联ban查询 | BFEAC + BFBan |

### 服务器管理

| 命令 | 说明 | 权限 |
|------|------|------|
| `/group create` | 创建群组 | 群主 |
| `/group delete` | 删除群组 | 服主 |
| `/group owner <qq>` | 更换服主 | 服主 |
| `/group bind <gid...>` | 绑定服务器 | 服主 |
| `/group unbind <gid>` | 解绑服务器 | 服主 |
| `/group admin add <qq...>` | 添加群组管理员 | 服主 |
| `/group admin rm <qq>` | 删除群组管理员 | 服主 |
| `/admin kick <name>` | 踢出玩家 | 管理员 |
| `/admin ban [a] <name>` | 封禁 | 管理员 |
| `/admin unban [a] <name>` | 解封 | 管理员 |
| `/admin banall <name>` | 全服封禁 | 管理员 |
| `/admin unbanall <name>` | 全服解封 | 管理员 |
| `/admin cm <a> [idx]` | 切换地图 | 管理员 |
| `/admin qm <a>` | 查看地图池 | 管理员 |
| `/admin move <p> <1\|2> [a]` | 移动玩家到队伍 | 管理员 |
| `/admin sadmin add\|rm <p> [a]` | 添加/移除服务器管理员 | 管理员 |
| `/admin vip add\|rm <p> [a]` | 添加/移除 VIP | 管理员 |
| `/admin logs <a> [n]` | 查询服务器操作日志 | 管理员 |

### 操作日志

机器人会记录群内全部管理操作（踢人、封禁、换图、移动玩家、VIP/管理员变更、群组变更等），
支持按日期/时间范围查询：

```text
/logs                        # 今天的操作日志
/logs 2026-08-01             # 指定日期起
/logs 2026-08-01 2026-08-27  # 日期范围
/logs 2026-08-01 2026-08-27 15:30   # 精确到分钟
/logs 2026-08-01 2026-08-27 ban     # 按操作类型过滤
```

操作类型包括：`kick` / `ban` / `unban` / `banall` / `unbanall` / `cm` / `move` /
`sadmin_add` / `sadmin_rm` / `vip_add` / `vip_rm` / `group_create` 等。

### 认证

| 命令 | 说明 | 权限 |
|------|------|------|
| `/login` | 从配置登录 EA 账号 | 超级管理员 |
| `/logout` | 登出 | 超级管理员 |
| `/status` | 查看状态 | 所有人 |

## 依赖

| 包 | 用途 | 是否需要私有库 |
|------|------|--------|
| `bfhelper/core` | 类型 + 存储 + 错误定义 | ❌ |
| `bfhelper/public` | BTR/Bili22/BFEAC/BFBan API | ❌ |
| `bfhelper/battlefield` | EA 认证 + Sparta 管理 | ✅ `GoBattlefieldAPI` |


## 免责声明

> 任何的查询不会导致被查询用户被 EA 封禁。
>
> 本插件仅是一个工具，开发者仅提供这样的一个工具，并不是您攻击谩骂的对象。
>
> 如果您遭到了联合封禁，本插件不提供任何有关 BFEAC 及 BFBan 解封的实际帮助。
> 请前往 [BFEAC申诉](https://bfeac.com/#/about) 和 [BFBan申诉](mailto:ban-appeals@bfban.com)。

## License

[AGPL-3.0](LICENSE)
