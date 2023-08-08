<div align="center">
  <br>

  # BFHelper

  BFHelper是依赖于 [ZeroBot](https://github.com/wdvxdr1123/ZeroBot) 的插件
</div>

<p align="center">
    <a href=""></a>
    <a href="https://goreportcard.com/report/github.com/KomeiDiSanXian/BFHelper">
        <img src="https://goreportcard.com/badge/github.com/KomeiDiSanXian/BFHelper">
    </a> 
    <a href="https://github.com/wdvxdr1123/ZeroBot">
        <img src="https://img.shields.io/badge/zerobot-v1.7.4-black?style=flat-square&logo=go">
    </a>
    <a href="https://raw.githubusercontent.com/KomeiDiSanXian/BFHelper/master/LICENSE">
        <img alt="GitHub" src="https://img.shields.io/github/license/KomeiDiSanXian/BFHelper">
    </a>
    <a href="https://pkg.go.dev/github.com/KomeiDiSanXian/BFHelper">
        <img src="https://pkg.go.dev/badge/github.com/KomeiDiSanXian/BFHelper.svg" alt="Go Reference">
    </a>
</p>

<p align="center">
    <a href="https://www.ea.com/games/battlefield/battlefield-1">
        <img src="https://img.shields.io/badge/BattleField-1-yellow?logo=EA&logoColor=red">
    </a> 
    <a href="https://www.ea.com/games/battlefield/battlefield-5">
        <img src="https://img.shields.io/badge/BattleField-V-blue?logo=EA&logoColor=red">
    </a> 
</p>

---

## 声明
> 任何的查询不会导致被查询用户被EA封禁

> 请注意，任何开发者都**没有义务**回答您的问题
> 
> 本插件仅是一个工具，开发者只是提供了这样的一个工具，并不是您攻击谩骂的对象

开发者**不**负责解封

开发者**不**负责解封

开发者**不**负责解封

如果您遭到了联合封禁，本插件不提供任何有关BFEAC及BFBan解封的实际帮助

请前往[BFEAC申诉](https://bfeac.com/#/about)和[BFBan申诉](mailto:ban-appeals@bfban.com)

如果您因本插件未来添加的功能导致被某一服务器添加进该服Ban列，请进入该服务器的QQ群，联系管理员解封

---

## 功能

> 前往 [src](https://github.com/KomeiDiSanXian/BFHelper/tree/master/bfhelper) 查看更多

- [ ] 举报作弊行为
- [x] 查询玩家是否被联合封禁 (在 [BFEAC](https://bfeac.com/#/) & [BFban](https://bfban.gametools.network/) 中查询)

### 战地一 BattleField 1
#### 玩家
- [x] 查询玩家战绩 (通过 [BTR](https://battlefieldtracker.com/) 实现)
- [x] 查询玩家武器信息 (Battlefield Gateway 实现)
- [x] 查询玩家载具击杀信息 (Battlefield Gateway 实现)
- [x] 查询玩家最近游玩信息 (借助 @Bili22 的api实现)
- [x] 查询本期的交换信息 (Battlefield Gateway 实现)
- [x] 查询本期行动包 (Battlefield Gateway 实现)

#### 服务器
- [ ] 踢出玩家
- [ ] 封禁&解封玩家
- [ ] 切换地图
- [ ] 添加&删除VIP
- [ ] 修改服务器配置

### 战地五 BattleField V
> 目前计划支持

### 战地2042 BattleField 2042
> 目前没有计划

---

## 如何使用
> 玩家相关命令
> 
> 中括号内可填可不填
- [x] **.绑定 xxx** 用于将xxx绑定到发送该命令的用户, 便于以后的查询
- [x] **.战绩 [xxx]** 用于查询xxx的战绩, 如果没有xxx将会查询发送该命令用户的战绩
- [x] **.武器 [xxx]** 类似上者, 改为查询武器击杀数据
- [x] **.最近 [xxx]** 类似上者, 改为查询最近战绩
- [x] **.载具 [xxx]** 类似上者, 改为查询载具击杀数据
- [x] **.交换** 查询战地一本期交换皮肤
- [x] **.行动** 查询战地一本期行动包信息
- [x] **.cb [xxx]** 查询xxx联合封禁信息

> 战地一服务器相关命令
>
> 中括号内可填可不填. 服主权限于群主权限等同; 服管理员于群管理员权限等同
- [x] **.创建服务器群组 [qq号]** 让群聊可以开始绑定服务器, 这些服务器的服主是所填的qq, 不填则为发送人 **`需要群主及以上权限`**
- [x] **.删除服务器群组** 删除群聊所有的服务器绑定信息 **`需要群主及以上权限`**
- [x] **.更换服主 qq号** 更换群聊绑定的群组的所属人 **`需要服主及以上权限`**
- [x] **.绑定服务器 群号 gameid1 gameid2...**  将gameid1, gameid2...的服务器绑定到群号 **`需要超级管理员权限`**
- [x] **.添加管理 qq1 qq2...** 将qq1, qq2...添加为服务器群组的管理员 **`需要服主及以上权限`**
- [x] **.设置别名 gameid 别名** 设置gameid的服务器别名为 别名 **`需要服管理员及以上权限`**
- [x] **.解绑服务器 gameid** 解除为gameid的服务器与群聊的绑定 **`需要服主及以上权限`**
- [x] **.删除管理 qq** 解除qq的服管理员权限 **`需要服主及以上权限`**
- [ ] **.kick [别名] 玩家** 在 别名 的服务器踢出 玩家 **`需要服管理员及以上权限`**
- [ ] **.ban [别名] 玩家** 在 别名 的服务器将 玩家 封禁. 不填别名则在所有已经绑定的服务器封禁 **`需要服管理员及以上权限`**
- [ ] **.unban [别名] 玩家** 在 别名 的服务器将 玩家 解封. 不填别名则在所有已经绑定的服务器解封 **`需要服管理员及以上权限`**
- [ ] **.changemap 别名 地图号** 将 别名 的服务器的地图切换到地图号 **`需要服管理员及以上权限`**

---

## 如何安装

> **注意**: release 中的插件有且仅有本插件 
>
> 对于Windows 系统，仅支持win 7 (win server 2008 R2) 及以上
### a. 下载二进制程序

1. 前往 [release](https://github.com/KomeiDiSanXian/BFHelper/releases) 或 [CI](https://github.com/KomeiDiSanXian/BFHelper/actions/workflows/go.yml) 下载符合您系统的版本
2. 启动应用程序
> **注意**: 第一次启动会生成配置文件 `botcongfig.yaml`，请修改该配置

3. 修改配置后，重新启动应用，同时启动 OneBot 框架 (如 [go-cqhttp](https://github.com/Mrs4s/go-cqhttp))
4. 修改 `data/battlefield` 中的 `settings.yml`
> **注意**: 如果没有该文件，请使用一次本插件, 插件将会生成一份 `settings.yml`
>
> 修改后无需重启

### b. 本地编译

1. 下载并安装最新的 [golang](https://studygolang.com/dl) 环境
2. clone [FloatTech/ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin)
3. 编辑`main.go`文件中的import, 在其中添加

```go
_ "github.com/KomeiDiSanXian/BFHelper/bfhelper"
```
4. 下载本项目中的data文件夹，复制进 `ZeroBot-Plugin` 并对其中的 `data/battlefield/settings.yml` 按需编辑
5. 根据你所使用的平台进行编译
6. 运行 OneBot 框架 然后运行你编译的文件

### c. 使用RemiliaBot
> [RemiliaBot](https://github.com/KomeiDiSanXian/RemiliaBot) 是 [FloatTech/ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin) 的 fork 分支

1. 下载 [RemiliaBot](https://github.com/KomeiDiSanXian/RemiliaBot/releases)
2. 编辑其中的 `data/battlefield/settings.yml`
3. 编辑 RemiliaBot (参考 RemiliaBot 的 [README.md](https://github.com/KomeiDiSanXian/RemiliaBot/blob/master/README.md))
4. 启动 OneBot 框架和 RemiliaBot

----
## 特别感谢
- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
- [ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin)
- [Bili22](mailto:b22lengfeng@qq.com)
- [SakuraKooi](https://github.com/SakuraKoi)
- [GameTools](https://github.com/Community-network)
