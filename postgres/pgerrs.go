package postgres

type IntegrityViolationError interface {
	error
	IsIntegrityViolationError()
}

type UniqueViolationError interface {
	IntegrityViolationError
	IsUniqueViolationError()
}

type RestrictViolationError interface {
	IntegrityViolationError
	IsRestrictViolationError()
}

type ForeignKeyViolationError interface {
	IntegrityViolationError
	IsForeignKeyViolationError()
}

type NotNullViolationError interface {
	IntegrityViolationError
	IsNotNullViolationError()
}

type CheckViolationError interface {
	IntegrityViolationError
	IsCheckViolationError()
}

type ExclusionViolationError interface {
	IntegrityViolationError
	IsExclusionViolationError()
}

type ErrIntegrityViolation struct{}

func (euv ErrIntegrityViolation) Error() string {
	return "integrity violation"
}
func (euv ErrIntegrityViolation) IsIntegrityViolationError() {}

type ErrUniqueViolation struct {
	Name string
}

func (euv ErrUniqueViolation) Error() string {
	return "unique violation " + euv.Name
}
func (euv ErrUniqueViolation) IsIntegrityViolationError() {}
func (euv ErrUniqueViolation) IsUniqueViolationError()    {}

type ErrRestrictViolation struct {
	Name string
}

func (erv ErrRestrictViolation) Error() string {
	return "restrict violation " + erv.Name
}
func (erv ErrRestrictViolation) IsIntegrityViolationError() {}
func (erv ErrRestrictViolation) IsRestrictViolationError()  {}

type ErrForeignKeyViolation struct {
	Name string
}

func (efv ErrForeignKeyViolation) Error() string {
	return "foreign key violation " + efv.Name
}
func (efv ErrForeignKeyViolation) IsIntegrityViolationError()  {}
func (efv ErrForeignKeyViolation) IsForeignKeyViolationError() {}

type ErrNotNullViolation struct {
	Name string
}

func (env ErrNotNullViolation) Error() string {
	return "not null violation " + env.Name
}
func (env ErrNotNullViolation) IsIntegrityViolationError()  {}
func (env ErrNotNullViolation) IsForeignKeyViolationError() {}

type ErrCheckViolation struct {
	Name string
}

func (ecv ErrCheckViolation) Error() string {
	return "check violation " + ecv.Name
}
func (ecv ErrCheckViolation) IsIntegrityViolationError() {}
func (ecv ErrCheckViolation) IsCheckViolationError()     {}

type ErrExclusionViolation struct {
	Name string
}

func (eev ErrExclusionViolation) Error() string {
	return "exclusion violation " + eev.Name
}
func (eev ErrExclusionViolation) IsIntegrityViolationError() {}
func (eev ErrExclusionViolation) IsExclusionViolationError() {}
