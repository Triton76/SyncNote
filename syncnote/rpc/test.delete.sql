USE syncnote;
-- SyncNote service test data cleanup
-- Safe to run multiple times.
DELETE FROM notes
WHERE user_id LIKE 'api_test_user_%'
	OR user_id LIKE 'rpc_test_user_%'
	OR user_id = 'test_user_01';
-- Optional: clean up any leftover fixed test note ids if you add them later.
-- DELETE FROM notes WHERE note_id IN ('test-note-1', 'test-note-2');