package main_test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	g "github.com/tsanton/goflake-client/goflake"
	i "github.com/tsanton/goflake-client/goflake/integration"
	a "github.com/tsanton/goflake-client/goflake/models/assets"
	d "github.com/tsanton/goflake-client/goflake/models/describables"
	dg "github.com/tsanton/goflake-client/goflake/models/describables/grants"
	e "github.com/tsanton/goflake-client/goflake/models/entities"
	eg "github.com/tsanton/goflake-client/goflake/models/entities/grants"
	"github.com/tsanton/goflake-client/goflake/models/enums"
	u "github.com/tsanton/goflake-client/goflake/utilities"
)

func TestMain(m *testing.M) {
	/* Setup */
	t := &testing.T{}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		PluginDir:    "/workspace/.providercache",
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"snowflake_uid":     os.Getenv("SNOWFLAKE_UID"),
			"snowflake_pwd":     os.Getenv("SNOWFLAKE_PWD"),
			"snowflake_account": os.Getenv("SNOWFLAKE_ACCOUNT"),
			"snowflake_region":  os.Getenv("SNOWFLAKE_REGION"),
			"db_prefix":         "INTEGRATION",
			"db_identifier":     "TEST_DB",
			"schema_name":       "TESTING_SCHEMA",
		},
		EnvVars: map[string]string{},
		NoColor: true,
	})
	defer terraform.Destroy(t, terraformOptions)
	if _, err := terraform.InitAndApplyE(t, terraformOptions); err != nil {
		os.Exit(1)
	}
	/* Run */
	// exitCode := m.Run()
	m.Run()

	/* Teardown: if you os.exit() defer function are not run */
	// os.Exit(exitCode)
}

func Goflake() (*g.GoflakeClient, u.Stack[a.ISnowflakeAsset]) {
	stack := u.Stack[a.ISnowflakeAsset]{}
	cli := g.GoflakeClient{
		SnowflakeHost: os.Getenv("SNOWFLAKE_HOST"),
		SnowflakeUid:  os.Getenv("SNOWFLAKE_UID"),
		SnowflakePwd:  os.Getenv("SNOWFLAKE_PWD"),
		SnowflakeRole: os.Getenv("SNOWFLAKE_ROLE"),
		SnowflakeWh:   os.Getenv("SNOWFLAKE_WH"),
	}

	cli.Open()

	return &cli, stack
}

func Test_schema_exists(t *testing.T) {
	/* Arrange */
	a := assert.New(t)
	cli, _ := Goflake()

	/* Act */
	sch, err := g.Describe[e.Schema](cli, &d.Schema{DatabaseName: "INTEGRATION_TEST_DB", SchemaName: "TESTING_SCHEMA"})
	i.ErrorFailNow(t, err)

	/* Assert */
	a.Nil(err)
	a.Equal(sch.Name, "TESTING_SCHEMA")
	a.Equal(sch.DatabaseName, "INTEGRATION_TEST_DB")
	a.Equal(sch.Owner, "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RWC")
}

func Test_schema_access_role_r_exists(t *testing.T) {
	/* Arrange */
	a := assert.New(t)
	cli, _ := Goflake()

	/* Act */
	role, err := g.Describe[e.Role](cli, &d.Role{Name: "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_R"})

	/* Assert */
	a.Nil(err)
	a.Equal(role.Name, "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_R")
	a.Equal(role.Owner, "INTEGRATION_TEST_DB_USER_ADMIN")
}

func Test_schema_access_role_r_privileges(t *testing.T) {
	/* Arrange */
	a := assert.New(t)
	cli, _ := Goflake()

	/* Act */
	res, err := g.Describe[*eg.RoleGrants](cli, &dg.RoleGrant{RoleName: "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_R"})

	/* Assert */
	a.Nil(err)
	assert.Len(t, res.Grants, 2)

	dbUsage, ok := lo.Find(res.Grants, func(i eg.RoleGrant) bool {
		return i.Privilege == enums.PrivilegeUsage && i.GrantedOn == enums.SnowflakeObjectDatabase && i.GrantTargetName == "INTEGRATION_TEST_DB"
	})
	assert.True(t, ok)
	assert.Equal(t, "INTEGRATION_TEST_DB_SYS_ADMIN", dbUsage.GrantedBy)

	schemaUsage, ok := lo.Find(res.Grants, func(i eg.RoleGrant) bool {
		return i.Privilege == enums.PrivilegeUsage && i.GrantedOn == enums.SnowflakeObjectSchema && i.GrantTargetName == "INTEGRATION_TEST_DB.TESTING_SCHEMA"
	})
	assert.True(t, ok)
	assert.Equal(t, "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RWC", schemaUsage.GrantedBy)
}

func Test_schema_access_role_rw_exists(t *testing.T) {
	/* Arrange */
	a := assert.New(t)
	cli, _ := Goflake()

	/* Act */
	role, err := g.Describe[e.Role](cli, &d.Role{Name: "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RW"})

	/* Assert */
	a.Nil(err)
	a.Equal(role.Name, "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RW")
	a.Equal(role.Owner, "INTEGRATION_TEST_DB_USER_ADMIN")
}

