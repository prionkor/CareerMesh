#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if [[ ! -f .env.test ]]; then
	echo "Error: .env.test not found" >&2
	exit 1
fi

set -a
source .env.test
set +a

if [[ "${APP_ENV:-}" != "test" ]]; then
	echo "Error: APP_ENV must be 'test'" >&2
	exit 1
fi

if [[ "${TEST_DATABASE_URL:-}" != *"/careermesh_test?"* ]]; then
	echo "Error: TEST_DATABASE_URL must target careermesh_test" >&2
	exit 1
fi

go test ./tests/integration/app -count=1
