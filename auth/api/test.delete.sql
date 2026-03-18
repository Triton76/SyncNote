USE syncnote;
-- Auth service test data cleanup
-- Safe to run multiple times.
DELETE FROM users
WHERE email = 'test01@example.com'
	OR email LIKE 'test%@example.com'
	OR username = 'test_user_01'
	OR username LIKE 'test_user_%'
	OR username LIKE 'api_test_user_%'
	OR username LIKE 'rpc_test_user_%';