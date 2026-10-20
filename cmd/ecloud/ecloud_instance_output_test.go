package ecloud

import (
	"testing"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/cli/test"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_instanceAdditionalColumns(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		expected []string
	}{
		{
			name:     "NoFlags",
			flags:    []string{},
			expected: nil,
		},
		{
			name:     "WithTags",
			flags:    []string{"--with-tags"},
			expected: []string{"tags"},
		},
		{
			name:     "WithBackup",
			flags:    []string{"--with-backup"},
			expected: []string{"backup_vm_enabled", "backup_agent_enabled", "backup_agent_installed", "backup_gateway_id"},
		},
		{
			name:     "WithMonitoring",
			flags:    []string{"--with-monitoring"},
			expected: []string{"monitoring_enabled", "monitoring_gateway_id"},
		},
		{
			name:     "WithBackupAndMonitoring",
			flags:    []string{"--with-monitoring", "--with-backup"},
			expected: []string{"backup_vm_enabled", "backup_agent_enabled", "backup_agent_installed", "backup_gateway_id", "monitoring_enabled", "monitoring_gateway_id"},
		},
		{
			name:     "WithAll",
			flags:    []string{"--with-all"},
			expected: []string{"tags", "backup_vm_enabled", "backup_agent_enabled", "backup_agent_installed", "backup_gateway_id", "monitoring_enabled", "monitoring_gateway_id"},
		},
		{
			// Combining --with-all with individual flags should not duplicate columns
			name:     "WithAllAndBackup",
			flags:    []string{"--with-all", "--with-backup"},
			expected: []string{"tags", "backup_vm_enabled", "backup_agent_enabled", "backup_agent_installed", "backup_gateway_id", "monitoring_enabled", "monitoring_gateway_id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ecloudInstanceListCmd(nil)
			assert.Nil(t, cmd.ParseFlags(tt.flags))

			assert.Equal(t, tt.expected, instanceAdditionalColumns(cmd))
		})
	}
}

func Test_instanceOutputOptions(t *testing.T) {
	t.Run("NoFlags_ReturnsNil", func(t *testing.T) {
		cmd := ecloudInstanceListCmd(nil)

		assert.Nil(t, instanceOutputOptions(cmd))
	})

	t.Run("WithBackup_ReturnsOption", func(t *testing.T) {
		cmd := ecloudInstanceListCmd(nil)
		assert.Nil(t, cmd.ParseFlags([]string{"--with-backup"}))

		assert.Len(t, instanceOutputOptions(cmd), 1)
	})
}

