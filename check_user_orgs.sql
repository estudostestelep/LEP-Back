SELECT
    uo.id,
    uo.user_id,
    uo.organization_id,
    uo.role,
    uo.active,
    u.email,
    u.name
FROM user_organizations uo
JOIN users u ON u.id = uo.user_id
WHERE uo.user_id = '123e4567-e89b-12d3-a456-426614174010';
