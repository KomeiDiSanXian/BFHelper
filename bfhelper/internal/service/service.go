// Package service 业务逻辑代码
package service

import (
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/dao"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/logger"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// Service 业务
type Service struct {
	zctx *zero.Ctx
	dao  *dao.Dao
}

// New 新建业务
func New(zctx *zero.Ctx) *Service {
	svc := Service{zctx: zctx}
	svc.dao = dao.New(global.DB)
	return &svc
}

// Log 日志记录
//
// 多业务使用同一个log
func (s *Service) Log() *logger.Logger {
	return global.Logger
}