func Test_instanceOutputFlags_TableOutput(t *testing.T) {
	instance := ecloud.Instance{
		ID:                  "i-abcdef12",
		Name:                "testinstance",
		BackupEnabled:       true,
		BackupVMEnabled:     true,
		BackupAgentEnabled:  true,
		SecureBackup:        true,
		BackupGatewayID:     "bgw-abcdef12",
		MonitoringEnabled:   true,
		MonitoringGatewayID: "mgw-abcdef12",
	}

	tests := []struct {
		name        string
		flags       []string
		contains    []string
		notContains []string
	}{
		{
			name:        "NoFlags_ExcludesBackupAndMonitoring",
			flags:       []string{},
			contains:    []string{"i-abcdef12"},
			notContains: []string{"BACKUP VM ENABLED", "BACKUP AGENT INSTALLED", "BACKUP GATEWAY ID", "bgw-abcdef12", "MONITORING GATEWAY ID", "mgw-abcdef12"},
		},
		{
			name:        "WithBackup_IncludesBackup",
			flags:       []string{"--with-backup"},
			contains:    []string{"BACKUP VM ENABLED", "BACKUP AGENT ENABLED", "BACKUP AGENT INSTALLED", "BACKUP GATEWAY ID", "bgw-abcdef12"},
			notContains: []string{"| BACKUP ENABLED", "SECURE BACKUP", "MONITORING GATEWAY ID", "mgw-abcdef12"},
		},
		{
			name:        "WithMonitoring_IncludesMonitoring",
			flags:       []string{"--with-monitoring"},
			contains:    []string{"MONITORING ENABLED", "MONITORING GATEWAY ID", "mgw-abcdef12"},
			notContains: []string{"BACKUP GATEWAY ID", "bgw-abcdef12"},
		},
		{
			name:        "WithAll_IncludesAll",
			flags:       []string{"--with-all"},
			contains:    []string{"TAGS", "BACKUP VM ENABLED", "BACKUP AGENT INSTALLED", "BACKUP GATEWAY ID", "bgw-abcdef12", "MONITORING GATEWAY ID", "mgw-abcdef12"},
			notContains: []string{"| BACKUP ENABLED", "SECURE BACKUP"},
		},
	}

	for _, tt := range tests {
		t.Run("InstanceList_"+tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			service := mocks.NewMockECloudService(mockCtrl)
			cmd := ecloudInstanceListCmd(nil)
			assert.Nil(t, cmd.ParseFlags(tt.flags))

			service.EXPECT().GetInstances(gomock.Any()).Return([]ecloud.Instance{instance}, nil)

			out := test.CatchStdOut(t, func() {
				assert.Nil(t, ecloudInstanceList(service, cmd, []string{}))
			})

			assertOutput(t, out, tt.contains, tt.notContains)
		})

		t.Run("InstanceShow_"+tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			service := mocks.NewMockECloudService(mockCtrl)
			cmd := ecloudInstanceShowCmd(nil)
			assert.Nil(t, cmd.ParseFlags(tt.flags))

			service.EXPECT().GetInstance("i-abcdef12").Return(instance, nil)

			out := test.CatchStdOut(t, func() {
				assert.Nil(t, ecloudInstanceShow(service, cmd, []string{"i-abcdef12"}))
			})

			assertOutput(t, out, tt.contains, tt.notContains)
		})

		t.Run("VPCInstanceList_"+tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			service := mocks.NewMockECloudService(mockCtrl)
			cmd := ecloudVPCInstanceListCmd(nil)
			assert.Nil(t, cmd.ParseFlags(tt.flags))

			service.EXPECT().GetVPCInstances("vpc-abcdef12", gomock.Any()).Return([]ecloud.Instance{instance}, nil)

			out := test.CatchStdOut(t, func() {
				assert.Nil(t, ecloudVPCInstanceList(service, cmd, []string{"vpc-abcdef12"}))
			})

			assertOutput(t, out, tt.contains, tt.notContains)
		})
	}
}

func assertOutput(t *testing.T, out string, contains []string, notContains []string) {
	t.Helper()

	for _, s := range contains {
		assert.Contains(t, out, s)
	}
	for _, s := range notContains {
		assert.NotContains(t, out, s)
	}
}

func Test_InstanceCollection_BackupFieldProperties(t *testing.T) {
	instances := InstanceCollection{{ID: "i-abcdef12", BackupEnabled: true, SecureBackup: true}}

	tests := []struct {
		name     string
		property string
		expected string
	}{
		{
			// secure_backup is renamed to backup_agent_installed
			name:     "BackupAgentInstalled_ReturnsSecureBackupValue",
			property: "backup_agent_installed",
			expected: "true\n",
		},
		{
			name:     "SecureBackup_NotOutput",
			property: "secure_backup",
			expected: "\n",
		},
		{
			// backup_enabled is a legacy field and is hidden
			name:     "BackupEnabled_NotOutput",
			property: "backup_enabled",
			expected: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String("output", "", "")
			cmd.Flags().StringSlice("property", nil, "")
			assert.Nil(t, cmd.Flags().Set("output", "value"))
			assert.Nil(t, cmd.Flags().Set("property", tt.property))

			out := test.CatchStdOut(t, func() {
				assert.Nil(t, output.CommandOutput(cmd, instances))
			})

			assert.Equal(t, tt.expected, out)
		})
	}
}
