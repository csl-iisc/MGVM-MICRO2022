package tracing

import (
	"fmt"
	"reflect"

	"gitlab.com/akita/akita"
)

// NamedHookable represent something both have a name and can be hooked
type NamedHookable interface {
	akita.Named
	akita.Hookable
	InvokeHook(akita.HookCtx)
}

// A list of hook poses for the hooks to apply to
var (
	HookPosTaskStart = &akita.HookPos{Name: "HookPosTaskStart"}
	HookPosTaskStep  = &akita.HookPos{Name: "HookPosTaskStep"}
	HookPosTaskEnd   = &akita.HookPos{Name: "HookPosTaskEnd"}
)

// StartTask notifies the hooks that hook to the domain about the start of a
// task.
func StartTask(
	id string,
	parentID string,
	now akita.VTimeInSec,
	domain NamedHookable,
	kind string,
	what string,
	detail interface{},
) {
	task := Task{
		ID:        id,
		ParentID:  parentID,
		StartTime: now,
		Kind:      kind,
		What:      what,
		Where:     domain.Name(),
		Detail:    detail,
	}
	ctx := akita.HookCtx{
		Now:    now,
		Domain: domain,
		Item:   task,
		Pos:    HookPosTaskStart,
	}
	domain.InvokeHook(ctx)
}

// AddTaskStep marks that a milestone has been reached when processing a task.
func AddTaskStep(
	id string,
	now akita.VTimeInSec,
	domain NamedHookable,
	what string,
) {
	step := TaskStep{
		Time: now,
		What: what,
	}
	task := Task{
		ID:    id,
		Steps: []TaskStep{step},
	}
	ctx := akita.HookCtx{
		Now:    now,
		Domain: domain,
		Item:   task,
		Pos:    HookPosTaskStep,
	}
	domain.InvokeHook(ctx)
}

// EndTask notifies the hooks about the end of a task.
func EndTask(
	id string,
	now akita.VTimeInSec,
	domain NamedHookable,
) {
	task := Task{
		ID:      id,
		EndTime: now,
		Where:   domain.Name(),
	}
	ctx := akita.HookCtx{
		Now:    now,
		Domain: domain,
		Item:   task,
		Pos:    HookPosTaskEnd,
	}
	domain.InvokeHook(ctx)
}

// MsgIDAtReceiver generates a standard ID for the message task at the
// message receiver.
func MsgIDAtReceiver(msg akita.Msg, domain NamedHookable) string {
	return fmt.Sprintf("%s@%s", msg.Meta().ID, domain.Name())
}

// TraceReqInitiate generatse a new task. The new task has Type="req_out",
// What=[the type name of the message]. This function is to be called by the
// sender of the message.
func TraceReqInitiate(
	msg akita.Msg,
	now akita.VTimeInSec,
	domain NamedHookable,
	taskParentID string,
) {
	StartTask(
		msg.Meta().ID+"_req_out",
		taskParentID,
		now,
		domain,
		"req_out",
		reflect.TypeOf(msg).String(),
		msg,
	)
}

func StartTracingNetworkReq(
	req akita.Msg,
	now akita.VTimeInSec,
	domain NamedHookable,
	parent akita.Msg,
) {
	StartTask(req.Meta().ID+"-trace-trans-req",
		MsgIDAtReceiver(parent, domain),
		now, domain, "trace-trans-req",
		reflect.TypeOf(req).String(), req)
}

func StopTracingNetworkReq(
	req akita.Msg,
	now akita.VTimeInSec,
	domain NamedHookable,
) {
	EndTask(req.Meta().ID+"-trace-trans-req", now, domain)
}

// TraceReqReceive generates a new task for the message handling. The kind of
// the task is always "req_in".
func TraceReqReceive(
	msg akita.Msg,
	now akita.VTimeInSec,
	domain NamedHookable,
) {
	StartTask(
		MsgIDAtReceiver(msg, domain),
		msg.Meta().ID+"_req_out",
		now,
		domain,
		"req_in",
		reflect.TypeOf(msg).String(),
		msg,
	)
}

// TraceReqComplete terminates the message handling task.
func TraceReqComplete(
	msg akita.Msg,
	now akita.VTimeInSec,
	domain NamedHookable,
) {
	EndTask(MsgIDAtReceiver(msg, domain), now, domain)
}

// TraceReqFinalize terminates the message task. This function should be called
// when the sender receives the response.
func TraceReqFinalize(
	msg akita.Msg,
	now akita.VTimeInSec,
	domain NamedHookable,
) {
	EndTask(msg.Meta().ID+"_req_out", now, domain)
}
