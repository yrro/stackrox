package main

import (
	"fmt"
	"log"
	"net"

	"github.com/boltdb/bolt"
	pb "github.com/stackrox/rox/generated/api/v1"
	sensorAPI "github.com/stackrox/rox/generated/internalapi/sensor"
	"google.golang.org/grpc"
)

var (
	port          = 9999
	dbPath        = "/tmp/collector-test.db"
	processBucket = "Process"
)

type signalServer struct {
	db *bolt.DB
}

func newServer(db *bolt.DB) *signalServer {
	return &signalServer{
		db: db,
	}
}

func (s *signalServer) PushSignals(stream sensorAPI.SignalService_PushSignalsServer) error {
	for {
		signal, err := stream.Recv()
		if err != nil {
			return err
		}
		var processSignal *pb.ProcessSignal
		if signal != nil && signal.GetSignal() != nil && signal.GetSignal().GetProcessSignal() != nil {
			processSignal = signal.GetSignal().GetProcessSignal()
		}

		fmt.Printf("%v\n", signal.GetSignal().GetProcessSignal())
		s.Update(processSignal)
	}
}

func boltDB(path string) (db *bolt.DB, err error) {
	db, err = bolt.Open(path, 0600, nil)
	return db, err
}

func (s *signalServer) Update(processSignal *pb.ProcessSignal) error {
	s.db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(processBucket))
		return b.Put([]byte(processSignal.Name), []byte(processSignal.ExecFilePath))
	})
	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	db, err := boltDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	grpcServer := grpc.NewServer()
	sensorAPI.RegisterSignalServiceServer(grpcServer, newServer(db))
	grpcServer.Serve(lis)
}
