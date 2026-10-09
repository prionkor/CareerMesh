package authorization

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestPermissionsForRoles(t *testing.T) {
	resources := []string{"users", "profiles", "experiences", "projects", "skills", "education", "certifications", "languages"}
	actions := []string{"read", "create", "update", "delete"}

	expectedForScope := func(scopes ...string) []string {
		permissions := make([]string, 0, len(resources)*len(actions)*len(scopes)+2)
		for _, resource := range resources {
			for _, action := range actions {
				for _, scope := range scopes {
					permissions = append(permissions, resource+":"+action+":"+scope)
				}
			}
		}
		sort.Strings(permissions)
		return permissions
	}

	userExpected := expectedForScope("own")
	adminExpected := expectedForScope("own", "all")
	adminExpected = removePermission(adminExpected, PermUsersCreateAll, PermUsersDeleteAll)
	superadminExpected := append(append(expectedForScope("own", "all"), PermRoleUpdateAll), PermUserRolesCreateAll, PermUserRolesUpdateAll)
	superadminExpected = removePermission(superadminExpected, PermUsersCreateAll, PermUsersDeleteAll)
	sort.Strings(superadminExpected)

	tests := []struct {
		name  string
		roles []Role
		want  []string
	}{
		{name: "user", roles: []Role{RoleUser}, want: userExpected},
		{name: "admin", roles: []Role{RoleAdmin}, want: adminExpected},
		{name: "superadmin", roles: []Role{RoleSuperadmin}, want: superadminExpected},
		{name: "unknown", roles: []Role{"unknown"}, want: []string{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := PermissionsForRoles(test.roles)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("PermissionsForRoles(%v) = %v, want %v", test.roles, got, test.want)
			}
		})
	}

	userPermissions := PermissionsForRoles([]Role{RoleUser})
	for _, permission := range userPermissions {
		if strings.HasSuffix(permission, ":all") || strings.HasPrefix(permission, "user_roles:") {
			t.Errorf("regular user received elevated permission %q", permission)
		}
	}
	if HasPermission(PermissionsForRoles([]Role{RoleAdmin}), PermRoleUpdateAll) {
		t.Error("admin received role update permission")
	}
	if !HasPermission(PermissionsForRoles([]Role{RoleSuperadmin}), PermRoleUpdateAll) {
		t.Error("superadmin did not receive role update permission")
	}

	first := PermissionsForRoles([]Role{RoleSuperadmin, RoleUser, RoleSuperadmin})
	second := PermissionsForRoles([]Role{RoleUser, RoleSuperadmin})
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("role order or duplicates changed permissions: %v != %v", first, second)
	}
}

func removePermission(permissions []string, excluded ...string) []string {
	filtered := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		remove := false
		for _, excludedPermission := range excluded {
			if permission == excludedPermission {
				remove = true
				break
			}
		}
		if !remove {
			filtered = append(filtered, permission)
		}
	}
	sort.Strings(filtered)
	return filtered
}

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
