# 📚 MCP RAPI 项目文档总览

> **MCP RAPI** - 现代化API权限管理系统完整文档指南

## 🎯 文档导航

### 🚀 快速开始

| 文档 | 描述 | 适用人群 | 预计时间 |
|------|------|----------|----------|
| [📖 README.md](README.md) | 项目介绍和基础信息 | 所有用户 | 5分钟 |
| [⚡ QUICK_START.md](QUICK_START.md) | 快速启动指南 | 新手用户 | 10分钟 |
| [📋 DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) | **完整部署指南** | 开发者、运维 | 30分钟 |

### 🏗️ 系统架构

| 文档 | 描述 | 适用人群 | 详细程度 |
|------|------|----------|----------|
| [🏗️ API多租户授权管理系统架构文档.md](API多租户授权管理系统架构文档.md) | 系统整体架构设计 | 架构师、开发者 | ⭐⭐⭐⭐⭐ |
| [🗃️ DATABASE_INIT.md](DATABASE_INIT.md) | 数据库初始化指南 | 数据库管理员 | ⭐⭐⭐⭐ |
| [🗃️ DATABASE_INIT_GUIDE.md](DATABASE_INIT_GUIDE.md) | 数据库详细配置 | 数据库管理员 | ⭐⭐⭐⭐⭐ |

### 🐳 部署运维

| 文档 | 描述 | 适用人群 | 复杂度 |
|------|------|----------|---------|
| [🐳 DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) | Docker专项部署 | DevOps工程师 | ⭐⭐⭐ |
| [📋 DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) | **完整部署指南** | 所有技术人员 | ⭐⭐⭐⭐⭐ |

### 🔧 功能模块

| 文档 | 描述 | 适用人群 | 重要性 |
|------|------|----------|---------|
| [🏢 新增部门流程文档.md](新增部门流程文档.md) | 部门管理流程 | 业务管理员 | ⭐⭐⭐⭐ |

### 🤖 MCP集成

| 文档 | 描述 | 适用人群 | 技术难度 |
|------|------|----------|----------|
| [🤖 mcp-example/FASTMCP_HTTP_CONFIG.md](mcp-example/FASTMCP_HTTP_CONFIG.md) | MCP集成配置 | AI开发者 | ⭐⭐⭐ |

### 📋 项目管理

| 文档 | 描述 | 适用人群 | 类型 |
|------|------|----------|------|
| [🤝 CONTRIBUTING.md](CONTRIBUTING.md) | 贡献指南 | 开源贡献者 | 规范 |
| [📜 CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) | 行为准则 | 社区成员 | 规范 |
| [🔒 SECURITY.md](SECURITY.md) | 安全政策 | 安全研究员 | 政策 |
| [📄 LICENSE](LICENSE) | 开源协议 | 法务、开发者 | 法律 |
| [Makefile](Makefile) | 项目构建和管理命令 | 开发者、运维 | 工具 |
| [scripts/check_project_status.sh](scripts/check_project_status.sh) | 项目状态检查脚本 | 所有用户 | 工具 |

**快速命令:**
- `make status` - 检查项目完整状态
- `make help` - 查看所有可用命令
- `make docs` - 查看文档总览

---

## 🎯 按角色推荐阅读路径

### 👨‍💻 开发者路径

1. **入门阶段**
   - [📖 README.md](README.md) - 了解项目概况
   - [⚡ QUICK_START.md](QUICK_START.md) - 快速搭建开发环境
   - [📋 DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) - 学习部署方式

2. **深入阶段**
   - [🏗️ 系统架构文档](API多租户授权管理系统架构文档.md) - 理解系统设计
   - [🗃️ DATABASE_INIT_GUIDE.md](DATABASE_INIT_GUIDE.md) - 掌握数据库结构
   - [🤖 MCP集成指南](mcp-example/FASTMCP_HTTP_CONFIG.md) - 学习MCP集成

3. **贡献阶段**
   - [🤝 CONTRIBUTING.md](CONTRIBUTING.md) - 了解贡献流程
   - [📜 CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) - 遵守社区规范

### 🔧 运维工程师路径

1. **部署阶段**
   - [📋 DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) - **核心文档**
   - [🐳 DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) - Docker专项
   - [🗃️ DATABASE_INIT.md](DATABASE_INIT.md) - 数据库配置

