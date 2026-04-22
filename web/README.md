# MCP Skill Hub Web

React 前端管理界面

## 快速开始

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## 功能特性

- 📊 技能浏览和搜索
- 🔐 用户认证（登录/注册）
- 📦 技能发布和管理
- ⭐ 评分和评论
- 👤 用户个人中心
- 🔑 API Key 管理

## 技术栈

- **React 18** - UI 框架
- **Vite** - 构建工具
- **React Router** - 路由
- **Zustand** - 状态管理
- **TanStack Query** - 数据获取
- **Axios** - HTTP 客户端
- **Tailwind CSS** - 样式
- **Lucide React** - 图标

## 项目结构

```
web/
├── src/
│   ├── components/       # 可复用组件
│   ├── pages/           # 页面组件
│   ├── hooks/           # 自定义 Hooks
│   ├── stores/          # Zustand stores
│   ├── api/             # API 客户端
│   ├── utils/           # 工具函数
│   └── App.jsx          # 根组件
├── index.html
├── vite.config.js
├── tailwind.config.js
└── package.json
```

## 开发规范

- 使用函数组件 + Hooks
- 使用 Tailwind CSS 进行样式设计
- 使用 Lucide React 图标
- 遵循 React 最佳实践
