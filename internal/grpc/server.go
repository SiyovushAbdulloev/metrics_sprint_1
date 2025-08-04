package grpc

import (
	"context"
	"crypto/rsa"
	"log"
	"net"

	pb "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase"
	"github.com/golang/protobuf/ptypes/empty"
)

// GRPCServer реализует gRPC сервис.
type GRPCServer struct {
	pb.UnimplementedMetricsServiceServer
	uc         usecase.MetricUseCase
	privateKey *rsa.PrivateKey
}

// NewGRPCServer создаёт новый сервер.
func NewGRPCServer(uc usecase.MetricUseCase, priv *rsa.PrivateKey) *GRPCServer {
	return &GRPCServer{uc: uc, privateKey: priv}
}

// SendMetrics — обработчик RPC от агента.
func (s *GRPCServer) SendMetrics(ctx context.Context, req *pb.MetricsRequest) (*empty.Empty, error) {
	for _, m := range req.Metrics {
		// Декодируем данные
		metric := entity.Metrics{
			ID:    m.Id,
			MType: m.Type,
		}

		if m.Type == entity.Gauge {
			v := m.Value
			metric.Value = &v
		} else if m.Type == entity.Counter {
			d := m.Delta
			metric.Delta = &d
		}

		// Если данные приходят зашифрованными – расшифруем
		//if s.privateKey != nil && len(m.Encrypted) > 0 {
		//	decrypted, err := crypto.DecryptWithPrivateKey(m.Encrypted, s.privateKey)
		//	if err != nil {
		//		log.Printf("Decrypt error: %v", err)
		//		continue
		//	}
		//
		//	if err := json.Unmarshal(decrypted, &metric); err != nil {
		//		log.Printf("JSON unmarshal error: %v", err)
		//		continue
		//	}
		//}

		_, err := s.uc.StoreMetric(metric)
		if err != nil {
			log.Printf("StoreMetric error: %v", err)
		}
	}

	return &empty.Empty{}, nil
}

// RunGRPCServer запускает gRPC сервер на указанном адресе.
func RunGRPCServer(address string, uc usecase.MetricUseCase, priv *rsa.PrivateKey) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterMetricsServiceServer(grpcServer, NewGRPCServer(uc, priv))

	log.Printf("gRPC server started on %s", address)
	return grpcServer.Serve(lis)
}
