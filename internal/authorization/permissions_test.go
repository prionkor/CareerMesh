package authorization

import "testing"

func TestHasPermission_ExactMatch(t *testing.T) {
	granted := []string{PermProfilesUpdateOwn}

	if !HasPermission(granted, PermProfilesUpdateOwn) {
		t.Fatal("expected exact match to satisfy the required permission")
	}
}

func TestHasPermission_ResourceWildcard(t *testing.T) {
	granted := []string{"*:update:own"}

	if !HasPermission(granted, PermProfilesUpdateOwn) {
		t.Fatal("expected resource wildcard to satisfy the required permission")
	}
}

func TestHasPermission_ActionWildcard(t *testing.T) {
	granted := []string{PermProfilesAnyOwn}

	if !HasPermission(granted, PermProfilesUpdateOwn) {
		t.Fatal("expected action wildcard to satisfy the required permission")
	}
}

func TestHasPermission_ResourceAndActionWildcard(t *testing.T) {
	granted := []string{"*:*:own"}

	if !HasPermission(granted, PermProfilesUpdateOwn) {
		t.Fatal("expected resource+action wildcard to satisfy the required permission")
	}
}

func TestHasPermission_AllScopeSatisfiesOwn(t *testing.T) {
	granted := []string{PermProfilesUpdateAll}

	if !HasPermission(granted, PermProfilesUpdateOwn) {
		t.Fatal("expected an 'all' scoped permission to satisfy a required 'own' scope")
	}
}

func TestHasPermission_OwnScopeDoesNotSatisfyAll(t *testing.T) {
	granted := []string{PermProfilesUpdateOwn}

	if HasPermission(granted, PermProfilesUpdateAll) {
		t.Fatal("expected an 'own' scoped permission to NOT satisfy a required 'all' scope")
	}
}

func TestHasPermission_UnrelatedResourceOrActionDoesNotMatch(t *testing.T) {
	granted := []string{PermProjectsUpdateAll, PermProfilesReadAll}

	if HasPermission(granted, PermProfilesUpdateOwn) {
		t.Fatal("expected unrelated resource/action permissions to not satisfy the required permission")
	}
}

func TestHasPermission_EmptyPermissions(t *testing.T) {
	if HasPermission(nil, PermProfilesUpdateOwn) {
		t.Fatal("expected no permissions to never satisfy a required permission")
	}
	if HasPermission([]string{}, PermProfilesUpdateOwn) {
		t.Fatal("expected an empty permission slice to never satisfy a required permission")
	}
}

func TestHasPermission_GlobalWildcardReadAllSatisfiesResourceReadAll(t *testing.T) {
	granted := []string{PermAllReadAll}

	if !HasPermission(granted, PermProfilesReadAll) {
		t.Fatal("expected *:read:all to satisfy profiles:read:all")
	}
}

func TestHasPermission_GlobalWildcardReadAllSatisfiesResourceReadOwn(t *testing.T) {
	granted := []string{PermAllReadAll}

	if !HasPermission(granted, PermProfilesReadOwn) {
		t.Fatal("expected *:read:all to satisfy profiles:read:own")
	}
}

func TestHasPermission_GlobalWildcardReadAllDoesNotSatisfyUpdate(t *testing.T) {
	granted := []string{PermAllReadAll}

	if HasPermission(granted, PermProfilesUpdateOwn) {
		t.Fatal("expected *:read:all to NOT satisfy profiles:update:own")
	}
}

func TestHasPermission_GlobalWildcardAnyAllSatisfiesResourceUpdateAll(t *testing.T) {
	granted := []string{PermAllAnyAll}

	if !HasPermission(granted, PermProfilesUpdateAll) {
		t.Fatal("expected *:*:all to satisfy profiles:update:all")
	}
}

func TestHasPermission_GlobalWildcardAnyAllSatisfiesResourceUpdateOwn(t *testing.T) {
	granted := []string{PermAllAnyAll}

	if !HasPermission(granted, PermProfilesUpdateOwn) {
		t.Fatal("expected *:*:all to satisfy profiles:update:own")
	}
}

func TestHasPermission_GlobalWildcardAnyAllSatisfiesAnyResourceAndAction(t *testing.T) {
	granted := []string{PermAllAnyAll}

	if !HasPermission(granted, PermUsersDeleteAll) {
		t.Fatal("expected *:*:all to satisfy users:delete:all")
	}
}
