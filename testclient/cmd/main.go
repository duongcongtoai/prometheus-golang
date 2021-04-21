/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	pb "testclient/routeguide"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	port       = flag.Int("port", 10001, "http server port")
	metricPort = flag.Int("metric_port", 4001, "metric port")
)

// printFeature gets the feature for the given point.
func printFeature(client pb.RouteGuideClient, point *pb.Point) (*pb.Feature, error) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.GetFeature(ctx, point)
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}

func main() {
	flag.Parse()
	fmt.Println(*serverAddr)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *metricPort), nil))
	}()
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	var conn *grpc.ClientConn

	for i := 0; i < 3; i++ {
		tryconn, err := grpc.Dial(*serverAddr, opts...)
		if err != nil {
			fmt.Printf("fail to dial: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		conn = tryconn
		break
	}
	if conn == nil {
		log.Fatal("Failed to dial server")
	}

	fmt.Printf("Connected to grpc at %s\n", *serverAddr)
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)

	// Looking for a valid feature
	e := echo.New()
	e.GET("/get", getHandler(client))
	e.Start(fmt.Sprintf(":%d", *port))
}

func getHandler(client pb.RouteGuideClient) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		_, err := printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
		if err != nil {
			return err
		}
		return nil
	}
}
