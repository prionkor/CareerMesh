Run all integration tests:

```sh
APP_ENV=test TEST_DATABASE_URL='postgres://careermesh:careermesh@localhost:5432/careermesh_test?sslmode=disable' go test ./tests/integration/...
```

Run only the new auth lifecycle tests:

```sh
APP_ENV=test TEST_DATABASE_URL='postgres://careermesh:careermesh@localhost:5432/careermesh_test?sslmode=disable' go test ./tests/integration/app -run '^TestAuthLifecycle_' -count=1
```

The auth command selects tests by name; the package’s shared setup still runs migrations against `careermesh_test`.
