package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "show cluster health",
	Long:  `esctl get health`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showHealth()
	},
}

var topRecoveryCmd = &cobra.Command{
	Use:     "recovery",
	Aliases: []string{"recov"},
	Short:   "Watch elasticsearch recovery queue",
	Run: func(cmd *cobra.Command, args []string) {
		topRecovery()
	},
}

func topRecovery() {
	c := make(chan string)
	go watcher(catRecovery)
	go scan(c)
	<-c
}

func watcher(f func() error) {
	for {
		if err := f(); err != nil {
			return
		}
		fmt.Printf("\n\nHit enter to stop\n")
		time.Sleep(time.Second * 2)
		screen.Clear()
	}
}

func scan(c chan string) {
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		c <- in.Text()
	}
}

func catRecovery() error {
	resp, err := client.Cat.Recovery(client.Cat.Recovery.WithActiveOnly(true),
		client.Cat.Recovery.WithBytes("gb"),
		client.Cat.Recovery.WithV(true),
		client.Cat.Recovery.WithS("time:desc,target_node,source_node,index"),
		client.Cat.Recovery.WithH("index,shard,time,type,stage,source_node,target_node,files,files_recovered,files_percent,bytes_total,bytes_percent,translog_ops_recovered,translog_ops,translog_ops_percent"),
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

func showHealth() error {
	resp, err := client.Cluster.Health(client.Cluster.Health.WithPretty())
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

func init() {
	getCmd.AddCommand(healthCmd)
	topCmd.AddCommand(topRecoveryCmd)
}
