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

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/GoSimplicity/AI-CloudOps/internal/model"
	"go.uber.org/zap"

	"github.com/GoSimplicity/AI-CloudOps/internal/workorder/dao"
)

type FormDesignService interface {
	CreateFormDesign(ctx context.Context, formDesignReq *model.CreateWorkorderFormDesignReq) error
	UpdateFormDesign(ctx context.Context, formDesignReq *model.UpdateWorkorderFormDesignReq) error
	DeleteFormDesign(ctx context.Context, id int) error
	GetFormDesign(ctx context.Context, id int) (*model.WorkorderFormDesign, error)
	ListFormDesign(ctx context.Context, req *model.ListWorkorderFormDesignReq) (*model.ListResp[*model.WorkorderFormDesign], error)
}

type formDesignService struct {
	dao         dao.WorkorderFormDesignDAO
	categoryDao dao.WorkorderCategoryDAO
	logger      *zap.Logger
}

func NewFormDesignService(dao dao.WorkorderFormDesignDAO, categoryDao dao.WorkorderCategoryDAO, logger *zap.Logger) FormDesignService {
	return &formDesignService{
		dao:         dao,
		categoryDao: categoryDao,
		logger:      logger,
	}
}

func (s *formDesignService) CreateFormDesign(ctx context.Context, formDesignReq *model.CreateWorkorderFormDesignReq) error {
	exists, err := s.dao.CheckFormDesignNameExists(ctx, formDesignReq.Name)
	if err != nil {
		s.logger.Error("检查表单设计名称是否存在失败", zap.Error(err), zap.String("name", formDesignReq.Name))
		return fmt.Errorf("检查表单设计名称失败: %w", err)
	}
	if exists {
		s.logger.Warn("表单设计名称已存在", zap.String("name", formDesignReq.Name))
		return dao.ErrFormDesignNameExists
	}

	if len(formDesignReq.Schema.Fields) == 0 {
		s.logger.Error("表单结构不能为空")
		return errors.New("表单结构不能为空")
	}

	// 生成表单字段ID
	s.generateFieldIDs(&formDesignReq.Schema)

	if len(formDesignReq.Tags) > 0 {
		for _, tag := range formDesignReq.Tags {
			if strings.TrimSpace(tag) == "" {
				return errors.New("标签不能为空字符串")
			}
		}
	}

	schemaJSON, err := json.Marshal(formDesignReq.Schema)
	if err != nil {
		s.logger.Error("表单结构序列化失败", zap.Error(err))
		return fmt.Errorf("表单结构序列化失败: %w", err)
	}

	// 转换为JSONMap
	var schemaMap model.JSONMap
	err = json.Unmarshal(schemaJSON, &schemaMap)
	if err != nil {
		s.logger.Error("转换表单结构为JSONMap失败", zap.Error(err))
		return fmt.Errorf("转换表单结构失败: %w", err)
	}

	formDesign := &model.WorkorderFormDesign{
		Name:         formDesignReq.Name,
		Description:  formDesignReq.Description,
		Schema:       schemaMap,
		Status:       formDesignReq.Status,
		CategoryID:   formDesignReq.CategoryID,
		OperatorID:   formDesignReq.OperatorID,
		OperatorName: formDesignReq.OperatorName,
		Tags:         formDesignReq.Tags,
		IsTemplate:   formDesignReq.IsTemplate,
	}

	if err := s.dao.CreateFormDesign(ctx, formDesign); err != nil {
		s.logger.Error("创建表单设计失败", zap.Error(err), zap.String("name", formDesignReq.Name))
		return err
	}

	return nil
}

