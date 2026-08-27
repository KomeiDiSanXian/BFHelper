package remilia

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	eventctx "github.com/KomeiDiSanXian/remilia/core/context"
	"github.com/KomeiDiSanXian/remilia/platform"
	"github.com/tidwall/gjson"
)

// handleAdminLogs 查询服务器 RSP 操作日志。
// 用法: /admin logs <别名> [数量]
func (p *Plugin) handleAdminLogs(ctx *eventctx.Context) error {
	if !p.requireLoggedIn(ctx) {
		return nil
	}
	gid, ok := p.requireGroupExists(ctx)
	if !ok {
		return nil
	}
	if !p.checkAdmin(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 3 {
		ctx.Reply(platform.TextMessage("用法: /admin logs <别名> [数量]"))
		return nil
	}
	alias := args[2]
	limit := 20
	if len(args) >= 4 {
		if n, err := strconv.Atoi(args[3]); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > 50 {
		limit = 50
	}

	sv, err := p.GetServerByAlias(gid, alias)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("未找到别名为 %s 的服务器", alias)))
		return nil
	}

	result, err := p.GetLogs(sv.GameID, 0, limit)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询日志失败: %v", err)))
		return nil
	}

	entries := extractLogEntries(result)
	if len(entries) == 0 {
		raw := result.Raw
		if len(raw) > 800 {
			raw = raw[:800] + "..."
		}
		ctx.Reply(platform.TextMessage(fmt.Sprintf("%s 暂无操作日志或响应格式未知:\n%s", alias, raw)))
		return nil
	}

	var msg strings.Builder
	fmt.Fprintf(&msg, "%s 的操作日志:\n", alias)
	for i, e := range entries {
		if i >= limit {
			break
		}
		msg.WriteString(formatLogEntry(e))
		msg.WriteString("\n")
	}
	ctx.Reply(platform.TextMessage(msg.String()))
	return nil
}

// logEntry RSP 日志条目（字段解析结果）。
type logEntry struct {
	Time    string
	Type    string
	Message string
	Player  string
}

// extractLogEntries 尝试从 RSP.getLogs 响应中提取日志条目列表。
func extractLogEntries(result gjson.Result) []logEntry {
	var arr gjson.Result
	switch {
	case result.IsArray():
		arr = result
	case result.Get("log").IsArray():
		arr = result.Get("log")
	case result.Get("logs").IsArray():
		arr = result.Get("logs")
	case result.Get("logEntries").IsArray():
		arr = result.Get("logEntries")
	case result.Get("entries").IsArray():
		arr = result.Get("entries")
	default:
		return nil
	}

	items := arr.Array()
	entries := make([]logEntry, 0, len(items))
	for _, it := range items {
		e := logEntry{
			Time:    firstString(it, "timestamp", "time", "logTime", "dateTime", "createdAt", "createdDate"),
			Type:    firstString(it, "logType", "type", "action"),
			Message: firstString(it, "message", "msg", "description"),
			Player:  firstString(it, "playerName", "personaName", "player", "target"),
		}
		if t := it.Get("logTime"); t.Exists() {
			e.Time = formatLogTime(t)
		}
		if t := it.Get("timestamp"); t.Exists() && e.Time == "" {
			e.Time = formatLogTime(t)
		}
		if e.Time == "" || e.Type == "" {
			// 未知字段结构时保留原始 JSON，便于排查。
			raw := it.Raw
			if len(raw) > 200 {
				raw = raw[:200] + "..."
			}
			e.Message = raw
		}
		entries = append(entries, e)
	}
	return entries
}

// firstString 返回第一个存在的字符串字段值。
func firstString(r gjson.Result, paths ...string) string {
	for _, path := range paths {
		if v := r.Get(path); v.Exists() {
			return v.String()
		}
	}
	return ""
}

