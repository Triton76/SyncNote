#!/usr/bin/env bash
set -euo pipefail

API_HOST="${API_HOST:-http://127.0.0.1:8888}"
AUTH_HOST="${AUTH_HOST:-http://127.0.0.1:8889}"
TS="$(date +%s)"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd curl
need_cmd jq

register_user() {
  local username="$1"
  local email="$2"
  local password="$3"

  local resp
  resp="$({
    curl -sS -X POST "$AUTH_HOST/auth/register" \
      -H "Content-Type: application/json" \
      -d "{\"username\":\"$username\",\"email\":\"$email\",\"password\":\"$password\",\"captcha\":\"233\"}"
  } )"

  echo "$resp" | jq . >/dev/null

  local token
  token="$(echo "$resp" | jq -r '.token // .Token // empty')"
  if [[ -z "$token" ]]; then
    echo "Register failed: $resp"
    exit 1
  fi

  printf '%s' "$token"
}

echo "[1/4] Register user A and B ..."
user_a="authz_user_a_${TS}"
user_b="authz_user_b_${TS}"
email_a="authz_a_${TS}@example.com"
email_b="authz_b_${TS}@example.com"
password="123456"

token_a="$(register_user "$user_a" "$email_a" "$password")"
token_b="$(register_user "$user_b" "$email_b" "$password")"

echo "[2/4] User A creates a note ..."
create_resp="$({
  curl -sS -X POST "$API_HOST/api/note/create" \
    -H "Authorization: Bearer $token_a" \
    -H "Content-Type: application/json" \
    -d '{"title":"authz smoke title","content":"authz smoke content"}'
} )"

echo "$create_resp" | jq . >/dev/null
note_id="$(echo "$create_resp" | jq -r '.noteId // .NoteId // empty')"
if [[ -z "$note_id" ]]; then
  echo "Create note failed: $create_resp"
  exit 1
fi

echo "[3/4] User B tries to read user A note (should be 403) ..."
status_code="$(curl -sS -o /tmp/syncnote_authz_body.json -w '%{http_code}' \
  -H "Authorization: Bearer $token_b" \
  "$API_HOST/api/note/$note_id")"

if [[ "$status_code" != "403" ]]; then
  echo "Expected 403, got $status_code"
  echo "Response body:"
  cat /tmp/syncnote_authz_body.json
  exit 1
fi

echo "[4/4] Validate response payload ..."
cat /tmp/syncnote_authz_body.json | jq . >/dev/null
resp_code="$(cat /tmp/syncnote_authz_body.json | jq -r '.code // empty')"
if [[ "$resp_code" != "403" ]]; then
  echo "Expected response code=403 in JSON body, got '$resp_code'"
  cat /tmp/syncnote_authz_body.json
  exit 1
fi

echo "API authz smoke test passed."
