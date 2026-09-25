// Package authorization provides simple RBAC permission checks. Resource
// ownership is enforced separately in each domain's service layer.
package authorization

import (
	"errors"
	"strings"

	"uuid"
)

// ErrForbidden indicates the authenticated user does not own the resource
// and does not have a permission scoped to "all".
var ErrForbidden = errors.New("forbidden: resource does not belong to the authenticated user")

// Permission keys use the form "<resource>:<action>:<scope>", where scope is
// "own" (the authenticated user's own resources) or "all" (any resource).
const (
	PermUsersReadOwn   = "users:read:own"
	PermUsersReadAll   = "users:read:all"
	PermUsersCreateOwn = "users:create:own"
	PermUsersCreateAll = "users:create:all"
	PermUsersUpdateOwn = "users:update:own"
	PermUsersUpdateAll = "users:update:all"
	PermUsersDeleteOwn = "users:delete:own"
	PermUsersDeleteAll = "users:delete:all"
	PermUsersAnyOwn    = "users:*:own"
	PermUsersAnyAll    = "users:*:all"

	PermProfilesReadOwn   = "profiles:read:own"
	PermProfilesReadAll   = "profiles:read:all"
	PermProfilesCreateOwn = "profiles:create:own"
	PermProfilesCreateAll = "profiles:create:all"
	PermProfilesUpdateOwn = "profiles:update:own"
	PermProfilesUpdateAll = "profiles:update:all"
	PermProfilesDeleteOwn = "profiles:delete:own"
	PermProfilesDeleteAll = "profiles:delete:all"
	PermProfilesAnyOwn    = "profiles:*:own"
	PermProfilesAnyAll    = "profiles:*:all"

	PermExperiencesReadOwn   = "experiences:read:own"
	PermExperiencesReadAll   = "experiences:read:all"
	PermExperiencesCreateOwn = "experiences:create:own"
	PermExperiencesCreateAll = "experiences:create:all"
	PermExperiencesUpdateOwn = "experiences:update:own"
	PermExperiencesUpdateAll = "experiences:update:all"
	PermExperiencesDeleteOwn = "experiences:delete:own"
	PermExperiencesDeleteAll = "experiences:delete:all"
	PermExperiencesAnyOwn    = "experiences:*:own"
	PermExperiencesAnyAll    = "experiences:*:all"

	PermProjectsReadOwn   = "projects:read:own"
	PermProjectsReadAll   = "projects:read:all"
	PermProjectsCreateOwn = "projects:create:own"
	PermProjectsCreateAll = "projects:create:all"
	PermProjectsUpdateOwn = "projects:update:own"
	PermProjectsUpdateAll = "projects:update:all"
	PermProjectsDeleteOwn = "projects:delete:own"
	PermProjectsDeleteAll = "projects:delete:all"
	PermProjectsAnyOwn    = "projects:*:own"
	PermProjectsAnyAll    = "projects:*:all"

	PermSkillsReadOwn   = "skills:read:own"
	PermSkillsReadAll   = "skills:read:all"
	PermSkillsCreateOwn = "skills:create:own"
	PermSkillsCreateAll = "skills:create:all"
	PermSkillsUpdateOwn = "skills:update:own"
	PermSkillsUpdateAll = "skills:update:all"
	PermSkillsDeleteOwn = "skills:delete:own"
	PermSkillsDeleteAll = "skills:delete:all"
	PermSkillsAnyOwn    = "skills:*:own"
	PermSkillsAnyAll    = "skills:*:all"

	PermEducationReadOwn   = "education:read:own"
	PermEducationReadAll   = "education:read:all"
	PermEducationCreateOwn = "education:create:own"
	PermEducationCreateAll = "education:create:all"
	PermEducationUpdateOwn = "education:update:own"
	PermEducationUpdateAll = "education:update:all"
	PermEducationDeleteOwn = "education:delete:own"
	PermEducationDeleteAll = "education:delete:all"
	PermEducationAnyOwn    = "education:*:own"
	PermEducationAnyAll    = "education:*:all"

	PermCertificationsReadOwn   = "certifications:read:own"
	PermCertificationsReadAll   = "certifications:read:all"
	PermCertificationsCreateOwn = "certifications:create:own"
	PermCertificationsCreateAll = "certifications:create:all"
	PermCertificationsUpdateOwn = "certifications:update:own"
	PermCertificationsUpdateAll = "certifications:update:all"
	PermCertificationsDeleteOwn = "certifications:delete:own"
	PermCertificationsDeleteAll = "certifications:delete:all"
	PermCertificationsAnyOwn    = "certifications:*:own"
	PermCertificationsAnyAll    = "certifications:*:all"

	PermLanguagesReadOwn   = "languages:read:own"
	PermLanguagesReadAll   = "languages:read:all"
	PermLanguagesCreateOwn = "languages:create:own"
	PermLanguagesCreateAll = "languages:create:all"
	PermLanguagesUpdateOwn = "languages:update:own"
	PermLanguagesUpdateAll = "languages:update:all"
	PermLanguagesDeleteOwn = "languages:delete:own"
	PermLanguagesDeleteAll = "languages:delete:all"
	PermLanguagesAnyOwn    = "languages:*:own"
	PermLanguagesAnyAll    = "languages:*:all"

	// PermAll* grant permissions across every resource via the "*" wildcard.
	PermAllReadOwn = "*:read:own"
	PermAllReadAll = "*:read:all"
	PermAllAnyOwn  = "*:*:own"
	PermAllAnyAll  = "*:*:all"
)

// HasPermission reports whether permissions contains a permission that
// satisfies required. Permissions have the form "<resource>:<action>:<scope>".
// "*" in the resource or action segment of a granted permission matches any
// value in that position. A granted scope of "all" satisfies a required
// scope of "own" or "all"; a granted scope of "own" only satisfies a
// required scope of "own".
func HasPermission(permissions []string, required string) bool {
	for _, p := range permissions {
		if permissionSatisfies(p, required) {
			return true
		}
	}

	return false
}

// permissionSatisfies reports whether granted authorizes required.
func permissionSatisfies(granted, required string) bool {
	gResource, gAction, gScope, ok := splitPermission(granted)
	if !ok {
		return false
	}

	rResource, rAction, rScope, ok := splitPermission(required)
	if !ok {
		return false
	}

	if gResource != "*" && gResource != rResource {
		return false
	}
	if gAction != "*" && gAction != rAction {
		return false
	}

	return gScope == rScope || (gScope == "all" && rScope == "own")
}

// splitPermission splits a permission string into its resource, action, and
// scope segments. ok is false if the string is not exactly three segments.
func splitPermission(permission string) (resource, action, scope string, ok bool) {
	parts := strings.SplitN(permission, ":", 3)
	if len(parts) != 3 {
		return "", "", "", false
	}

	return parts[0], parts[1], parts[2], true
}

func CanAccessResource(
	userID uuid.UUID,
	ownerID uuid.UUID,
	permissions []string,
	ownPermission string,
	allPermission string,
) bool {
	if ownerID == userID {
		return HasPermission(permissions, ownPermission)
	}

	return HasPermission(permissions, allPermission)
}
