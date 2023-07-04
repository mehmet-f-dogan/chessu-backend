package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"

	"mehmetfd.dev/chessu-backend/database"
	"mehmetfd.dev/chessu-backend/models"
)

func AssignCompletionHandlers(app *fiber.App) {
	app.Post("/completion/content/:contentId/user/:userId/complete", handleCompleteContent)
	app.Get("/completion/course/:courseId/user/:userId/verify", handleVerifyCourseCompletion)
	app.Get("/completion/chapter/:chapterId/user/:userId/verify", handleVerifyChapterCompletion)
	app.Get("/completion/content/:contentId/user/:userId/verify", handleVerifyContentCompletion)
}

func handleCompleteContent(c *fiber.Ctx) error {
	clerkUserId := c.Params("userId")

	contentIdString := c.Params("contentId")
	contentId, err := uuid.Parse(contentIdString)
	if err != nil {
		return c.SendStatus(fiber.StatusOK)
	}

	var user models.AppUser
	if err := database.DB.Where(&models.AppUser{ClerkId: clerkUserId}).First(&user).Error; err != nil {
		return c.SendStatus(fiber.StatusOK)
	}

	// Find the content within the materials
	var foundContent *models.Content
	for _, course := range database.Materials {
		for _, chapter := range course.Chapters {
			for _, content := range chapter.Contents {
				if content.Id.Bytes == contentId {
					foundContent = &content
					break
				}
			}
		}
	}

	if foundContent == nil {
		return c.SendStatus(fiber.StatusOK)
	}

	// Check if the content is already completed
	for _, completedContentId := range user.CompletedContentId.Elements {
		if completedContentId.Bytes == contentId {
			return c.SendStatus(fiber.StatusOK)
		}
	}

	// Mark the content as completed
	user.CompletedContentId.Append(contentIdString)

	if err := database.DB.Save(&user).Error; err != nil {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusOK)
}

func handleVerifyCourseCompletion(c *fiber.Ctx) error {
	clerkUserId := c.Params("userId")

	courseIdString := c.Params("courseId")
	courseId, err := uuid.Parse(courseIdString)
	if err != nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	var user models.AppUser
	if err := database.DB.Where(&models.AppUser{ClerkId: clerkUserId}).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	coursePtr := database.GetCourse(courseId)

	if coursePtr == nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	// Check if all chapters in the course are completed
	for _, chapter := range coursePtr.Chapters {
		for _, content := range chapter.Contents {
			found := false
			for _, completedContentId := range user.CompletedContentId.Elements {
				if content.Id.Bytes == completedContentId.Bytes {
					found = true
					break
				}
			}
			if !found {
				return c.JSON(fiber.Map{
					"verified": false,
				})
			}
		}

	}

	return c.JSON(fiber.Map{
		"verified": true,
	})
}

func handleVerifyChapterCompletion(c *fiber.Ctx) error {
	clerkUserId := c.Params("userId")

	chapterIdString := c.Params("chapterId")
	chapterId, err := uuid.Parse(chapterIdString)
	if err != nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	var user models.AppUser
	if err := database.DB.Where(&models.AppUser{ClerkId: clerkUserId}).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	var chapterPtr *models.Chapter = nil

outer:
	for _, course := range database.Materials {
		for _, chapter := range course.Chapters {
			if chapter.Id.UUID.Bytes == chapterId {
				chapterPtr = &chapter
				break outer
			}
		}
	}

	if chapterPtr == nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	// Check if all contents in the chapter are completed
	for _, content := range chapterPtr.Contents {
		found := false
		for _, completedContentId := range user.CompletedContentId.Elements {
			if completedContentId.Bytes == content.Id.Bytes {
				found = true
				break
			}
		}
		if !found {
			return c.JSON(fiber.Map{
				"verified": false,
			})
		}
	}

	return c.JSON(fiber.Map{
		"verified": true,
	})
}

func handleVerifyContentCompletion(c *fiber.Ctx) error {
	clerkUserId := c.Params("userId")

	contentIdString := c.Params("contentId")
	contentId, err := uuid.Parse(contentIdString)
	if err != nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	var user models.AppUser
	if err := database.DB.Where(&models.AppUser{ClerkId: clerkUserId}).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	// Check if the content is completed
	for _, completedContentId := range user.CompletedContentId.Elements {
		if completedContentId.Bytes == contentId {
			return c.JSON(fiber.Map{
				"verified": true,
			})
		}
	}

	return c.JSON(fiber.Map{
		"verified": false,
	})
}
