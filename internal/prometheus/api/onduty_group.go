/*
 * MIT License
 *
 * Copyright (c) 2024 Bamboo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package api

import (
	"github.com/GoSimplicity/AI-CloudOps/internal/model"
	"github.com/GoSimplicity/AI-CloudOps/internal/prometheus/service/alert"
	"github.com/GoSimplicity/AI-CloudOps/pkg/base"
	"github.com/GoSimplicity/AI-CloudOps/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type OnDutyGroupHandler struct {
	alertOnDutyService alert.AlertManagerOnDutyService
}

func NewOnDutyGroupHandler(alertOnDutyService alert.AlertManagerOnDutyService) *OnDutyGroupHandler {
	return &OnDutyGroupHandler{
		alertOnDutyService: alertOnDutyService,
	}
}

func (h *OnDutyGroupHandler) RegisterRouters(server *gin.Engine) {
	monitorGroup := server.Group("/api/monitor")
	{
		monitorGroup.GET("/onduty_groups/list", h.GetMonitorOnDutyGroupList)
		monitorGroup.POST("/onduty_groups/create", h.CreateMonitorOnDutyGroup)
		monitorGroup.POST("/onduty_groups/changes", h.CreateMonitorOnDutyGroupChange)
		monitorGroup.GET("/onduty_groups/changes/:id", h.GetMonitorOnDutyGroupChangeList)
		monitorGroup.PUT("/onduty_groups/update/:id", h.UpdateMonitorOnDutyGroup)
		monitorGroup.DELETE("/onduty_groups/delete/:id", h.DeleteMonitorOnDutyGroup)
		monitorGroup.GET("/onduty_groups/detail/:id", h.GetMonitorOnDutyGroup)
		monitorGroup.GET("/onduty_groups/future_plan/:id", h.GetMonitorOnDutyGroupFuturePlan)
		monitorGroup.GET("/onduty_groups/history/:id", h.GetMonitorOnDutyHistory)
	}
}

func (h *OnDutyGroupHandler) GetMonitorOnDutyGroupList(ctx *gin.Context) {
	var req model.GetMonitorOnDutyGroupListReq

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.alertOnDutyService.GetMonitorOnDutyGroupList(ctx.Request.Context(), &req)
	})
}

func (h *OnDutyGroupHandler) CreateMonitorOnDutyGroup(ctx *gin.Context) {
	var req model.CreateMonitorOnDutyGroupReq

	uc := ctx.MustGet("user").(jwt.UserClaims)
	req.UserID = uc.Uid
	req.CreateUserName = uc.Username

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.alertOnDutyService.CreateMonitorOnDutyGroup(ctx.Request.Context(), &req)
	})
}

func (h *OnDutyGroupHandler) CreateMonitorOnDutyGroupChange(ctx *gin.Context) {
	var req model.CreateMonitorOnDutyGroupChangeReq

	uc := ctx.MustGet("user").(jwt.UserClaims)
	req.UserID = uc.Uid
	req.CreateUserName = uc.Username

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.alertOnDutyService.CreateMonitorOnDutyGroupChange(ctx.Request.Context(), &req)
	})
}

func (h *OnDutyGroupHandler) UpdateMonitorOnDutyGroup(ctx *gin.Context) {
	var req model.UpdateMonitorOnDutyGroupReq

	id, err := base.GetParamID(ctx)
	if err != nil {
		base.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.ID = id

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.alertOnDutyService.UpdateMonitorOnDutyGroup(ctx.Request.Context(), &req)
	})
}

func (h *OnDutyGroupHandler) DeleteMonitorOnDutyGroup(ctx *gin.Context) {
	var req model.DeleteMonitorOnDutyGroupReq

	id, err := base.GetParamID(ctx)
	if err != nil {
		base.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.ID = id

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.alertOnDutyService.DeleteMonitorOnDutyGroup(ctx.Request.Context(), &req)
	})
}

func (h *OnDutyGroupHandler) GetMonitorOnDutyGroup(ctx *gin.Context) {
	var req model.GetMonitorOnDutyGroupReq

	id, err := base.GetParamID(ctx)
	if err != nil {
		base.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.ID = id

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.alertOnDutyService.GetMonitorOnDutyGroup(ctx.Request.Context(), &req)
	})
}

// GetMonitorOnDutyGroupFuturePlan 获取指定值班组的未来值班计划
func (h *OnDutyGroupHandler) GetMonitorOnDutyGroupFuturePlan(ctx *gin.Context) {
	var req model.GetMonitorOnDutyGroupFuturePlanReq

	id, err := base.GetParamID(ctx)
	if err != nil {
		base.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.ID = id

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.alertOnDutyService.GetMonitorOnDutyGroupFuturePlan(ctx.Request.Context(), &req)
	})
}

func (h *OnDutyGroupHandler) GetMonitorOnDutyHistory(ctx *gin.Context) {
	var req model.GetMonitorOnDutyHistoryReq

	id, err := base.GetParamID(ctx)
	if err != nil {
		base.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.OnDutyGroupID = id

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.alertOnDutyService.GetMonitorOnDutyHistory(ctx.Request.Context(), &req)
	})
}

func (h *OnDutyGroupHandler) GetMonitorOnDutyGroupChangeList(ctx *gin.Context) {
	var req model.GetMonitorOnDutyGroupChangeListReq

	id, err := base.GetParamID(ctx)
	if err != nil {
		base.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.OnDutyGroupID = id

	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.alertOnDutyService.GetMonitorOnDutyGroupChangeList(ctx.Request.Context(), &req)
	})
}
