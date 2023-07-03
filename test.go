package main

import (
	"context"
	"log"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
)

func getDataType(value *ua.DataValue) string {
	if value.Status != ua.StatusOK {
		return value.Status.Error()
	}

	switch value.Value.NodeID().IntID() {
	case id.DateTime:
		return "time.Time"

	case id.Boolean:
		return "bool"

	case id.Int32:
		return "int32"
	case id.Float:
		return "float"
	}

	return value.Value.NodeID().String()
}
func main() {
	endpoint := "opc.tcp://GUROLGENCEL.SASA.LOCAL:4334/POC-Sample"
	nodeID := "ns=1;i=1001"

	ctx := context.Background()

	c := opcua.NewClient(endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		log.Fatalf("invalid node id: %v", err)
	}
	for i := 1; i <= 10; i++ {
		req := &ua.ReadRequest{
			MaxAge:             2000,
			NodesToRead:        []*ua.ReadValueID{{NodeID: id}},
			TimestampsToReturn: ua.TimestampsToReturnBoth,
		}

		resp, err := c.Read(req)
		if err != nil {
			log.Fatalf("Read failed: %s", err)
		}
		if resp.Results[0].Status != ua.StatusOK {
			log.Fatalf("Status not OK: %v", resp.Results[0].Status)
		}
		log.Printf("%#v", resp.Results[0].Value.Value())
		time.Sleep(1 * time.Second)
	}

}
