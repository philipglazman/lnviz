package main

import (
	"flag"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"github.com/philipglazman/lnviz/data"
	"github.com/philipglazman/lnviz/report"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	/* ------ Flags ------ */
	var certFile string
	var macaroonFile string
	var node string

	flag.StringVar(&certFile, "cert", "", "path to file containing TLS certificate for LND node")
	flag.StringVar(&macaroonFile, "macaroon", "", "path to file containing macaroon for LND node")
	flag.StringVar(&node, "host", "", "host:port of LND node")

	flag.Parse()

	// Connect to lightning client.
	nodeTLS, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatal(err)
	}

	macaroonBytes, err := ioutil.ReadFile(macaroonFile)
	if err != nil {
		log.Fatal(err)
	}
	nodeMacaroon := macaroon.Macaroon{}
	if err := nodeMacaroon.UnmarshalBinary(macaroonBytes); err != nil {
		log.Fatal(err)
	}

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(nodeTLS),
		grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(&nodeMacaroon)),
	}

	conn, err := grpc.Dial(node, options...)
	if err != nil {
		log.Fatal(err)
	}

	client := lnrpc.NewLightningClient(conn)

	events, nodes, err := data.BuildData(client)
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Create("report.html")

	// Create report.
	if err := components.NewPage().AddCharts(
		report.DailyRoutesProcessed(events),
		report.DailyRouteFees(events),
		report.DailyRoutingVolume(events),
		report.RouteFeePerChanOut(events, nodes),
		report.RouteFeePerChanIn(events, nodes),
		report.CumulativeRoutingFees(events),
	).SetLayout(components.PageFlexLayout).Render(io.MultiWriter(f)); err != nil {
		return
	}
}
