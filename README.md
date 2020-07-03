# Chaos Commander

A chaos test framework. The manager should run on an individual server which has different ansible versions.

## Properties explanation

### Ansible

Ansible is used for deploying TiDB cluster which will run before each test task.

#### Ansible workload

```shell
ansible-playbook local-prepare.yml
ansible-playbook stop.yml
ansible-playbook deploy.yml
ansible-playbook start.yml
```

**Make sure running `ansible-playbook bootstrap.yml` manually before.**

#### Ansible table structure

| Field | Default Value | Description |
| - | - | - |
| id | 1 | ansible ID is `1` |
| path | /home/tidb/ansible/tidb-ansible | Path of ansible which will be used later |

### Load

Load is an executable file which add some presure and checking cluster available and working status. The load program will run after ansible prepare and before chaos job.

#### Load table structure

| Field | Value | Description |
| - | - | - |
| id | 1 | load ID is `1` |
| type | 1 | load type is `1`, stand for executable file |
| path | /home/tidb/load/bank | executable file path |
| param | -addr 172.11.45.14:4000 | the param given to executable file when running |

### Resource

The resource is for job schedule. When there are multiple chaos jobs planned in the same machine, there should be a job execute queue. [TODO]

#### Resource tabel structure

| Field | Value | Description |
| - | - | - |
| id | 1 | load ID is `1` |
| host | 172.191.98.19 | hostname or IP of resource |
| port | 22 | ssh port is `22` |
| cgroup | 2 | the cgroup version of machine, must be `1` or `2` |
| password |  | login password of root user |
| key | /home/tidb/.ssh/id_rsa | private key to login server with root user |

### Blue print

The chaos plan is defined by blue print.

#### Blue print table structure

| Field | Value | Description |
| - | - | - |
| id | 1 | blue print ID is `1` |
| blue_print | {} | JSON format blue print |

#### Content explanation

```json
{
  "hosts": [
    {
      "name": "tikv1",
      "host": "172.16.5.206",
      "ether": "ens18"
    },
    {
      "name": "tikv2",
      "host": "172.16.5.207",
      "ether": "ens18"
    },
    {
      "name": "tikv3",
      "host": "172.16.5.208",
      "ether": "ens18"
    },
    {
      "name": "tidb1",
      "host": "172.16.5.206",
      "ether": "ens18"
    },
    {
      "name": "pd1",
      "host": "172.16.5.206",
      "ether": "ens18"
    }
  ],
  "regions": [
    {
      "name": "tikv",
      "nodes": ["tikv1", "tikv2", "tikv3"]
    },
    {
      "name": "tikv12",
      "nodes": ["tikv1", "tikv2"]
    },
    {
      "name": "tikv23",
      "nodes": ["tikv2", "tikv3"]
    }
  ],
  "relations": [
    {
      "name": "partition_tikv3_down",
      "from": "tikv3",
      "to": "tikv12",
      "rule": "iptables_drop"
    }
  ],
  "exceptions": [
    {
      "name": "tc_delay1",
      "host": "tikv2",
      "rule": "tc_network_delay"
    },
    {
      "name": "io_wbps_limit1",
      "host": "tikv2",
      "rule": "cgroupv2_limit_io_wbps"
    },
    {
      "name": "time_travel_2",
      "host": "tikv2",
      "rule": "time_travel"
    }
  ],
  "steps": [
    {
      "name": "time_travel_2",
      "duration": 60,
      "args": []
    },
    {
      "name": "time_travel_2",
      "duration": 60,
      "args": []
    },
    {
      "name": "time_travel_2",
      "duration": 60,
      "args": []
    },
    {
      "name": "time_travel_2",
      "duration": 60,
      "args": []
    },
    {
      "name": "time_travel_2",
      "duration": 60,
      "args": []
    },
    {
      "name": "io_wbps_limit1",
      "duration": 600,
      "args": ["sda", "tikv", 1000000]
    },
    {
      "name": "partition_tikv3_down",
      "duration": 600,
      "args": []
    }
  ]
}
```

##### blueprint.hosts

Hosts define all components and it's host address.

##### blueprint.regions

Regions define host groups, it's for easier use.

##### blueprint.relations

Relations define relation-type chaos rules.

```json
{
  "name": "partition_tikv3_down",
  "from": "tikv3",
  "to": "tikv12",
  "rule": "iptables_drop"
}
```

`name` field will be called in steps. Same named relations or exceptions will be applied in the same time when called from step.

In this case, `from` field `tikv3` can be found in hosts definition and `to` field `tikv12` can be found in region definition. According to `rule` field with value `iptables_drop`, there will be 2 rules created in step duration:

- tikv3 drop tikv1
- tikv3 drop tikv2

This demonstrate the unilateral network scene, and you can also d this to illustrate a full broken network between 2 regions:

```json
{
  "name": "partition_tikv3_down",
  "from": "tikv3",
  "to": "tikv12",
  "rule": "iptables_drop"
},
{
  "name": "partition_tikv3_down",
  "from": "tikv12",
  "to": "tikv3",
  "rule": "iptables_drop"
}
```

Notice that the two relation rules both named `partition_tikv3_down`, when a step call `partition_tikv3_down`, the 2 relation rules will both be applied, which means 4 rules will be created in the step duration:

- tikv3 drop tikv1
- tikv3 drop tikv2
- tikv1 drop tikv3
- tikv2 drop tikv3

##### blueprint.exceptions

Unlike relation rules will do chaos action between 2 hosts, exceptions is much easier which will only apply a rule in 1 host.

The `name` field of exception is as same as relation's. If the exception shares the same name with a relation, then will be called from step in the same time.

##### blueprint.steps

This is the real chaos plan which will be execute one by one.

There will be 3 sub step in a single step execution:

- Do chaos actions for relations and exception with step's `name` field.
- Sleep for seconds defined by step's `duration` field.
- Clear chaos applied in step1

### Chaos Job

| Field | Value | Description |
| - | - | - |
| name | demo test | the name of the test |
| running_plan | daily | this test will be triggered everyday |
| blue_print | 1 | the chaos operation will follow #1 blue print |
| tidb | master | using master TiDB |
| tikv | master | using master TiKV |
| pd | master | using master PD |
| ansible | 1 | the deployment will follow #1 ansible |
| blue_print | 1 | the chaos operation plan will follow #1 blue print |
| load | 1 | the test program will follow #1 load |
