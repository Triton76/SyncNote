#!/usr/bin/env bash
set -euo pipefail

AUTH_HOST="${AUTH_HOST:-http://127.0.0.1:8000}"
USER_API_HOST="${USER_API_HOST:-http://127.0.0.1:8888}"
AUTH_EMAIL="${AUTH_EMAIL:-rebuild_userapi_$(date +%s)@example.com}"
AUTH_PASSWORD="${AUTH_PASSWORD:-123456}"
NEW_USERNAME="${NEW_USERNAME:-tester_$(date +%s)}"
NEW_SYNOPSIS="${NEW_SYNOPSIS:-updated by userapi smoke}"
NEW_AVATAR_URL="${NEW_AVATAR_URL:-https://example.com/avatar.png}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd curl
need_cmd jq

http_json_request() {
  local method="$1"
  local url="$2"
  local token="${3:-}"
  local body="${4:-}"

  local tmp_body
  tmp_body="$(mktemp)"

  local code
  if [[ -n "$token" ]]; then
    if [[ -n "$body" ]]; then
      code="$(curl -sS -o "$tmp_body" -w "%{http_code}" -X "$method" "$url" -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "$body")"
    else
      code="$(curl -sS -o "$tmp_body" -w "%{http_code}" -X "$method" "$url" -H "Authorization: Bearer $token")"
    fi
  else
    if [[ -n "$body" ]]; then
      code="$(curl -sS -o "$tmp_body" -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -d "$body")"
    else
      code="$(curl -sS -o "$tmp_body" -w "%{http_code}" -X "$method" "$url")"
    fi
  fi

  local resp
  resp="$(cat "$tmp_body")"
  rm -f "$tmp_body"

  if [[ "$code" -lt 200 || "$code" -ge 300 ]]; then
    echo "HTTP request failed: $method $url" >&2
    echo "status=$code" >&2
    if [[ -n "$resp" ]]; then
      echo "response=$resp" >&2
    else
      echo "response=<empty>" >&2
    fi
    return 1
  fi

  if ! echo "$resp" | jq . >/dev/null 2>&1; then
    echo "Response is not valid JSON: $method $url" >&2
    echo "status=$code" >&2
    if [[ -n "$resp" ]]; then
      echo "response=$resp" >&2
    else
      echo "response=<empty>" >&2
    fi
    return 1
  fi

  printf "%s" "$resp"
}

echo "[1/7] Register user via authapi"
register_resp="$(http_json_request "POST" "$AUTH_HOST/auth/register" "" "{\"email\":\"$AUTH_EMAIL\",\"password\":\"$AUTH_PASSWORD\",\"captcha\":\"\"}")"
user_id="$(echo "$register_resp" | jq -r '.userId // .UserId // empty')"
if [[ -z "$user_id" ]]; then
  echo "Register failed: $register_resp"
  exit 1
fi

echo "[2/7] Login via authapi"
login_resp="$(http_json_request "POST" "$AUTH_HOST/auth/login" "" "{\"email\":\"$AUTH_EMAIL\",\"password\":\"$AUTH_PASSWORD\",\"captcha\":\"\"}")"
token="$(echo "$login_resp" | jq -r '.token // .Token // empty')"
if [[ -z "$token" ]]; then
  echo "Login failed: $login_resp"
  exit 1
fi

echo "[3/7] Get self info"
me_resp="$(http_json_request "GET" "$USER_API_HOST/api/user/me" "$token")"
me_user_id="$(echo "$me_resp" | jq -r '.userId // .UserId // empty')"
if [[ "$me_user_id" != "$user_id" ]]; then
  echo "GetSelfInfo mismatch: expected $user_id, got $me_user_id"
  echo "Response: $me_resp"
  exit 1
fi

echo "[4/7] Edit self info"
edit_resp="$(http_json_request "POST" "$USER_API_HOST/api/user/edit/me" "$token" "{\"username\":\"$NEW_USERNAME\",\"synopsis\":\"$NEW_SYNOPSIS\",\"avatarUrl\":\"$NEW_AVATAR_URL\"}")"

echo "[5/7] Verify self info after edit"
me_after_resp="$(http_json_request "GET" "$USER_API_HOST/api/user/me" "$token")"
me_after_username="$(echo "$me_after_resp" | jq -r '.username // .Username // empty')"
if [[ "$me_after_username" != "$NEW_USERNAME" ]]; then
  echo "Edit verification failed: expected username=$NEW_USERNAME, got $me_after_username"
  echo "Response: $me_after_resp"
  exit 1
fi

echo "[6/7] Get user info by userId"
by_id_resp="$(http_json_request "GET" "$USER_API_HOST/api/user/$user_id/info" "$token")"
by_id_user_id="$(echo "$by_id_resp" | jq -r '.userId // .UserId // empty')"
if [[ "$by_id_user_id" != "$user_id" ]]; then
  echo "GetUserInfo mismatch: expected $user_id, got $by_id_user_id"
  echo "Response: $by_id_resp"
  exit 1
fi

echo "[7/7] Search user by email"
search_resp="$(http_json_request "GET" "$USER_API_HOST/api/user/search?email=$AUTH_EMAIL" "$token")"
search_total="$(echo "$search_resp" | jq -r '.total // .Total // 0')"
if [[ "$search_total" == "0" ]]; then
  echo "SearchUser returned empty result"
  echo "Response: $search_resp"
  exit 1
fi

found_id="$(echo "$search_resp" | jq -r --arg uid "$user_id" '[(.infoList // .InfoList // [])[]? | select((.userId // .UserId // "") == $uid)] | length')"
if [[ "$found_id" == "0" ]]; then
  echo "SearchUser did not include expected userId=$user_id"
  echo "Response: $search_resp"
  exit 1
fi

echo "userapi smoke test passed"
echo "userId=$user_id"
echo "email=$AUTH_EMAIL"
