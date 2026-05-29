package user

// type Service struct {
// 	Store store.Store
// }

// func NewUserService(store store.Store) *Service {
// 	return &Service{
// 		Store: store,
// 	}
// }
// func (u *Service) RegistionUser(
// 	ctx context.Context,
// 	input dto.RegistrationInput,
// ) (*model.User, error) {
// 	if err := u.CheckEmail(ctx, input.Email); err != nil {
// 		return nil, err
// 	}

// 	hashPassword, err := hashPassword(input.HashPassword)
// 	if err != nil {
// 		return nil, fmt.Errorf("technical error: %w", err)
// 	}
// 	input.HashPassword = string(hashPassword)

// 	user, err := u.Store.CreateUser(ctx, input)
// 	if err != nil {
// 		return nil, fmt.Errorf("DB error: %w", err)
// 	}

// 	return user, nil
// }

// func (u *Service) UpdateUserInfo(
// 	ctx context.Context,
// 	input dto.UpdateUserInput,
// ) (*dto.UserResponse, error) {
// 	postUser, err := u.Store.UpdateUser(ctx, input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user := mapper.UserResponse(postUser)
// 	return &user, nil
// }

// func (u *Service) CheckEmail(ctx context.Context, Email string) error {
// 	user, err := u.Store.GetByEmail(ctx, Email)
// 	if err == nil {
// 		if user.ID > 0 {
// 			return fmt.Errorf("This %s email is taken", Email)
// 		}
// 		return nil
// 	}
// 	if !errors.Is(err, pgx.ErrNoRows) {
// 		return fmt.Errorf("technical error: %w", err)
// 	}
// 	return nil
// }

// func hashPassword(password string) ([]byte, error) {
// 	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return HashPassword, nil
// }
