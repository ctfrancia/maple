package ports

// TournamentRepositoryProvider is an interface for providing thread safe access to the tournament repository
type TournamentRepositoryProvider interface {
	WriteTx(func(TournamentRepository) error) error
	ReadTx(func(TournamentRepository) error) error
}
