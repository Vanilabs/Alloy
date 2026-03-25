package messaging

import (
	"alloy/internal/shared/router"
	"alloy/internal/shared/socket"

	"alloy/internal/shared/constants"

	"alloy/internal/shared/database/models"

	"alloy/internal/shared/middlewares"

	"strconv"

	"time"

	"context"

	"strings"

	"encoding/json"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

type Handler struct {
	service Service
	env     *router.Environment
	connectionManager *socket.ConnectionManager
}

func NewHandler(service Service, socketManager *socket.ConnectionManager) *Handler {
	return &Handler{
		service: service,
		connectionManager:  socketManager,
	}
}

func addOrUpdateAdmin(email string, dto *CreateChannelDTO) {
    found := false

    for i := range dto.Members {
        if dto.Members[i].Email == email {
            dto.Members[i].Role = constants.GroupAdmin
            found = true
            break 
        }
    }

    if !found {
        dto.Members = append(dto.Members, channelMember{
            Email: email,
            Role:  constants.GroupAdmin,
        })
    }
}


func (h *Handler) Init(basePath string, env *router.Environment) error {

	h.env = env

	chatGroup := env.Fiber.Group(basePath + "/chat")

	chatGroup.Use(middlewares.JWTMiddleware(env))

	chatGroup.Get("/send_message", websocket.New(h.handleActiveChat))
	chatGroup.Get("/messages", h.getConversationMessages)
	chatGroup.Get("/conversations", h.getActiveConversations)
	chatGroup.Post("/create", h.InitiateChat)
	chatGroup.Post("/messages/last_read", h.UpdateLastReadMessage)

	return nil
}


func (h *Handler) getConversationMessages(c *fiber.Ctx) error {

	conversationId := c.Query("conversation_id")
	limit := c.Query("limit")

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Limit must be a vaild Number"})
	}

	messages, err := h.service.ListConversationMessages(conversationId, parsedLimit)
	if err != nil {
		h.env.Logger.Error("Message History Retrieval Error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch messages"})
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}

func (h *Handler) handleActiveChat(c *websocket.Conn) {

	user, ok := c.Locals("user").(*models.User)
    if !ok {
        c.WriteJSON(fiber.Map{"error": "user context missing"})
		c.Close()
		return
    }

	ctx := context.Background()

	socketID := uuid.NewString()

	h.connectionManager.Add(socketID, c)
	if err := h.connectionManager.SocketTracker.AddSocketForUser(user.ID.String(), socketID); err != nil {
		h.connectionManager.Remove(socketID)
			_ = c.Close()
			return
		}

	h.env.Logger.Info("DM WebSocket connected", zap.String("userID", user.ID.String()))

	h.connectionManager.StartHeartbeat(ctx, socketID, c)

	c.SetPongHandler(func(string) error {
			return nil
		})

	offlineMessages, err := h.service.GetOfflineMessages(ctx, user.ID, 100)
	if err != nil {
		h.env.Logger.Error("Fetch Offline Message error", zap.Error(err))
		}
	if offlineMessages != nil || len(offlineMessages) != 0  {
		h.env.Logger.Info("Fetched Offline Message")
		for _, offmsg := range offlineMessages {
			
			if err != nil {
				h.env.Logger.Error("Fetch Offline Message error", zap.Error(err))
				continue
			}
			message := MessageMetadata {
				ID: GocqlToUUID(offmsg.ID),
			ConversationID: GocqlToUUID(offmsg.ConversationID),
			Text: offmsg.Text,
			Timestamp: offmsg.Timestamp,
			SenderID: GocqlToUUID(offmsg.SenderID),
			RecipientID: user.ID,
						}

			h.env.Logger.Info("Now delivering Offline Message")
			h.service.SendMessage(ctx, socketID, message)
		}
	}

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			h.env.Logger.Info("DM WebSocket disconnected", zap.String("userID", user.ID.String()))
			h.connectionManager.CleanupDeadSocket(socketID)
			break
		}
		var payload ChatDTO
		err = json.Unmarshal(msg, &payload)
		if err != nil {
    		c.WriteMessage(websocket.TextMessage, []byte("Invalid Payload"))
    		h.env.Logger.Warn("Invalid message payload", zap.String("userID", user.ID.String()), zap.Error(err))
    		continue
}		

		active_chat, err := h.service.ValidateConversation(ctx, payload.ConversationID)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Invalid Payload"))
			h.env.Logger.Warn("Invalid Conversation ID", zap.Error(err))
    		continue
		}

		msgID, err := h.service.SaveMessage(payload.Text, payload.ConversationID, user.ID, payload.Timestamp)
		if err != nil {
			h.env.Logger.Error("Failed to Save Message", zap.Error(err))
			c.WriteMessage(websocket.TextMessage, []byte("Message Not Delivered, Resend!!"))
		}

		members, err := h.service.GetChannelMembers(ctx, active_chat.channelId)
		if err != nil {
			h.env.Logger.Error("Failed to Retrieve Channel Members", zap.Error(err))
    		continue
		}

		for _, member := range members {
			
			message := MessageMetadata {
			ID: GocqlToUUID(*msgID),
			ConversationID: payload.ConversationID,
			Text: payload.Text,
			Timestamp: payload.Timestamp,
			SenderID: user.ID,
			RecipientID: member.UserID,
						}

			err = h.service.SendMessage(ctx, socketID, message)
			if err != nil {
				h.env.Logger.Error("Send Message Error", zap.Error(err))
				}

	}
	
	}
}


func (h *Handler) InitiateChat(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
    if !ok {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user context missing"})
    }
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	var payload CreateChannelDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	switch payload.Type {
	case constants.ChannelDM:
	if payload.Name != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "dm channels cannot have a name"})
	}
	if len(payload.Members)  != 1{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "can only add 1 participant to your dm"})
	}
	case constants.ChannelGroup:
	if payload.Name == nil || strings.TrimSpace(*payload.Name) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "group channels require a name"})
	}
}
	addOrUpdateAdmin(user.Email, &payload)

	chat_metadata, err := h.service.InitiateChat(ctx, payload); 
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(chat_metadata)
}


func (h *Handler) getActiveConversations(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user context missing"})
	}
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	conversations, err := h.service.GetActiveConversations(ctx, user.ID)
	if err != nil {
		h.env.Logger.Error("Failed to fetch active conversations", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch conversations"})
	}

	return c.Status(fiber.StatusOK).JSON(conversations)
}

func (h *Handler) UpdateLastReadMessage(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
    if !ok {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user context missing"})
    }
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	var payload UpdateLastReadDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.service.UpdateLastReadMessage(ctx, user.ID, payload.ConversationID,
				payload.MessageID, payload.ReadTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(nil)
}