// formatLogTime 将时间字段统一格式化为 "2006-01-02 15:04"。
func formatLogTime(v gjson.Result) string {
	switch v.Type {
	case gjson.Number:
		ms := v.Int()
		if ms == 0 {
			return ""
		}
		if ms < 1e12 { // 秒级时间戳
			return time.Unix(ms, 0).Format("2006-01-02 15:04")
		}
		return time.UnixMilli(ms).Format("2006-01-02 15:04")
	case gjson.String:
		s := strings.TrimSpace(v.String())
		if s == "" {
			return ""
		}
		for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02 15:04", "2006-01-02"} {
			if t, err := time.Parse(layout, s); err == nil {
				return t.Format("2006-01-02 15:04")
			}
		}
		return s
	default:
		return ""
	}
}

// formatLogEntry 渲染单条日志。
func formatLogEntry(e logEntry) string {
	var b strings.Builder
	if e.Time != "" {
		b.WriteString("[" + e.Time + "] ")
	}
	if e.Type != "" {
		b.WriteString(e.Type + " ")
	}
	if e.Player != "" {
		b.WriteString(e.Player + " ")
	}
	if e.Message != "" {
		b.WriteString(e.Message)
	}
	s := strings.TrimSpace(b.String())
	if s == "" {
		return "(空条目)"
	}
	return s
}

// parseLogQuery 解析 /logs 的时间范围与操作类型过滤参数。
//
// 支持的输入：
//
//	/logs                      今天
//	/logs 2026-08-01           指定日期
//	/logs 2026-08-01 2026-08-27
//	/logs 2026-08-01 2026-08-27 15:30  时间范围到分钟
//	/logs 2026-08-01 2026-08-27 ban   按操作类型过滤
func parseLogQuery(args []string) (from, to *time.Time, action string) {
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var times []time.Time
	for _, arg := range args {
		if t, ok := parseLogTime(arg); ok {
			times = append(times, t)
			continue
		}
		if action == "" && isActionFilter(arg) {
			action = arg
		}
	}

	switch len(times) {
	case 0:
		return &dayStart, &dayEnd, action
	case 1:
		t := times[0]
		return &t, &dayEnd, action
	default:
		from, to = &times[0], &times[1]
		return from, to, action
	}
}

// parseLogTime 解析日期或日期时间字符串。
func parseLogTime(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	for _, layout := range []string{"2006-01-02 15:04", "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// isActionFilter 判断是否为已知的操作类型过滤词。
func isActionFilter(s string) bool {
	switch strings.ToLower(s) {
	case "kick", "ban", "unban", "banall", "unbanall", "cm", "move",
		"sadmin_add", "sadmin_rm", "vip_add", "vip_rm", "group_create",
		"group_delete", "group_owner", "group_bind", "group_unbind",
		"group_admin_add", "group_admin_rm":
		return true
	}
	return false
}

// handleLogs 查询本群操作日志（支持日期/时间范围过滤）。
// 用法: /logs [起始] [结束] [操作类型]
func (p *Plugin) handleLogs(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid, ok := p.requireGroupExists(ctx)
	if !ok {
		return nil
	}
	if !p.checkAdmin(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	var rest []string
	if len(args) > 1 {
		rest = args[1:]
	}
	from, to, action := parseLogQuery(rest)

	logs, err := p.QueryOperationLogs(gid, from, to, action, 20)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询日志失败: %v", err)))
		return nil
	}

	actionText := "全部操作"
	if action != "" {
		actionText = action
	}
	if len(logs) == 0 {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("在 %s ~ %s 内没有 %s 记录", from.Format("2006-01-02 15:04"), to.Format("2006-01-02 15:04"), actionText)))
		return nil
	}

	var msg strings.Builder
	fmt.Fprintf(&msg, "📋 操作日志（%s）%s ~ %s 显示 %d 条（最多 20）:\n",
		actionText, from.Format("2006-01-02 15:04"), to.Format("2006-01-02 15:04"), len(logs))
	for _, log := range logs {
		icon := "✅"
		if !log.Success {
			icon = "❌"
		}
		line := fmt.Sprintf("%s %s %s[%d] %s", log.CreatedAt.Format("01-02 15:04"), icon, log.Action, log.OperatorQQ, log.Target)
		if log.ServerName != "" {
			line += " @" + log.ServerName
		}
		if log.Detail != "" {
			line += " " + log.Detail
		}
		msg.WriteString(line + "\n")
	}
	ctx.Reply(platform.TextMessage(msg.String()))
	return nil
}
