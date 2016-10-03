* Show short container description. Description should be derived
from the containers REAMDE.md file.
* Fix runtime error:

````
registry:/var/log/supervisord_children# tail -f cntrbrowserd-stderr.log
panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xb code=0x1 addr=0x0 pc=0x5350f3]

goroutine 5 [running]:
panic(0x81cb00, 0xc82000a0b0)
        /usr/lib/go/src/runtime/panic.go:481 +0x3e6
github.com/zaunerc/cntrbrowserd/consul.runCleanUpTask(0x8b7280, 0xe, 0x5)
        /home/admin/gowork/src/github.com/zaunerc/cntrbrowserd/consul/janitor.go:53 +0x9b3
created by github.com/zaunerc/cntrbrowserd/consul.ScheduleCleanUpTask
        /home/admin/gowork/src/github.com/zaunerc/cntrbrowserd/consul/janitor.go:20 +0x133

````