func (s *formDesignService) UpdateFormDesign(ctx context.Context, formDesignReq *model.UpdateWorkorderFormDesignReq) error {
	existingFormDesign, err := s.dao.GetFormDesignByName(ctx, formDesignReq.Name)
	if err != nil {
		if errors.Is(err, dao.ErrFormDesignNotFound) {
		} else {
			s.logger.Error("获取表单设计失败", zap.Error(err), zap.String("name", formDesignReq.Name))
			return fmt.Errorf("获取表单设计失败: %w", err)
		}
	}

	if existingFormDesign != nil && existingFormDesign.ID != formDesignReq.ID {
		s.logger.Warn("表单设计名称已存在", zap.String("name", formDesignReq.Name))
		return dao.ErrFormDesignNameExists
	}

	if len(formDesignReq.Schema.Fields) == 0 {
		s.logger.Error("表单结构不能为空")
		return errors.New("表单结构不能为空")
	}

	// 生成表单字段ID
	s.generateFieldIDs(&formDesignReq.Schema)

	if len(formDesignReq.Tags) > 0 {
		for _, tag := range formDesignReq.Tags {
			if strings.TrimSpace(tag) == "" {
				return errors.New("标签不能为空字符串")
			}
		}
	}

	schemaJSON, err := json.Marshal(formDesignReq.Schema)
	if err != nil {
		s.logger.Error("表单结构序列化失败", zap.Error(err))
		return fmt.Errorf("表单结构序列化失败: %w", err)
	}

	// 转换为JSONMap
	var schemaMap model.JSONMap
	err = json.Unmarshal(schemaJSON, &schemaMap)
	if err != nil {
		s.logger.Error("转换表单结构为JSONMap失败", zap.Error(err))
		return fmt.Errorf("转换表单结构失败: %w", err)
	}

	formDesign := &model.WorkorderFormDesign{
		Model:       model.Model{ID: formDesignReq.ID},
		Name:        formDesignReq.Name,
		Description: formDesignReq.Description,
		Schema:      schemaMap,
		CategoryID:  formDesignReq.CategoryID,
		Status:      formDesignReq.Status,
		Tags:        formDesignReq.Tags,
		IsTemplate:  formDesignReq.IsTemplate,
	}

	if err := s.dao.UpdateFormDesign(ctx, formDesign); err != nil {
		s.logger.Error("更新表单设计失败", zap.Error(err), zap.Int("id", formDesignReq.ID))
		return err
	}

	return nil
}

func (s *formDesignService) DeleteFormDesign(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("表单设计ID无效")
	}

	_, err := s.dao.GetFormDesign(ctx, id)
	if err != nil {
		s.logger.Error("获取表单设计失败", zap.Error(err), zap.Int("id", id))
		return err
	}

	if err := s.dao.DeleteFormDesign(ctx, id); err != nil {
		s.logger.Error("删除表单设计失败", zap.Error(err), zap.Int("id", id))
		return err
	}

	return nil
}

func (s *formDesignService) GetFormDesign(ctx context.Context, id int) (*model.WorkorderFormDesign, error) {
	if id <= 0 {
		return nil, errors.New("表单设计ID无效")
	}

	formDesign, err := s.dao.GetFormDesign(ctx, id)
	if err != nil {
		s.logger.Error("获取表单设计失败", zap.Error(err), zap.Int("id", id))
		return nil, err
	}

	return formDesign, nil
}

func (s *formDesignService) ListFormDesign(ctx context.Context, req *model.ListWorkorderFormDesignReq) (*model.ListResp[*model.WorkorderFormDesign], error) {
	formDesigns, total, err := s.dao.ListFormDesign(ctx, req)
	if err != nil {
		s.logger.Error("获取表单设计列表失败", zap.Error(err))
		return nil, err
	}

	return &model.ListResp[*model.WorkorderFormDesign]{
		Items: formDesigns,
		Total: total,
	}, nil
}

// generateFieldIDs 生成字段ID
func (s *formDesignService) generateFieldIDs(schema *model.FormSchema) {
	for i := range schema.Fields {
		schema.Fields[i].ID = strconv.Itoa(i + 1)
	}
}
