package cmd

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

var (
	engUsers = []string{"test-eng-01", "test-eng-02"}
	userCfg  = `{
		"password":"p@ssw0rd",
	}`
	roleMapName = "test-admin"
	roleMapping = `{
		"roles": ["superuser"],
		"enabled": true,
		"rules": {
		   "field" : { "username" : "test-eng-*" }
		}
	  }`
)

func TestRoles_ListRoles(t *testing.T) {
	if err := listRoles(); err != nil {
		t.Error(err)
	}
}

func TestRoles_GetRoles(t *testing.T) {
	b, err := getRoles("watcher_admin")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

func TestRoles_ListRoleMapping(t *testing.T) {
	if err := roleMappingSetup(); err != nil {
		t.Error(err)
	}
	if err := listRoleMapping(); err != nil {
		t.Error(err)
	}
	b, err := getRoleMapping(roleMapName)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

func roleMappingSetup() error {
	// delete if they already exist
	for _, user := range engUsers {
		client.Security.DeleteUser(user)
	}
	buf := bytes.NewBufferString(userCfg)
	for _, user := range engUsers {
		client.Security.PutUser(user, buf)
	}
	// delete role mapping if it already exists
	client.Security.DeleteRoleMapping(roleMapName)
	buf = bytes.NewBufferString(roleMapping)
	resp, err := client.Security.PutRoleMapping(roleMapName, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
