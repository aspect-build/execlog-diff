package bazel8

import (
	"github.com/aspect-build/execlog-diff/third_party/github.com/bazelbuild/bazel8/src/main/protobuf/spawn"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func has(message protoreflect.ProtoMessage, descriptor protoreflect.FieldDescriptor) bool {
	return message.ProtoReflect().Has(descriptor)
}

func descriptor(message interface{}, name protoreflect.Name) protoreflect.FieldDescriptor {
	switch v := message.(type) {
	case *spawn.ExecLogEntry:
		return v.ProtoReflect().Descriptor().Fields().ByName(name)
	case *protoreflect.Message:
		mes := *v
		return mes.Descriptor().Fields().ByName(name)
	case *protoreflect.Value:
		val := *v
		mes := val.Message()
		return descriptor(&mes, name)
	default:
		return nil
	}
}

func get(val protoreflect.Value, name protoreflect.Name) protoreflect.Value {
	message := val.Message()
	descriptor := descriptor(&message, name)
	return message.Get(descriptor)
}
