package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gofrs/uuid"
)

type IKGenerator struct {
	Snowflake *snowflake.Node
}

func NewIKGenerator() (g *IKGenerator, err error) {
	g = &IKGenerator{}

	g.Snowflake, err = snowflake.NewNode(1)
	if err != nil {
		//errors.New("Error creating snowflake node")
		return
	}

	return
}

func (g *IKGenerator) genSnowflakeID() int64 {
	id := g.Snowflake.Generate()
	return id.Int64()
}

func (g *IKGenerator) genUuidV1() string {
	id := uuid.Must(uuid.NewV1())
	return id.String()
}

func (g *IKGenerator) genUuidV4() string {
	id := uuid.Must(uuid.NewV4())
	return id.String()
}
