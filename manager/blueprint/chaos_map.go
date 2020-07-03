package blueprint

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/manager/blueprint/client"
	"github.com/you06/chaos-commander/manager/executor/content"
	"github.com/you06/chaos-commander/pkg/types"
)

type relationFunc func(ctn *content.Content, resource *types.Resource, addrStr string, args... interface{}) error
type exceptionFunc func(ctn *content.Content, resource *types.Resource, args... interface{}) error

type relationCaller struct {
	Fn relationFunc
	Arg int
}
type exceptionCaller struct {
	Fn exceptionFunc
	Arg int
}

type relationRule struct {
	Do *relationCaller
	Clear *relationCaller
}
type exceptionRule struct {
	Do *exceptionCaller
	Clear *exceptionCaller
}

type relationMap map[string]*relationRule
type exceptionMap map[string]*exceptionRule

// Call execute the function
func (r *relationCaller)Call(ctn *content.Content, resource *types.Resource, addrStr string, args []interface{}) error {
	if len(args) < r.Arg {
		return errors.New("not enough args")
	}
	return errors.Trace(r.Fn(ctn, resource, addrStr, args[:r.Arg]...))
}

// Call execute the function
func (e *exceptionCaller)Call(ctn *content.Content, resource *types.Resource, args []interface{}) error {
	if len(args) < e.Arg {
		return errors.New("not enough args")
	}
	return errors.Trace(e.Fn(ctn, resource, args[:e.Arg]...))
}

func initRelationMap() relationMap {
	var m relationMap = make(map[string]*relationRule)

	m["iptables_drop"] = &relationRule{
		Do: &relationCaller{
			Fn: client.IpTablesDrop,
			Arg: 0,
		},
		Clear: &relationCaller{
			Fn: client.IpTablesDropClear,
			Arg: 0,
		},
	}

	m["iptables_reject"] = &relationRule{
		Do: &relationCaller{
			Fn: client.IpTablesReject,
			Arg: 0,
		},
		Clear: &relationCaller{
			Fn: client.IpTablesRejectClear,
			Arg: 0,
		},
	}

	return m
}

func initExceptionMap() exceptionMap {
	var m exceptionMap = make(map[string]*exceptionRule)

	m["cgroupv2_limit_io_wbps"] = &exceptionRule{
		Do: &exceptionCaller{
			Fn: client.CGroupV2ioWbps,
			Arg: 3,
		},
		Clear: &exceptionCaller{
			Fn: client.CGroupV2ioClear,
			Arg: 1,
		},
	}

	m["sudden_death"] = &exceptionRule{
		Do: &exceptionCaller{
			Fn: client.SuddenDeathDo,
			Arg: 0,
		},
		Clear: &exceptionCaller{
			Fn:client.SuddenDeathClear,
			Arg: 0,
		},
	}

	m["time_travel"] = &exceptionRule{
		Do: &exceptionCaller{
			Fn: client.TravelTimeDo,
			Arg: 0,
		},
		Clear: &exceptionCaller{
			Fn: client.TravelTimeClear,
			Arg: 0,
		},
	}

	return m
}
