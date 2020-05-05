package main

import (
	"context"
	"log"
	"time"

	"github.com/go-sif/sif"
	"github.com/go-sif/sif/cluster"
	"github.com/go-sif/sif/datasource/file"
	"github.com/go-sif/sif/datasource/parser/jsonl"
	ops "github.com/go-sif/sif/operations/transform"
	"github.com/go-sif/sif/operations/util"
	"github.com/go-sif/sif/schema"
)

func main() {
	schema := schema.CreateSchema()
	schema.CreateColumn("coords.x", &sif.Float64ColumnType{})
	schema.CreateColumn("coords.z", &sif.Float64ColumnType{})
	schema.CreateColumn("date", &sif.TimeColumnType{Format: "2006-01-02 15:04:05"})

	parser := jsonl.CreateParser(&jsonl.ParserConf{
		PartitionSize: 128,
	})
	conf := &file.DataSourceConf{Glob: "/testenv/*.jsonl"}
	frame := file.CreateDataFrame(conf, parser, schema)

	frame, err := frame.To(
		ops.AddColumn("count", &sif.Uint32ColumnType{}),
		ops.Map(func(row sif.Row) error {
			err := row.SetInt32("count", int32(1))
			if err != nil {
				return err
			}
			return nil
		}),
		ops.Reduce(func(row sif.Row) ([]byte, error) {
			return []byte{byte(1)}, nil
		}, func(lrow sif.Row, rrow sif.Row) error {
			lval, err := lrow.GetInt32("count")
			if err != nil {
				return err
			}
			rval, err := rrow.GetInt32("count")
			if err != nil {
				return err
			}
			return lrow.SetInt32("count", lval+rval)
		}),
		util.Collect(1),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Define a node
	// Sif will read the SIF_NODE_TYPE environment variable to
	// determine whether this copy of the binary
	// is a "coordinator" or "worker".
	opts := &cluster.NodeOptions{
		NumWorkers:        2,
		CoordinatorHost:   "sif-coordinator",
		WorkerJoinTimeout: time.Duration(10) * time.Second,
	}
	node, err := cluster.CreateNode(opts)
	if err != nil {
		log.Fatal(err)
	}
	// start this node in the background and run the DataFrame
	defer node.GracefulStop()
	go func() {
		err := node.Start(frame)
		if err != nil {
			log.Fatal(err)
		}
	}()
	result, err := node.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// result only exists on coordinator
	if node.IsCoordinator() {
		for _, part := range result.Collected {
			row := part.GetRow(0)
			count, err := row.GetUint32("count")
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("There have been %d system discoveries in Elite Dangerous in the last 7 days.", count)
			break
		}
	}
}
