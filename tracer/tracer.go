package tracer

import "github.com/getsentry/raven-go"

func Stack(err error) []string {
	return StackPacket(err).Fingerprint
}

func StackPacket(err error) *raven.Packet {
	stacktrace := raven.NewException(err, raven.NewStacktrace(2, 3, nil))
	packet := raven.NewPacket("stack trace", stacktrace)
	return packet
}
