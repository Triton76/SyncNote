#!/usr/bin/env bash
set -euo pipefail

API_HOST="${API_HOST:-http://127.0.0.1:8888}"
USER_ID="${USER_ID:-api_test_user_$(date +%s)}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd curl
need_cmd jq

echo "[1/5] Create note via API..."
create_resp="$({
  curl -sS -X POST "$API_HOST/api/note/create" \
    -H "Content-Type: application/json" \
    -d "{\"userId\":\"$USER_ID\",\"title\":\"api smoke title\",\"content\":\"api smoke content\"}"
} )"

echo "$create_resp" | jq . >/dev/null
note_id="$(echo "$create_resp" | jq -r '.noteId // .NoteId // empty')"
version="$(echo "$create_resp" | jq -r '.version // .Version // empty')"
if [[ -z "$note_id" || -z "$version" ]]; then
  echo "CreateNote failed: $create_resp"
  exit 1
fi

echo "Created note_id=$note_id version=$version"

echo "[2/5] Get note via API..."
get_resp="$(curl -sS "$API_HOST/api/note/$note_id")"
echo "$get_resp" | jq . >/dev/null
get_note_id="$(echo "$get_resp" | jq -r '.noteId // .NoteId // empty')"
if [[ "$get_note_id" != "$note_id" ]]; then
  echo "GetNote mismatch: expected $note_id, got $get_note_id"
  exit 1
fi

echo "[3/5] Save note success path..."
save_resp="$({
  curl -sS -X POST "$API_HOST/api/note/save" \
    -H "Content-Type: application/json" \
    -d "{\"noteId\":\"$note_id\",\"userId\":\"$USER_ID\",\"content\":\"api updated content\",\"expectedVersion\":$version}"
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

echo "[4/5] Save note conflict path..."
conflict_resp="$({
  curl -sS -X POST "$API_HOST/api/note/save" \
    -H "Content-Type: application/json" \
    -d "{\"noteId\":\"$note_id\",\"userId\":\"$USER_ID\",\"content\":\"api conflict content\",\"expectedVersion\":$version}"
} )"
echo "$conflict_resp" | jq . >/dev/null
conflict_code="$(echo "$conflict_resp" | jq -r '.code // .Code // empty')"
if [[ "$conflict_code" != "SAVE_CODE_VERSION_CONFLICT" ]]; then
  echo "Conflict code mismatch, expected SAVE_CODE_VERSION_CONFLICT, got '$conflict_code'"
  echo "Response: $conflict_resp"
  exit 1
fi

echo "[5/5] Get user notes via API..."
list_resp="$(curl -sS "$API_HOST/api/user/$USER_ID/notes")"
echo "$list_resp" | jq . >/dev/null
found="$(echo "$list_resp" | jq -r --arg id "$note_id" '[(.notes // .Notes // [])[]? | select((.noteId // .NoteId // "") == $id)] | length')"
if [[ "$found" == "0" ]]; then
  echo "GetUserNotes did not return created note"
  echo "Response: $list_resp"
  exit 1
fi

echo "API smoke test passed."
