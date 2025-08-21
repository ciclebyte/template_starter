-- 更新admin用户的密码为正确的bcrypt哈希 (密码: admin123)
UPDATE users 
SET password_hash = '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi' 
WHERE username = 'admin';

-- 确认更新
SELECT username, password_hash FROM users WHERE username = 'admin';