func Test_schema_access_role_rw_privileges(t *testing.T) {
	/* Arrange */
	a := assert.New(t)
	cli, _ := Goflake()

	/* Act */
	res, err := g.Describe[*eg.RoleGrants](cli, &dg.RoleGrant{RoleName: "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RW"})

	/* Assert */
	a.Nil(err)
	assert.Len(t, res.Grants, 3) //db usage, schema usage, access_role_r usage

	dbUsage, ok := lo.Find(res.Grants, func(i eg.RoleGrant) bool {
		return i.Privilege == enums.PrivilegeUsage && i.GrantedOn == enums.SnowflakeObjectDatabase && i.GrantTargetName == "INTEGRATION_TEST_DB"
	})
	assert.True(t, ok)
	assert.Equal(t, "INTEGRATION_TEST_DB_SYS_ADMIN", dbUsage.GrantedBy)

	schemaUsage, ok := lo.Find(res.Grants, func(i eg.RoleGrant) bool {
		return i.Privilege == enums.PrivilegeUsage && i.GrantedOn == enums.SnowflakeObjectSchema && i.GrantTargetName == "INTEGRATION_TEST_DB.TESTING_SCHEMA"
	})
	assert.True(t, ok)
	assert.Equal(t, "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RWC", schemaUsage.GrantedBy)
}

func Test_schema_access_role_rwc_exists(t *testing.T) {
	/* Arrange */
	a := assert.New(t)
	cli, _ := Goflake()

	/* Act */
	role, err := g.Describe[e.Role](cli, &d.Role{Name: "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RW"})

	/* Assert */
	a.Nil(err)
	a.Equal(role.Name, "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RW")
	a.Equal(role.Owner, "INTEGRATION_TEST_DB_USER_ADMIN")
}

func Test_schema_access_role_rwc_privileges(t *testing.T) {
	/* Arrange */
	a := assert.New(t)
	cli, _ := Goflake()

	/* Act */
	res, err := g.Describe[*eg.RoleGrants](cli, &dg.RoleGrant{RoleName: "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RWC"})

	/* Assert */
	a.Nil(err)
	assert.Len(t, res.Grants, 4) //db usage, schema ownership, schema usage, access_role_rw usage

	dbUsage, ok := lo.Find(res.Grants, func(i eg.RoleGrant) bool {
		return i.Privilege == enums.PrivilegeUsage && i.GrantedOn == enums.SnowflakeObjectDatabase && i.GrantTargetName == "INTEGRATION_TEST_DB"
	})
	assert.True(t, ok)
	assert.Equal(t, "INTEGRATION_TEST_DB_SYS_ADMIN", dbUsage.GrantedBy)

	schemaOwnership, ok := lo.Find(res.Grants, func(i eg.RoleGrant) bool {
		return i.Privilege == enums.PrivilegeOwnership && i.GrantedOn == enums.SnowflakeObjectSchema && i.GrantTargetName == "INTEGRATION_TEST_DB.TESTING_SCHEMA"
	})
	assert.True(t, ok)
	assert.Equal(t, "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RWC", schemaOwnership.GrantedBy)

	schemaUsage, ok := lo.Find(res.Grants, func(i eg.RoleGrant) bool {
		return i.Privilege == enums.PrivilegeUsage && i.GrantedOn == enums.SnowflakeObjectSchema && i.GrantTargetName == "INTEGRATION_TEST_DB.TESTING_SCHEMA"
	})
	assert.True(t, ok)
	assert.Equal(t, "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RWC", schemaUsage.GrantedBy)
}

// Test that role hierarchy is intact: access_role_r -> access_role_rw -> access_role_rwc -> db_sys_admin -> sysadmin -> accountadmin
func Test_schema_access_role_hierarchy(t *testing.T) {
	/* Arrange */
	a := assert.New(t)
	cli, _ := Goflake()

	/* Act */
	hierarchy, err := g.Describe[e.RoleHierarchy](cli, &d.RoleHierarchy{RoleName: "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_R"})

	/* Assert */
	a.Nil(err)
	rw, ok := lo.Find(hierarchy.InheritingRoles, func(i e.InheritedRole) bool {
		return i.ParentRoleName == "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RW"
	})
	a.True(ok)
	a.Equal(rw.DistanceFromSource, 0)

	rwc, ok := lo.Find(hierarchy.InheritingRoles, func(i e.InheritedRole) bool {
		return i.ParentRoleName == "INTEGRATION_TEST_DB_TESTING_SCHEMA_ACCESS_ROLE_RWC"
	})
	a.True(ok)
	a.Equal(rwc.DistanceFromSource, 1)

	dbsa, ok := lo.Find(hierarchy.InheritingRoles, func(i e.InheritedRole) bool { return i.ParentRoleName == "INTEGRATION_TEST_DB_SYS_ADMIN" })
	a.True(ok)
	a.Equal(dbsa.DistanceFromSource, 2)

	sa, ok := lo.Find(hierarchy.InheritingRoles, func(i e.InheritedRole) bool { return i.ParentRoleName == "SYSADMIN" })
	a.True(ok)
	a.Equal(sa.DistanceFromSource, 3)

	aa, ok := lo.Find(hierarchy.InheritingRoles, func(i e.InheritedRole) bool { return i.ParentRoleName == "ACCOUNTADMIN" })
	a.True(ok)
	a.Equal(aa.DistanceFromSource, 4)
}
