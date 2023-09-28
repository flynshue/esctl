package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var getSlmCmd = &cobra.Command{
	Use:   "slm [policy]",
	Short: "get detailed info about snapshot lifecycle management policy (slm)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply a slm policy")
		}
		b, err := getSlm(args[0])
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}

var getSnapshotCmd = &cobra.Command{
	Use:     "snapshots",
	Aliases: []string{"snapshot", "snap", "snaps"},
	Short:   "get detailed information about a snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply a snapshot name")
		}
		return getSnapshot(args[0])
	},
}

var getSnapshotRepoCmd = &cobra.Command{
	Use:     "repository [repo name]",
	Aliases: []string{"repo", "repos"},
	Short:   "get detailed info about snapshot repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply repo name")
		}
		return getSnapshotRepos(args[0])
	},
}

var listSlmCmd = &cobra.Command{
	Use:   "slm",
	Short: "list snapshot lifecycle management policies (slm)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listSlmPolicies()
	},
}

var listSnapshotCmd = &cobra.Command{
	Use:     "snapshots [repo name]",
	Aliases: []string{"snapshot", "snap", "snaps"},
	Short:   "list summary of snapshots for stored in one or more repositories",
	Example: `# list all snapshots for all repositories
esctl list snapshots

# list all snapshots stored under repository
esctl list snapshots test-elastic-fs
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return catSnapshots("*")
		}
		return catSnapshots(args[0])
	},
}

var listSnapshotRepoCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo", "repos"},
	Short:   "list snapshot repositories",
	Example: `# list all snapshot repositories
esctl list snapshot repository
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listSnapshotRepos()
	},
}

type snapshotRepo struct {
	Type string `json:"type"`
}

type catSnapshotResp struct {
	Snapshot string `json:"id"`
	Repo     string `json:"repository"`
	Start    string `json:"start_epoch"`
	End      string `json:"end_epoch"`
	Duration string `json:"duration"`
	Status   string `json:"status"`
}

type slmResp struct {
	LastSuccess `json:"last_success"`
}

type LastSuccess struct {
	Snapshot string `json:"snapshot_name"`
}

func catSnapshots(repo string) error {
	resp, err := client.Cat.Snapshots(client.Cat.Snapshots.WithFormat("json"),
		client.Cat.Snapshots.WithHuman(),
		client.Cat.Snapshots.WithRepository(repo),
		client.Cat.Snapshots.WithS("start_epoch"),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	snapshots := []catSnapshotResp{}
	if err := json.Unmarshal(b, &snapshots); err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "snapshot\t repo\t status\t start\t end\t duration\t")
	for _, s := range snapshots {
		start := parseUnixSecDate(s.Start, localTime)
		end := parseUnixSecDate(s.End, localTime)
		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t\n", s.Snapshot, s.Repo, s.Status, start, end, s.Duration)
	}
	w.Flush()
	return nil
}

func getSlm(policyID string) ([]byte, error) {
	resp, err := client.SlmGetLifecycle(client.SlmGetLifecycle.WithPretty(),
		client.SlmGetLifecycle.WithPolicyID(policyID),
		client.SlmGetLifecycle.WithHuman(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getSnapshotRepos(repo string) error {
	resp, err := client.Snapshot.GetRepository(client.Snapshot.GetRepository.WithPretty(),
		client.Snapshot.GetRepository.WithRepository(repo),
	)
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

func getSnapshot(snapshot string) error {
	resp, err := client.Snapshot.Get("*", []string{snapshot},
		client.Snapshot.Get.WithHuman(),
		client.Snapshot.Get.WithIgnoreUnavailable(true),
		client.Snapshot.Get.WithPretty(),
	)
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

func listSlmPolicies() error {
	b, err := getSlm("")
	if err != nil {
		return err
	}
	policies := map[string]slmResp{}
	if err := json.Unmarshal(b, &policies); err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "slm_policy\t last_snapshot\t")
	for name, p := range policies {
		fmt.Fprintf(w, "%s\t %s\t\n", name, p.Snapshot)
	}
	w.Flush()
	return nil
}

func listSnapshotRepos() error {
	resp, err := client.Snapshot.GetRepository(client.Snapshot.GetRepository.WithHuman())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	repos := map[string]snapshotRepo{}
	if err := json.Unmarshal(b, &repos); err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "repo\t type\t")
	for name, r := range repos {
		fmt.Fprintf(w, "%s\t %s\t\n", name, r.Type)
	}
	w.Flush()
	return nil
}

func init() {
	getCmd.AddCommand(getSnapshotCmd, getSlmCmd)
	getSnapshotCmd.AddCommand(getSnapshotRepoCmd)
	listCmd.AddCommand(listSnapshotCmd, listSlmCmd)
	listSnapshotCmd.AddCommand(listSnapshotRepoCmd)
	listSnapshotCmd.Flags().BoolVar(&localTime, "local", false, "display snapshot start/end times in local time")
}
