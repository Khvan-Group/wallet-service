package rabbitmq

import (
	"encoding/json"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/logger"
	"github.com/Khvan-Group/common-library/utils"
	"wallet-service/internal/models"
	"wallet-service/internal/service"
)

func ConsumeRabbitMQ(service service.WalletService) {
	go walletCreateQueueConsume(service)
	go walletUpdateQueueConsume(service)
}

func walletCreateQueueConsume(service service.WalletService) {
	msgs, err := RabbitMQChannel.Consume(
		utils.GetEnv(constants.RABBIT_CREATE_WALLET_QUEUE),
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logger.Logger.Fatal("Failed to consume RabbitMQ channel")
		return
	}

	go func() {
		for m := range msgs {
			var wallet models.Wallet

			if err = json.Unmarshal(m.Body, &wallet); err != nil {
				panic(err)
			}

			errSave := service.Save(wallet)
			if errSave != nil {
				panic(errSave)
			}
		}
	}()

	select {}
}

func walletUpdateQueueConsume(service service.WalletService) {
	msgs, err := RabbitMQChannel.Consume(
		utils.GetEnv(constants.RABBIT_UPDATE_WALLET_QUEUE),
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	go func() {
		for m := range msgs {
			var wallet models.WalletUpdate

			if err = json.Unmarshal(m.Body, &wallet); err != nil {
				panic(err)
			}

			errUpdate := service.Update(wallet)
			if errUpdate != nil {
				panic(errUpdate)
			}
		}
	}()

	select {}
}
