package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"os"
	"softcery/internal/pkg/entity"
)

type rabbitRepository struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitRepository() (*rabbitRepository, error) {
	conn, err := createConnection()
	if err != nil {
		return nil, err
	}

	ch, err := createChannel(conn)
	if err != nil {
		return nil, err
	}

	_, err = setupImgQueue(ch)
	if err != nil {
		return nil, err
	}

	return &rabbitRepository{conn: conn, ch: ch}, nil
}

func createConnection() (*amqp.Connection, error) {
	rabbitUrl := os.Getenv("RABBIT_URL")
	if rabbitUrl == "" {
		return nil, errors.New("Env var 'RABBIT_URL' isn't specified!")
	}

	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}

	return conn, nil
}

func createChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open a rabbit channel: %v", err)
	}

	return ch, nil
}

func setupImgQueue(ch *amqp.Channel) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		ImgQueueName(),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue: %v", err)
	}
	return &q, nil
}

func (r *rabbitRepository) PushInImgQueue(body []byte) error {
	err := r.ch.Publish(
		"",
		ImgQueueName(),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("Failed to publish a message: %v", err)
	}
	return nil
}

func (r *rabbitRepository) ImgQueueConsumeChan(cxt context.Context) (<-chan *entity.Image, error) {
	msgs, err := r.ch.Consume(
		ImgQueueName(),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to register a consumer: %v", err)
	}

	return convertRabbitMsgChanToImgChan(cxt, msgs), nil
}

func convertRabbitMsgChanToImgChan(cxt context.Context, msgs <-chan amqp.Delivery) chan *entity.Image {
	imgs := make(chan *entity.Image)

	go func(ch chan *entity.Image) {
		for {
			select {
			case <-cxt.Done():
				close(ch)
				return
			case msg, ok := <-msgs:
				if !ok {
					close(ch)
					return
				}

				if img, err := convertMsgToImg(msg); err != nil {
					zap.S().Errorf("Invalid unmarshalling img from rabbit: %v", err)
				} else {
					ch <- img
				}
			}
		}
	}(imgs)

	return imgs
}

func convertMsgToImg(msg amqp.Delivery) (*entity.Image, error) {
	img := &entity.Image{}
	if err := json.Unmarshal(msg.Body, img); err != nil {
		return nil, err
	}

	return img, nil
}

func (r *rabbitRepository) Close() {
	err := r.conn.Close()
	if err != nil {
		zap.S().Errorf("Error during closing rabbit connection: %s", err)
	}

	err = r.ch.Close()
	if err != nil {
		zap.S().Errorf("Error during closing rabbit channel: %s", err)
	}
}

func ImgQueueName() string {
	return viper.GetString("rabbit.queue.img.name")
}
