# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-04-24

### Added
- 🎉 Initial release
- 📦 Core skill management (publish, search, download, rate)
- 🔐 User authentication with JWT
- 🗄️ PostgreSQL database integration
- 📧 Email service with TLS support
- 📁 Object storage (MinIO/S3 compatible)
- 🐳 Docker deployment support
- 📊 Basic analytics and statistics
- 🔍 Full-text search for skills
- 📝 API documentation

### Security
- Password hashing with bcrypt
- JWT token authentication
- Rate limiting middleware
- Input validation and sanitization

### Dependencies
- Go 1.21+
- PostgreSQL 14+
- Redis 6+ (optional, for caching)
- MinIO or S3-compatible storage

---

## Future Roadmap

### [0.2.0] - Planned
- Webhook support
- Skill dependency management
- Advanced analytics dashboard
- Multi-language support

### [0.3.0] - Planned
- GraphQL API
- Real-time notifications
- Skill marketplace UI

### [1.0.0] - Planned
- Enterprise features (SSO, audit logs)
- High availability support
- Performance optimizations
