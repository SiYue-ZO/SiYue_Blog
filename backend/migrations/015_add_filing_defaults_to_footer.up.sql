UPDATE site_settings
SET footer = COALESCE(footer, '{}'::jsonb) || jsonb_build_object(
    'filing_icon_url', COALESCE(footer->>'filing_icon_url', ''),
    'icp_number', COALESCE(footer->>'icp_number', ''),
    'icp_url', COALESCE(footer->>'icp_url', 'https://beian.miit.gov.cn/'),
    'public_security_number', COALESCE(footer->>'public_security_number', ''),
    'public_security_url', COALESCE(footer->>'public_security_url', ''),
    'filing_information', COALESCE(footer->>'filing_information', '')
)
WHERE footer IS NOT NULL;
