package main

import (
	"context"
	"log"

	pb "github.com/harshav17/shippy/consignment-service/proto/consignment"
	vesselProto "github.com/harshav17/shippy/vessel-service/proto/vessel"
	"gopkg.in/mgo.v2"
)

type service struct {
	session      *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

//GetRepo -
func (s *service) GetRepo() Repository {
	return &ConsignmentRepository{s.session.Clone()}
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	defer s.GetRepo().Close()

	// Here we call a client instance of our vessel service with our consignment weight,
	// and the amount of containers as the capacity value
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
	if err = s.GetRepo().Create(req); err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req
	return nil
}

// GetConsignments -
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	defer s.GetRepo().Close()
	consignments, err := s.GetRepo().GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
