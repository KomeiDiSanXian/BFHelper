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
├── zerobot/           # ZeroBot 插件适配器（即将移除）
├── main.go            # 独立入口（默认 Remilia）
├── main_zerobot.go    # ZeroBot 入口（build tag: zerobot）
└── config.yaml        # 全局配置
```

## 构建

### Remilia 模式（默认，推荐）

```bash
# 需要配置私有库（remilia + GoBattlefieldAPI）
go build -o bfhelper.exe .
```

### ZeroBot 模式 ⚠️ 即将废弃

```bash
go build -tags zerobot -o bfhelper.exe .
```

> **ZeroBot 支持将在 v3.0.0 中移除。**
> 建议现有 ZeroBot 用户迁移到 Remilia 模式。
> 如有迁移困难请联系开发者。

### 作为外部插件的导入方式

```go
// Remilia 项目
import "github.com/KomeiDiSanXian/BFHelper/remilia"
pm.Register(remilia.New())
```

## 配置

编辑 `config.yaml`：

```yaml
# 运行模式: remilia（默认）
mode: remilia

# EA 账号凭据（首次需手动填写）
account:
  sid: "your-ea-sid-cookie"
  remid: "your-ea-remid-cookie"
  game: "bf1"        # bf1 | bfv | bf4 | bf2042 | bf6

# Remilia 模式配置
remilia:
  ws_server: "ws://127.0.0.1:6700"
  ws_token: ""
  super_users: [123456]

# ZeroBot 模式配置（即将废弃）
zerobot:
  bot_names: ["蕾米"]
  command_prefix: "."
  super_users: [123456]
  ws_server: "ws://127.0.0.1:6700"
  ws_token: ""

bfeac:
  api_key: ""

blaze:
  enabled: false
```

## 功能

### 玩家命令

| 命令 | 说明 | 数据源 |
|------|------|--------|
| `/bf bind <name>` | 绑定玩家名到 QQ | 本地数据库 |
| `/bf unbind` | 解绑 | 本地数据库 |
| `/bf stats [name]` | 查询玩家战绩 | BTR (tracker.gg) |
| `/bf weapons [name]` | 武器数据 | Sparta Gateway |
| `/bf vehicles [name]` | 载具数据 | Sparta Gateway |
| `/bf recent [name]` | 最近战绩 | Bili22 API |
| `/bf exchange` | 本期交换皮肤 | Sparta Gateway |
| `/bf campaign` | 本期行动包 | Sparta Gateway |
| `/bf cb [name]` | 联ban查询 | BFEAC + BFBan |

### 服务器管理

| 命令 | 说明 | 权限 |
|------|------|------|
| `/bf group create` | 创建群组 | 群主 |
| `/bf group delete` | 删除群组 | 群主 |
| `/bf group owner <qq>` | 更换服主 | 服主 |
| `/bf group bind <gid...>` | 绑定服务器 | 超级用户 |
| `/bf group unbind <gid>` | 解绑服务器 | 服主 |
| `/bf group admin add <qq...>` | 添加管理员 | 服主 |
| `/bf group admin rm <qq>` | 删除管理员 | 服主 |
| `/bf admin kick <name>` | 踢出玩家 | 管理员 |
| `/bf admin ban [a] <name>` | 封禁 | 管理员 |
| `/bf admin unban [a] <name>` | 解封 | 管理员 |
| `/bf admin cm <a> [idx]` | 切换地图 | 管理员 |
| `/bf admin qm <a>` | 查看地图池 | 所有人 |

### 认证

| 命令 | 说明 |
|------|------|
| `/bf login` | 从配置登录 EA 账号 |
| `/bf logout` | 登出 |
| `/bf status` | 查看状态 |

## 依赖

| 包 | 用途 | 是否需要私有库 |
|------|------|--------|
| `bfhelper/core` | 类型 + 存储 + 错误定义 | ❌ |
| `bfhelper/public` | BTR/Bili22/BFEAC/BFBan API | ❌ |
| `bfhelper/battlefield` | EA 认证 + Sparta 管理 | ✅ `GoBattlefieldAPI` |
| `remilia/` | Remilia 框架适配器 | ✅ `remilia` |
| `zerobot/` | ZeroBot 适配器（即将移除） | ❌ |

## 免责声明

> 任何的查询不会导致被查询用户被 EA 封禁。
>
> 本插件仅是一个工具，开发者仅提供这样的一个工具，并不是您攻击谩骂的对象。
>
> 如果您遭到了联合封禁，本插件不提供任何有关 BFEAC 及 BFBan 解封的实际帮助。
> 请前往 [BFEAC申诉](https://bfeac.com/#/about) 和 [BFBan申诉](mailto:ban-appeals@bfban.com)。

## License

[AGPL-3.0](LICENSE)
