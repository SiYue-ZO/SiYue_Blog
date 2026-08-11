# API 概览

所有响应采用 `{ "code": number, "message": string, "data": ... }`。管理端使用同源 HttpOnly Cookie，不要求前端保存令牌。生产默认关闭公开注册。

## 公开 API `/api/v1`

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 数据库与服务健康状态 |
| POST | `/auth/login` | 登录并设置访问/刷新 Cookie |
| POST | `/auth/refresh` | 原子轮换刷新会话 |
| POST | `/auth/logout` | 撤销当前刷新会话并清除 Cookie |
| GET | `/auth/me` | 当前登录用户 |
| GET | `/posts` | 已发布文章分页；支持 `page`、`page_size`、`search`、`tag` |
| GET | `/posts/featured` | 精选已发布文章 |
| GET | `/posts/:slug` | 已发布文章详情 |
| GET | `/tags` | 标签及文章数量 |
| GET | `/projects` | 项目列表 |
| GET | `/projects/featured` | 精选项目 |
| POST | `/contact` | 提交联系消息 |

仅当 `auth.allow_registration=true` 时提供 `POST /auth/register`。

## 管理 API `/api/admin`

所有路由均要求有效访问 Cookie 和 `admin` 角色。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/stats` | 后台统计 |
| POST | `/password` | 修改当前管理员密码并撤销所有刷新会话 |
| GET/POST | `/posts` | 文章列表/创建 |
| GET/PUT/DELETE | `/posts/:id` | 文章读取/更新/删除 |
| GET/POST | `/tags` | 标签列表/创建 |
| DELETE | `/tags/:id` | 删除未被文章使用的标签 |
| GET/POST | `/projects` | 项目列表/创建 |
| PUT/DELETE | `/projects/:id` | 项目更新/删除 |
| GET | `/messages` | 联系消息分页 |
| PATCH/DELETE | `/messages/:id` | 更新消息状态/删除 |

## 安全与限制

- JSON 写请求必须声明 `Content-Type: application/json`。
- 请求体默认最大 2 MiB；登录、刷新、联系表单及全局请求均有限速。
- Cookie 使用 HttpOnly、SameSite=Strict，生产环境强制 Secure。
- 公开文章接口不返回作者邮箱、草稿正文或内部管理字段。
- Markdown 在前端按白名单清洗，文章中的原始 HTML 不会直接执行。
