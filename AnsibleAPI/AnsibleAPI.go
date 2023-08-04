package AnsibleAPI

import (
	"context"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
)

const INVENTORY_PATH string = "./Inventory/"
const ANSIBLE_PLAYBOOKS = "AnsiblePlaybooks/"
const CONNECTION_TYPE string = "ssh"
const ROCKPI_USER = "rock"

func ExecuteAnsible(playbookFilename string) {

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: CONNECTION_TYPE,
		User:       ROCKPI_USER,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: INVENTORY_PATH,
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{INVENTORY_PATH + ANSIBLE_PLAYBOOKS + playbookFilename},
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		//PrivilegeEscalationOptions: ansiblePlaybookPrivilegeEscalationOptions,
		Options: ansiblePlaybookOptions,
		Exec: execute.NewDefaultExecute(
			//execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithTransformers(
				results.Prepend("Go-ansible example with become"),
			),
		),
	}

	err := playbook.Run(context.TODO())
	if err != nil {
		panic(err)
	}

}
