package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

// Config 邮件配置
type Config struct {
	SMTPHost     string
	SMTPPort     int
	Username     string
	Password     string
	FromEmail    string
	FromName     string
	EnableTLS    bool
}

// Service 邮件服务
type Service struct {
	config Config
}

// NewService 创建邮件服务
func NewService(config Config) *Service {
	return &Service{config: config}
}

// Email 邮件结构
type Email struct {
	To      []string
	Subject string
	Body    string
	HTML    bool
}

// Send 发送邮件
func (s *Service) Send(email *Email) error {
	// 构建邮件内容
	var body bytes.Buffer

	if email.HTML {
		body.WriteString("MIME-version: 1.0\r\n")
		body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	} else {
		body.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	}

	body.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromEmail))
	body.WriteString(fmt.Sprintf("To: %s\r\n", email.To))
	body.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	body.WriteString("\r\n")
	body.WriteString(email.Body)

	// 构建 SMTP 地址
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	// 认证
	var auth smtp.Auth
	if s.config.Username != "" && s.config.Password != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)
	}

	// 发送
	if s.config.EnableTLS {
		// TODO: 实现 TLS 支持
		return smtp.SendMail(addr+":587", auth, s.config.FromEmail, email.To, body.Bytes())
	}

	return smtp.SendMail(addr, auth, s.config.FromEmail, email.To, body.Bytes())
}

// SendWelcome 发送欢迎邮件
func (s *Service) SendWelcome(to, username string) error {
	subject := "欢迎加入 MCP Skill Hub！"

	body := fmt.Sprintf(`
你好 %s，

欢迎加入 MCP Skill Hub - MCP 技能市场！

在这里，你可以：
- 📦 发现和安装高质量的 MCP 技能
- 🔧 发布和分享你自己开发的技能
- ⭐ 为喜欢的技能评分和评论
- 📊 追踪技能的使用情况

开始探索：https://mcp-skill-hub.dev/skills

如有任何问题，请随时联系我们。

祝好，
MCP Skill Hub 团队
`, username)

	return s.Send(&Email{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		HTML:    false,
	})
}

// SendSkillApproved 发送技能审核通过通知
func (s *Service) SendSkillApproved(to, username, skillName string) error {
	subject := fmt.Sprintf("你的技能 \"%s\" 已通过审核", skillName)

	body := fmt.Sprintf(`
你好 %s，

好消息！你的技能 "%s" 已通过审核，现在可以在 MCP Skill Hub 上被其他用户发现和安装了。

查看你的技能：https://mcp-skill-hub.dev/skills/%s

感谢你的贡献！

祝好，
MCP Skill Hub 团队
`, username, skillName, skillName)

	return s.Send(&Email{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		HTML:    false,
	})
}

// SendSkillRejected 发送技能审核拒绝通知
func (s *Service) SendSkillRejected(to, username, skillName, reason string) error {
	subject := fmt.Sprintf("你的技能 \"%s\" 审核未通过", skillName)

	body := fmt.Sprintf(`
你好 %s，

很遗憾，你的技能 "%s" 未能通过审核。

原因：%s

你可以修改后重新提交。如有异议，请回复此邮件联系我们。

祝好，
MCP Skill Hub 团队
`, username, skillName, reason)

	return s.Send(&Email{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		HTML:    false,
	})
}

// SendPasswordReset 发送密码重置邮件
func (s *Service) SendPasswordReset(to, username, resetLink string) error {
	subject := "MCP Skill Hub 密码重置"

	body := fmt.Sprintf(`
你好 %s，

你请求重置 MCP Skill Hub 账户的密码。

点击以下链接重置密码（链接 24 小时内有效）：
%s

如果你没有请求重置密码，请忽略此邮件。

祝好，
MCP Skill Hub 团队
`, username, resetLink)

	return s.Send(&Email{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		HTML:    false,
	})
}

// SendNewVersion 发送新版本通知
func (s *Service) SendNewVersion(to, username, skillName, version string) error {
	subject := fmt.Sprintf("你关注的技能 \"%s\" 发布了新版本", skillName)

	body := fmt.Sprintf(`
你好 %s，

你关注的技能 "%s" 发布了新版本：%s

立即查看：https://mcp-skill-hub.dev/skills/%s

祝好，
MCP Skill Hub 团队
`, username, skillName, version, skillName)

	return s.Send(&Email{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		HTML:    false,
	})
}

// SendWeeklyDigest 发送周报邮件
func (s *Service) SendWeeklyDigest(to, username string, topSkills []DigestSkill) error {
	subject := "MCP Skill Hub 周报 - 本周热门技能"

	var skillsHTML string
	for i, skill := range topSkills {
		skillsHTML += fmt.Sprintf(`
<tr>
  <td>%d</td>
  <td><a href="https://mcp-skill-hub.dev/skills/%s">%s</a></td>
  <td>%s</td>
  <td>%d</td>
</tr>
`, i+1, skill.Slug, skill.Name, skill.Description, skill.Downloads)
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: Arial, sans-serif; }
    table { border-collapse: collapse; width: 100%%; }
    th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
    th { background-color: #f2f2f2; }
  </style>
</head>
<body>
  <h2>你好 %s，</h2>
  <p>这是本周的 MCP Skill Hub 热门技能排行榜：</p>
  <table>
    <tr>
      <th>排名</th>
      <th>技能名称</th>
      <th>描述</th>
      <th>下载量</th>
    </tr>
    %s
  </table>
  <p>立即探索更多技能：<a href="https://mcp-skill-hub.dev">mcp-skill-hub.dev</a></p>
  <p>祝好，<br/>MCP Skill Hub 团队</p>
</body>
</html>
`, username, skillsHTML)

	return s.Send(&Email{
		To:      []string{to},
		Subject: subject,
		Body:    htmlBody,
		HTML:    true,
	})
}

// DigestSkill 周报技能信息
type DigestSkill struct {
	Slug        string
	Name        string
	Description string
	Downloads   int64
}

// TemplateEmail 使用模板发送邮件
func (s *Service) SendTemplate(to, subject, templateName string, data interface{}) error {
	// TODO: 实现模板加载和渲染
	return nil
}

// BatchSend 批量发送邮件
func (s *Service) BatchSend(emails []*Email, batchSize int) error {
	for i := 0; i < len(emails); i += batchSize {
		end := i + batchSize
		if end > len(emails) {
			end = len(emails)
		}

		for _, email := range emails[i:end] {
			if err := s.Send(email); err != nil {
				// 记录错误，继续发送
				continue
			}
			// 避免触发速率限制
			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}
