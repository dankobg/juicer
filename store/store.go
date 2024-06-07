package store

type UserStore interface {
}

type GameStore interface {
	// Create(ctx context.Context, input models.CreateCountryInput) (*ent.Country, error)
	// Update(ctx context.Context, input models.UpdateCountryInput) (*ent.Country, error)
	// Delete(ctx context.Context, input models.DeleteCountryInput) (uuid.UUID, error)
	// BulkDelete(ctx context.Context, input models.BulkDeleteCountriesInput) (int, error)
	// Get(ctx context.Context, id uuid.UUID) (*ent.Country, error)
	// List(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, orderBy *ent.CountryOrder, where *ent.CountryWhereInput) (*ent.CountryConnection, error)
}

type Store interface {
	User() UserStore
	Game() GameStore
}
