package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcp-skill-hub/server/internal/skill"
)

// ListSkills 获取技能列表
func ListSkills(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		opts := skill.ListOptions{
			Page:     1,
			PageSize: 20,
			SortBy:   c.DefaultQuery("sort", "downloads"),
			Category: c.Query("category"),
			Query:    c.Query("q"),
		}

		if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
			opts.Page = page
		}
		if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil {
			opts.PageSize = pageSize
		}

		// 标签过滤
		if tags := c.QueryArray("tags"); len(tags) > 0 {
			opts.Tags = tags
		}

		list, err := s.ListSkills(c.Request.Context(), opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, list)
	}
}

// GetSkill 获取技能详情
func GetSkill(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的技能 ID"})
			return
		}

		skill, err := s.GetSkillByID(c.Request.Context(), uint(id))
		if err != nil {
			if err == skill.ErrSkillNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "技能不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, skill)
	}
}

// DownloadSkill 下载技能
func DownloadSkill(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的技能 ID"})
			return
		}

		// 获取技能
		skillObj, err := s.GetSkillByID(c.Request.Context(), uint(id))
		if err != nil {
			if err == skill.ErrSkillNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "技能不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 记录下载
		userID := getUserID(c)
		s.IncrementDownload(c.Request.Context(), uint(id), 0, userID, c.ClientIP(), c.Request.UserAgent())

		// 返回下载信息（实际下载由 storage 层处理）
		c.JSON(http.StatusOK, gin.H{
			"download_url": skillObj.LatestVersion,
			"message":      "下载链接已生成",
		})
	}
}

// SearchSkills 搜索技能
func SearchSkills(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
			return
		}

		opts := skill.ListOptions{
			Page:     1,
			PageSize: 20,
			Query:    query,
		}

		if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
			opts.Page = page
		}

		list, err := s.SearchSkills(c.Request.Context(), query, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, list)
	}
}

// ListCategories 获取分类列表
func ListCategories(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := s.ListCategories(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}

// PublishSkill 发布技能
func PublishSkill(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c)
		
		var req skill.Skill
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req.AuthorID = userID

		if err := s.CreateSkill(c.Request.Context(), &req); err != nil {
			if err == skill.ErrSkillAlreadyExists {
				c.JSON(http.StatusConflict, gin.H{"error": "技能名称已存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, req)
	}
}

// UpdateSkill 更新技能
func UpdateSkill(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的技能 ID"})
			return
		}

		userID := getUserID(c)
		
		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.UpdateSkill(c.Request.Context(), uint(id), updates, userID); err != nil {
			if err == skill.ErrUnauthorized {
				c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此技能"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
	}
}

// DeleteSkill 删除技能
func DeleteSkill(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的技能 ID"})
			return
		}

		userID := getUserID(c)

		if err := s.DeleteSkill(c.Request.Context(), uint(id), userID); err != nil {
			if err == skill.ErrUnauthorized {
				c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此技能"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	}
}

// RateSkill 评分
func RateSkill(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的技能 ID"})
			return
		}

		userID := getUserID(c)

		var req struct {
			Rating  int    `json:"rating" binding:"required,min=1,max=5"`
			Comment string `json:"comment"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.RateSkill(c.Request.Context(), uint(id), userID, req.Rating, req.Comment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "评分成功"})
	}
}

// AdminListSkills 管理员获取技能列表（包含待审核）
func AdminListSkills(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.DefaultQuery("status", "all")
		
		// 管理员可以查看所有状态的技能
		// 实现略...
		
		c.JSON(http.StatusOK, gin.H{"message": "管理员功能"})
	}
}

// ApproveSkill 批准技能
func ApproveSkill(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "批准功能"})
	}
}

// RejectSkill 拒绝技能
func RejectSkill(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "拒绝功能"})
	}
}

// GetAnalytics 获取分析数据
func GetAnalytics(s *skill.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"total_skills": 0,
			"total_downloads": 0,
			"trending_skills": []string{},
		})
	}
}

// getUserID 从上下文获取用户 ID
func getUserID(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		if uid, ok := id.(uint); ok {
			return uid
		}
	}
	return 0
}

// parseUint 解析无符号整数
func parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint(v), err
}
