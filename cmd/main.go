package main

import (
	"clean-architecture/initiator"
	"context"
)

func main() {
	initiator.Initialize(context.Background())
}
