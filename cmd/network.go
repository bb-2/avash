package cmd

import (
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/network"
	"github.com/spf13/cobra"
)

// NetworkCommand represents the network command
var NetworkCommand = &cobra.Command{
	Use:   "network",
	Short: "Tools for interacting with remote hosts.",
	Long:  `Tools for interacting with remote hosts.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// NetworkSSHCommand represents the network ssh command
var NetworkSSHCommand = &cobra.Command{
	Use: "ssh",
	Short: "Tools for interacting with remote hosts via SSH.",
	Long:  `Tools for interacting with remote hosts via SSH.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// SSHDeployCommand deploys a node through an SSH client
var SSHDeployCommand = &cobra.Command{
	Use: "deploy [node name] [SSH username] [IP address]",
	Short: "Deploys a remotely running node.",
	Long:  `Deploys a remotely running node to a specified host.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		c1, err := network.NewSSH(args[1], args[2])
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer c1.Close()
		if err := initHost(c1); err != nil {
			log.Error(err.Error())
			return
		}
		// New connection necessary to refresh user groups
		c2, err := c1.Clone()
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer c2.Close()
		if err := deploy(c2); err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Node successfully deployed!")
	},
}

func initHost(client *network.SSHClient) error {
	const cfp string = "./init.sh"
	cmds := []string{
		"chmod 777 " + cfp,
		cfp,
	}

	if err := client.CopyFile("network/init.sh", cfp); err != nil {
		return err
	}
	defer client.Remove(cfp)

	if err := client.Run(cmds); err != nil {
		return err
	}
	return nil
}

func deploy(client *network.SSHClient) error {
	const cfp string = "./startnode.sh"
	cmds := []string{
		"chmod 777 " + cfp,
		cfp + " --staking-tls-enabled=false",
	}

	if err := client.CopyFile("network/startnode.sh", cfp); err != nil {
		return err
	}
	defer client.Remove(cfp)

	if err := client.Run(cmds); err != nil {
		return err
	}
	return nil
}

func init() {
	NetworkSSHCommand.AddCommand(SSHDeployCommand)
	NetworkCommand.AddCommand(NetworkSSHCommand)
}