package httpserver

type Option func(*Server)

func WithAddress(address string) Option {
	return func(s *Server) {
		s.Address = address
	}
}
