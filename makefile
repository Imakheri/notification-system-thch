MOCKGEN=mockgen
MOCKS_DEST=internal/mocks

.PHONY: mocks
mocks:
	@echo "Generating mocks on $(MOCKS_DEST)..."
	@mkdir -p $(MOCKS_DEST)
	#Repositories
	$(MOCKGEN) -source=internal/domain/gateway/notification_repository.go -destination=$(MOCKS_DEST)/mock_notification_repository.go -package=mocks
	$(MOCKGEN) -source=internal/domain/gateway/user_repository.go -destination=$(MOCKS_DEST)/mock_user_repository.go -package=mocks
	$(MOCKGEN) -source=internal/domain/gateway/channel_repository.go -destination=$(MOCKS_DEST)/mock_channel_repository.go -package=mocks
	#Services
	$(MOCKGEN) -source=internal/domain/gateway/jwt_service.go -destination=$(MOCKS_DEST)/mock_jwt_service.go -package=mocks
	$(MOCKGEN) -source=internal/domain/gateway/simulated_api_service.go -destination=$(MOCKS_DEST)/mock_simulated_api_service.go -package=mocks
	$(MOCKGEN) -source=internal/domain/gateway/clock.go -destination=$(MOCKS_DEST)/mock_clock.go -package=mocks
	#Usecases
	$(MOCKGEN) -source=internal/domain/usecase/create_user_usecase.go -destination=$(MOCKS_DEST)/mock_create_user_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/domain/usecase/create_notification_usecase.go -destination=$(MOCKS_DEST)/mock_create_notification_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/domain/usecase/get_notification_usecase.go -destination=$(MOCKS_DEST)/mock_get_notification_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/domain/usecase/get_user_usecase.go -destination=$(MOCKS_DEST)/mock_get_user_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/domain/usecase/update_user_usecase.go -destination=$(MOCKS_DEST)/mock_update_user_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/domain/usecase/update_notification_usecase.go -destination=$(MOCKS_DEST)/mock_update_notification_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/domain/usecase/delete_notification_usecase.go -destination=$(MOCKS_DEST)/mock_delete_notification_usecase.go -package=mocks