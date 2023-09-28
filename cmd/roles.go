package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

type roleMappingResp struct {
	Enabled bool     `json:"enabled"`
	Roles   []string `json:"roles"`
}

var getRolesCmd = &cobra.Command{
	Use:     "roles",
	Aliases: []string{"role"},
	Short:   "get details for a role",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply a role name")
		}
		b, err := getRoles(args[0])
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}

var getRoleMappingCmd = &cobra.Command{
	Use:     "role-mappings [role mapping name]",
	Aliases: []string{"role-map", "role-mapping"},
	Short:   "get details about a role mapping",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply role mapping name")
		}
		b, err := getRoleMapping(args[0])
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}

var listRolesCmd = &cobra.Command{
	Use:     "roles",
	Aliases: []string{"role"},
	Short:   "list roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listRoles()
	},
}

var listRoleMappingCmd = &cobra.Command{
	Use:     "role-mappings",
	Aliases: []string{"role-map", "role-mapping"},
	Short:   "list role mappings",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listRoleMapping()
	},
}

func getRoles(role string) ([]byte, error) {
	resp, err := client.Security.GetRole(client.Security.GetRole.WithPretty(),
		client.Security.GetRole.WithName(role),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getRoleMapping(roleMap string) ([]byte, error) {
	resp, err := client.Security.GetRoleMapping(client.Security.GetRoleMapping.WithPretty(),
		client.Security.GetRoleMapping.WithName(roleMap),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func listRoles() error {
	b, err := getRoles("")
	if err != nil {
		return err
	}
	roles := map[string]any{}
	if err := json.Unmarshal(b, &roles); err != nil {
		return err
	}
	for role := range roles {
		fmt.Println(role)
	}
	return nil
}

func listRoleMapping() error {
	b, err := getRoleMapping("")
	if err != nil {
		return err
	}
	roleMaps := map[string]roleMappingResp{}
	if err := json.Unmarshal(b, &roleMaps); err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "role-mappping\t enabled\t roles\t")
	for name, r := range roleMaps {
		fmt.Fprintf(w, "%s\t %t\t %s\t\n", name, r.Enabled, strings.Join(r.Roles, ","))
	}
	w.Flush()
	return nil
}

func init() {
	getCmd.AddCommand(getRolesCmd, getRoleMappingCmd)
	listCmd.AddCommand(listRolesCmd, listRoleMappingCmd)
}
