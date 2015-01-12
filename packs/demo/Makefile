default:
	# run two instances to show the hello and goodbye messages
	go build demo.go
	./demo -logtostderr &
	sleep 1
	./demo -logtostderr &
	rm demo
