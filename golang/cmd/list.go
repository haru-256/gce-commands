package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/haru-256/gce-commands/pkg/config"
	"github.com/haru-256/gce-commands/pkg/gce"
	"github.com/haru-256/gce-commands/pkg/log"
	"github.com/haru-256/gce-commands/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var (
	purple = lipgloss.Color("99")
	gray   = lipgloss.Color("#fbfcfc ")

	headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center).Padding(0, 1)
	baseRowStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(gray)
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VM in settings",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		cnf, err := config.ParseConfig(CnfPath)
		if err != nil {
			utils.ErrorReport(fmt.Sprintf("Failed to parse config: %v\n", err))
			os.Exit(1)
		}
		log.Logger.Debug(fmt.Sprintf("Config: %+v", cnf))

		// Update VMs info, such as status and schedule policy
		ctx := context.Background()
		if err = gce.UpdateInstancesInfo(ctx, cnf.VMs); err != nil {
			utils.ErrorReport(fmt.Sprintf("Failed to update VMs info: %v\n", err))
			os.Exit(1)
		}

		// Prepare rows for display
		var rows [][]string
		for _, vm := range cnf.VMs {
			rows = append(rows, []string{
				vm.Name,
				vm.Project,
				vm.Zone,
				vm.MachineType,
				vm.Status,
				vm.SchedulePolicy,
			})
		}

		// render all VMs in settings
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(purple)).
			Headers("Name", "Project", "Zone", "Machine-Type", "Status", "Schedule").
			Rows(rows...).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == 0:
					return headerStyle
				case col == 4: // status
					return baseRowStyle.Align(lipgloss.Center)
				default:
					return baseRowStyle.Align(lipgloss.Left)
				}
			})
		fmt.Println(t)
	},
}
