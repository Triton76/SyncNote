#!/usr/bin/env bash
set -euo pipefail

API_HOST="${API_HOST:-http://127.0.0.1:8888}"
AUTH_HOST="${AUTH_HOST:-http://127.0.0.1:8889}"
USER_ID="${USER_ID:-api_test_user_$(date +%s)}"
AUTH_EMAIL="${AUTH_EMAIL:-api_test_$(date +%s)@example.com}"
AUTH_PASSWORD="${AUTH_PASSWORD:-123456}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd curl
need_cmd jq

echo "[0/9] Register auth user and get token..."
auth_resp="$({
  curl -sS -X POST "$AUTH_HOST/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$USER_ID\",\"email\":\"$AUTH_EMAIL\",\"password\":\"$AUTH_PASSWORD\",\"captcha\":\"233\"}"
} )"

echo "$auth_resp" | jq . >/dev/null
token="$(echo "$auth_resp" | jq -r '.token // .Token // empty')"
auth_user_id="$(echo "$auth_resp" | jq -r '.userId // .UserId // empty')"
if [[ -z "$token" || -z "$auth_user_id" ]]; then
  echo "Auth register failed: $auth_resp"
  exit 1
fi

echo "[1/9] Create note via API..."
create_resp="$({
  curl -sS -X POST "$API_HOST/api/note/create" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"api smoke title\",\"content\":\"api smoke content\"}"
} )"

echo "$create_resp" | jq . >/dev/null
note_id="$(echo "$create_resp" | jq -r '.noteId // .NoteId // empty')"
version="$(echo "$create_resp" | jq -r '.version // .Version // empty')"
if [[ -z "$note_id" || -z "$version" ]]; then
  echo "CreateNote failed: $create_resp"
  exit 1
fi

echo "Created note_id=$note_id version=$version"

echo "[2/9] Get note via API..."
get_resp="$(curl -sS -H "Authorization: Bearer $token" "$API_HOST/api/note/$note_id")"
echo "$get_resp" | jq . >/dev/null
get_note_id="$(echo "$get_resp" | jq -r '.noteId // .NoteId // empty')"
if [[ "$get_note_id" != "$note_id" ]]; then
  echo "GetNote mismatch: expected $note_id, got $get_note_id"
  exit 1
fi

echo "[3/9] Save note success path..."
save_resp="$({
  curl -sS -X POST "$API_HOST/api/note/save" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: application/json" \
    -d "{\"noteId\":\"$note_id\",\"content\":\"api updated content\",\"expectedVersion\":$version}"
} )"
echo "$save_resp" | jq . >/dev/null
save_success="$(echo "$save_resp" | jq -r '(.success // .Success // false | tostring)')"
if [[ "$save_success" != "true" ]]; then
  echo "SaveNote failed: $save_resp"
  exit 1
fi

new_version="$(echo "$save_resp" | jq -r '.note.version // .note.Version // .Note.version // .Note.Version // empty')"
if [[ -z "$new_version" ]]; then
  echo "SaveNote response missing note.version: $save_resp"
  exit 1
fi

echo "[4/9] Save note conflict path..."
conflict_resp="$({
  curl -sS -X POST "$API_HOST/api/note/save" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: application/json" \
    -d "{\"noteId\":\"$note_id\",\"content\":\"api conflict content\",\"expectedVersion\":$version}"
} )"
echo "$conflict_resp" | jq . >/dev/null
conflict_code="$(echo "$conflict_resp" | jq -r '.code // .Code // empty')"
if [[ "$conflict_code" != "SAVE_CODE_VERSION_CONFLICT" ]]; then
  echo "Conflict code mismatch, expected SAVE_CODE_VERSION_CONFLICT, got '$conflict_code'"
  echo "Response: $conflict_resp"
  exit 1
fi

echo "[5/9] Get user notes via API..."
list_resp="$(curl -sS -H "Authorization: Bearer $token" "$API_HOST/api/user/notes")"
echo "$list_resp" | jq . >/dev/null
found="$(echo "$list_resp" | jq -r --arg id "$note_id" '[(.notes // .Notes // [])[]? | select((.noteId // .NoteId // "") == $id)] | length')"
if [[ "$found" == "0" ]]; then
  echo "GetUserNotes did not return created note"
  echo "Response: $list_resp"
  exit 1
fi

target_user_id="collab_user_001"

echo "[6/9] Grant permission via API..."
grant_resp="$({
  curl -sS -X POST "$API_HOST/api/permission/grant" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: application/json" \
    -d "{\"noteId\":\"$note_id\",\"targetUserId\":\"$target_user_id\",\"role\":\"editor\"}"
} )"
echo "$grant_resp" | jq . >/dev/null
grant_success="$(echo "$grant_resp" | jq -r '(.success // .Success // false | tostring)')"
if [[ "$grant_success" != "true" ]]; then
  echo "GrantPermission failed: $grant_resp"
  exit 1
fi

echo "[7/9] List permissions via API..."
perm_list_resp="$(curl -sS -H "Authorization: Bearer $token" "$API_HOST/api/permission/list/$note_id")"
echo "$perm_list_resp" | jq . >/dev/null
perm_found="$(echo "$perm_list_resp" | jq -r --arg uid "$target_user_id" '[(.permissions // .Permissions // [])[]? | select((.userId // .UserId // "") == $uid and (.status // .Status // "") == "active")] | length')"
if [[ "$perm_found" == "0" ]]; then
  echo "ListPermissions did not return expected active permission"
  echo "Response: $perm_list_resp"
  exit 1
fi

echo "[8/9] Revoke permission via API..."
revoke_resp="$({
  curl -sS -X POST "$API_HOST/api/permission/revoke" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: application/json" \
    -d "{\"noteId\":\"$note_id\",\"targetUserId\":\"$target_user_id\"}"
} )"
echo "$revoke_resp" | jq . >/dev/null
revoke_success="$(echo "$revoke_resp" | jq -r '(.success // .Success // false | tostring)')"
if [[ "$revoke_success" != "true" ]]; then
  echo "RevokePermission failed: $revoke_resp"
  exit 1
fi

echo "[9/9] Get note events via API..."
events_resp="$(curl -sS -H "Authorization: Bearer $token" "$API_HOST/api/note/$note_id/events?startSeq=0&limit=20")"
echo "$events_resp" | jq . >/dev/null
if ! echo "$events_resp" | jq -e '.events // .Events' >/dev/null 2>&1; then
  echo "GetNoteEvents response missing events field"
  echo "Response: $events_resp"
  exit 1
fi

echo "API smoke test passed."
