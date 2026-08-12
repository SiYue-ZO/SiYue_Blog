UPDATE site_settings
SET footer = footer
    - 'filing_icon_url'
    - 'icp_number'
    - 'icp_url'
    - 'public_security_number'
    - 'public_security_url'
    - 'filing_information'
WHERE footer IS NOT NULL;
