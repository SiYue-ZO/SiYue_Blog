# API 概览

所有响应采用 `{ "code": number, "message": string, "data": ... }`。管理端使用同源 HttpOnly Cookie，不要求前端保存令牌。生产默认关闭公开注册。

## 公开 API `/api/v1`

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 数据库与服务健康状态 |
| GET | `/settings` | 公开站点设置（项目页状态、联系页社交账号与页脚） |
| GET | `/home` | 已发布的首页区块配置 |
| POST | `/auth/login` | 登录并设置访问/刷新 Cookie |
| POST | `/auth/refresh` | 原子轮换刷新会话 |
| POST | `/auth/logout` | 撤销当前刷新会话并清除 Cookie |
| GET | `/auth/me` | 当前登录用户 |
| GET | `/posts` | 已发布文章分页；支持 `page`、`page_size`、`search`、`tag` |
| GET | `/posts/featured` | 精选已发布文章 |
| GET | `/posts/:slug` | 已发布文章详情；纯读取，不增加阅读次数 |
| POST | `/posts/:slug/views` | 上报阅读事件，JSON `{ "event_id": "UUID v4" }`；返回 `view_count` |
| GET | `/tags` | 标签及文章数量 |
| GET | `/projects` | 项目列表 |
| GET | `/projects/featured` | 精选项目 |
| POST | `/contact` | 提交联系消息 |

仅当 `auth.allow_registration=true` 时提供 `POST /auth/register`。

阅读次数按同一标签页会话内首次打开某篇文章计数。阅读事件 ID 按文章 ID 保存在 `sessionStorage` 中，刷新、失败重试以及离开后返回该文章均复用原事件 ID，不再增加计数；新的独立标签页会话首次阅读才增加。
元数据请求、预取和直接 GET 文章不计数。浏览器禁止会话存储时跳过阅读上报，文章仍可正常阅读。
同一篇文章的相同 `event_id` 只计一次，去重记录与计数在数据库事务中同时提交，跨服务重启及多实例仍然有效。

## 管理 API `/api/admin`

所有路由均要求有效访问 Cookie 和 `admin` 角色。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/stats` | 后台统计 |
| GET/PATCH | `/settings` | 读取/更新站点设置；支持 `project_page_visible`、`social_links`、`footer` |
| GET | `/home` | 读取首页草稿、发布配置和版本信息 |
| PUT | `/home/draft` | 保存首页草稿配置 |
| POST | `/home/publish` | 发布当前首页草稿 |
| POST | `/home/reset` | 将首页草稿恢复为已发布版本 |
| POST | `/password` | 修改当前管理员密码并撤销所有刷新会话 |
| GET/POST | `/posts` | 文章列表/创建 |
| GET/PUT/DELETE | `/posts/:id` | 文章读取/更新/删除 |
| GET/POST | `/tags` | 标签列表/创建 |
| DELETE | `/tags/:id` | 删除未被文章使用的标签 |
| GET/POST | `/projects` | 项目列表/创建 |
| PUT/DELETE | `/projects/:id` | 项目更新/删除 |
| GET | `/messages` | 联系消息分页 |
| PATCH/DELETE | `/messages/:id` | 更新消息状态/删除 |

`footer` 支持以下可选备案字段：

```json
{
  "filing_icon_url": "/images/beian.png",
  "icp_number": "粤ICP备12345678号",
  "icp_url": "https://beian.miit.gov.cn/",
  "public_security_number": "粤公网安备 44000000000000号",
  "public_security_url": "https://www.beian.gov.cn/portal/registerSystemInfo",
  "filing_information": "备案附加信息"
}
```

图标支持站内绝对路径或 HTTPS 地址；备案链接只接受 HTTPS 地址。空字段不会在公开页脚中占位。

## 安全与限制

- JSON 写请求必须声明 `Content-Type: application/json`。
- 请求体默认最大 2 MiB；登录、刷新、联系表单及全局请求均有限速。
- Cookie 使用 HttpOnly、SameSite=Strict，生产环境强制 Secure。
- 公开文章接口不返回作者邮箱、草稿正文或内部管理字段。
- Markdown 在前端按白名单清洗，文章中的原始 HTML 不会直接执行。
