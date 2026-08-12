CREATE TABLE home_page_content (
    id               SMALLINT    PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    draft_config     JSONB       NOT NULL,
    published_config JSONB       NOT NULL,
    config_version   INTEGER     NOT NULL DEFAULT 1 CHECK (config_version > 0),
    updated_by       BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    published_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

WITH default_config AS (
    SELECT '{
      "schema_version": 1,
      "order": ["hero", "intro", "tech_stack", "featured_projects", "latest_posts", "contact_cta"],
      "sections": {
        "hero": {
          "enabled": true,
          "title_prefix": "你好，我是",
          "highlight": "",
          "subtitle": "",
          "primary_button": {"label": "阅读博客", "href": "/blog"},
          "secondary_button": {"label": "查看项目", "href": "/projects"},
          "show_social_links": true
        },
        "intro": {
          "enabled": true,
          "title": "关于我",
          "image_url": "",
          "paragraphs": [
            "我是一名热爱技术的全栈开发者，专注于构建现代化的 Web 应用。我喜欢探索新技术，分享开发经验，并通过代码创造有价值的产品。",
            "这个博客是我记录技术思考、分享项目经验的地方。希望我的文章能对你有所帮助。"
          ]
        },
        "tech_stack": {
          "enabled": true,
          "title": "技术栈",
          "subtitle": "我日常使用的技术和工具",
          "items": ["TypeScript", "React", "Next.js", "Go", "PostgreSQL", "Docker", "Tailwind CSS", "Git"]
        },
        "featured_projects": {
          "enabled": true,
          "title": "精选项目",
          "subtitle": "我最近参与的一些项目",
          "limit": 12
        },
        "latest_posts": {
          "enabled": true,
          "title": "最新文章",
          "subtitle": "分享我的技术思考和经验",
          "limit": 3,
          "source": "latest"
        },
        "contact_cta": {
          "enabled": true,
          "title": "有想法？聊聊吧",
          "description": "无论是项目合作、技术交流，还是简单的问候，我都很乐意收到你的消息。",
          "button": {"label": "联系我", "href": "/contact"}
        }
      }
    }'::jsonb AS value
)
INSERT INTO home_page_content (id, draft_config, published_config)
SELECT 1, value, value FROM default_config
ON CONFLICT (id) DO NOTHING;

CREATE TRIGGER trg_home_page_content_updated_at
    BEFORE UPDATE ON home_page_content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