2. **优化阶段**
   - [📋 性能优化章节](DEPLOYMENT_GUIDE.md#-性能优化) - 系统调优
   - [🔒 安全配置章节](DEPLOYMENT_GUIDE.md#-安全配置) - 安全加固
   - [📈 监控日志章节](DEPLOYMENT_GUIDE.md#-监控与日志) - 运维监控

3. **维护阶段**
   - [🛠️ 故障排除章节](DEPLOYMENT_GUIDE.md#️-故障排除) - 问题诊断
   - [🔒 SECURITY.md](SECURITY.md) - 安全响应

### 👨‍💼 产品经理路径

1. **产品理解**
   - [📖 README.md](README.md) - 产品功能概览
   - [🏗️ 系统架构文档](API多租户授权管理系统架构文档.md) - 技术架构

2. **业务流程**
   - [🏢 新增部门流程文档.md](新增部门流程文档.md) - 核心业务流程
   - [🤖 MCP集成指南](mcp-example/FASTMCP_HTTP_CONFIG.md) - 集成能力

### 🎓 新手用户路径

1. **第一步**: [📖 README.md](README.md) - 5分钟了解项目
2. **第二步**: [⚡ QUICK_START.md](QUICK_START.md) - 10分钟快速体验
3. **第三步**: [📋 DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) - 30分钟掌握部署
4. **进阶**: 根据需要选择其他专项文档

---

## 📊 文档质量评级

| 文档 | 完整性 | 准确性 | 易读性 | 维护状态 |
|------|--------|--------|--------|----------|
| [📋 DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 🟢 最新 |
| [📖 README.md](README.md) | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 🟢 最新 |
| [⚡ QUICK_START.md](QUICK_START.md) | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 🟢 最新 |
| [🏗️ 系统架构文档](API多租户授权管理系统架构文档.md) | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 🟡 需更新 |
| [🐳 DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 🟡 需更新 |

---

## 🔄 文档更新计划

### 🎯 近期计划 (本月)

- [x] ✅ 完善 [📋 DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)
- [x] ✅ 更新 [📖 README.md](README.md) 
- [ ] 🔄 更新 [🐳 DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)
- [ ] 🔄 完善 [🏗️ 系统架构文档](API多租户授权管理系统架构文档.md)

### 📅 中期计划 (下月)

- [ ] 📝 新增 API 详细文档
- [ ] 📝 新增 前端开发指南
- [ ] 📝 新增 测试指南
- [ ] 📝 新增 性能测试报告

### 🚀 长期计划 (季度)

- [ ] 📚 多语言文档支持
- [ ] 🎥 视频教程制作
- [ ] 📱 移动端文档优化
- [ ] 🤖 AI辅助文档生成

---

## 🤝 文档贡献

### 📝 如何贡献文档

1. **发现问题**: 在使用过程中发现文档问题
2. **提交Issue**: 在 [GitHub Issues](https://github.com/chuangyeshuo/mcprapi/issues) 中报告
3. **提交PR**: 直接修改文档并提交Pull Request
4. **参与讨论**: 在 [GitHub Discussions](https://github.com/chuangyeshuo/mcprapi/discussions) 中讨论

### 📋 文档规范

- **格式**: 使用 Markdown 格式
- **结构**: 遵循现有文档结构
- **语言**: 中文为主，关键术语保留英文
- **图片**: 使用 SVG 格式，存放在 `assets/` 目录
- **链接**: 使用相对路径，确保链接有效

### 🏆 贡献者

感谢所有为文档做出贡献的朋友们！

---

## 📞 获取帮助

如果您在阅读文档过程中遇到问题：

1. **📖 查看相关文档**: 先查看是否有相关文档解答
2. **🔍 搜索已知问题**: 在 Issues 中搜索类似问题
3. **💬 提问讨论**: 在 Discussions 中提出问题
4. **🐛 报告问题**: 如果是文档错误，请提交 Issue
5. **📧 联系我们**: support@mcprapi.com

---

## 🎯 快速链接

**🚀 立即开始：**
- [⚡ 一键启动](QUICK_START.md#-一键启动)
- [🐳 Docker部署](DEPLOYMENT_GUIDE.md#-docker-compose-一键部署推荐)
- [💻 手动部署](DEPLOYMENT_GUIDE.md#-手动部署)

**📚 深入学习：**
- [🏗️ 系统架构](API多租户授权管理系统架构文档.md)
- [🔒 安全配置](DEPLOYMENT_GUIDE.md#-安全配置)
- [📊 性能优化](DEPLOYMENT_GUIDE.md#-性能优化)

**🤝 参与贡献：**
- [🤝 贡献指南](CONTRIBUTING.md)
- [🐛 报告问题](https://github.com/chuangyeshuo/mcprapi/issues)
- [💬 参与讨论](https://github.com/chuangyeshuo/mcprapi/discussions)

---

<div align="center">

**📚 文档即代码，让每个人都能轻松上手 MCP RAPI！**

[⭐ Star 项目](https://github.com/chuangyeshuo/mcprapi) • [📖 改进文档](https://github.com/chuangyeshuo/mcprapi/issues) • [💬 加入讨论](https://github.com/chuangyeshuo/mcprapi/discussions)

</div>